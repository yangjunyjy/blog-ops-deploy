package api

import (
	"net/http"
	"time"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/middleware"
	"my-blog-backend/internal/pkg/session"
	"my-blog-backend/internal/pkg/token"
	"my-blog-backend/internal/pkg/utils"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SysAuthHandler RBAC认证Handler
type SysAuthHandler struct {
	userService       services.SysUserService
	permissionService services.SysPermissionService
	sessionManager    session.Manager
	tokenManager      token.Manager
}

func NewSysAuthHandler(
	userService services.SysUserService,
	permissionService services.SysPermissionService,
	sessionManager session.Manager,
	tokenManager token.Manager,
) *SysAuthHandler {
	return &SysAuthHandler{
		userService:       userService,
		permissionService: permissionService,
		sessionManager:    sessionManager,
		tokenManager:      tokenManager,
	}
}

// Login 管理后台登录（用户名+密码+邮箱验证码）
// @Summary 管理后台登录
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body request.SysLoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/auth/login [post]
func (h *SysAuthHandler) Login(c *gin.Context) {
	var req request.SysLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	logger.Info("Login attempt",
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("password", req.Password))

	// 验证用户
	logger.Info("Step 1: 开始验证用户", logger.String("username", req.Username))
	user, _, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		logger.Error("Step 1: 用户验证失败", logger.Err("eroor", err))
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误", err)
		return
	}
	logger.Info("Step 1: 用户验证成功", logger.Uint64("user_id", user.ID), logger.String("user_email", user.Email))

	// 验证邮箱是否匹配该用户
	logger.Info("Step 2: 验证邮箱匹配", logger.String("user_email", user.Email), logger.String("req_email", req.Email))
	if user.Email != req.Email {
		logger.Error("Step 2: 邮箱不匹配")
		response.Error(c, http.StatusUnauthorized, "邮箱与用户不匹配", nil)
		return
	}
	logger.Info("Step 2: 邮箱匹配成功")

	// 验证邮箱验证码
	logger.Info("Step 3: 开始验证验证码", logger.String("email", req.Email), logger.String("code", req.Code))
	verified, err := h.userService.VerifyEmailCode(c, req.Email, req.Code, true)
	if err != nil {
		logger.Error("Step 3: 验证码验证出错", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "验证验证码失败", err)
		return
	}
	logger.Info("Step 3: 验证码验证结果", logger.Bool("verified", verified))

	if !verified {
		logger.Error("Step 4: 验证码验证失败，返回401")
		response.Error(c, http.StatusUnauthorized, "验证码错误或已过期", nil)
		return
	}

	logger.Info("Step 4: 验证码验证成功，开始创建Session")
	// 创建Session并保存用户信息
	h.createSessionAndRespond(c, user)
}

// createSessionAndRespond 创建Session并返回响应
func (h *SysAuthHandler) createSessionAndRespond(c *gin.Context, user *models.SysUser) {
	// 获取用户角色
	userInfo, _ := h.userService.GetUserInfo(user.ID)

	var perms []string
	perms = utils.ExtractRolePerms(userInfo)
	logger.Info("权限码列表", logger.Any("perm_codes", perms))

	// 单会话模式：删除该用户的所有其他Session
	if err := h.sessionManager.DeleteAllByUserID(uint(user.ID)); err != nil {
		logger.Error("删除旧Session失败", logger.Err("error", err))
		// 不阻断登录流程，继续执行
		response.Error(c, http.StatusInternalServerError, "创建会话失败", err)
		return
	}

	// 创建Session
	sessionInfo := &session.SessionInfo{
		UserID:    uint(user.ID),
		Email:     user.Email,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: time.Now(),
		PermCodes: perms,
	}

	if len(userInfo.Roles) > 0 {
		var roleIDs []uint
		for _, role := range userInfo.Roles {
			roleIDs = append(roleIDs, uint(role.ID))
		}
		sessionInfo.RoleIDs = roleIDs
	}

	sessionID, err := h.sessionManager.Create(sessionInfo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建会话失败", err)
		return
	}

	// 设置Cookie
	h.sessionManager.SetCookie(c, sessionID)

	response.Success(c, gin.H{
		"user":  userInfo,
		"token": sessionID,
	}, "登录成功")
}

// Logout 登出
// @Summary 登出
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/auth/logout [post]
func (h *SysAuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	// 清除Session
	sessionID := h.sessionManager.GetCookie(c)
	if sessionID != "" {
		h.sessionManager.Delete(sessionID)
	}

	// 清除用户的Token
	h.tokenManager.BatchDelete(c, userID)

	// 调用用户服务登出
	h.userService.Logout(userID)

	// 清除Cookie
	h.sessionManager.ClearCookie(c)

	response.Success(c, nil, "登出成功")
}

// SendEmailCode 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body request.SysSendEmailCodeRequest true "邮箱信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/auth/send-email-code [post]
func (h *SysAuthHandler) SendEmailCode(c *gin.Context) {
	var req request.SysSendEmailCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 生成并发送验证码（根据邮箱查找用户以获取用户名）
	user, err := h.userService.GetUserInfoByEmail(req.Email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "邮箱未注册", nil)
		return
	}

	if err := h.userService.GenerateEmailCode(c, req.Email, user.Username); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "验证码已发送到邮箱")
}

// GetInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/auth/info [get]
func (h *SysAuthHandler) GetInfo(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在", err)
		return
	}
	logger.Info("controller层角色信息", logger.Any("user", user))
	response.Success(c, user, "获取成功")
}

// GetMenu 获取当前用户菜单
// @Summary 获取当前用户菜单
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/auth/menu [get]
func (h *SysAuthHandler) GetMenu(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	menus, err := h.permissionService.GetUserMenus(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取菜单失败", err)
		return
	}

	response.Success(c, menus, "获取成功")
}

// GenerateToken 生成一次性Token
// @Summary 生成一次性Token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/auth/generate-token [post]
func (h *SysAuthHandler) GenerateToken(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	username := middleware.GetCurrentUsername(c)

	// 检查是否启用TOTP
	if middleware.IsTOTPEnabled() {
		// 使用TOTP生成Token
		totpInstance := middleware.GetTOTPInstance()
		if totpInstance == nil {
			response.Error(c, http.StatusInternalServerError, "TOTP实例未初始化", nil)
			return
		}

		token, err := totpInstance.Generate()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "生成TOTP Token失败", err)
			return
		}

		// 返回TOTP Token（6位数字）
		response.Success(c, gin.H{
			"token": token,
			"type":  "totp",
		}, "生成成功")
		return
	}

	// 否则使用传统的Redis Token
	tokenInfo, err := h.tokenManager.Generate(c, userID, username, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "生成Token失败", err)
		return
	}

	response.Success(c, gin.H{
		"token_id":   tokenInfo.TokenID,
		"expires_at": tokenInfo.ExpiresAt,
		"type":       "redis",
	}, "生成成功")
}

// CheckPermission 检查权限
// @Summary 检查权限
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param menu_code query string true "菜单编码"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/auth/check-permission [get]
func (h *SysAuthHandler) CheckPermission(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	menuCode := c.Query("menu_code")
	if menuCode == "" {
		response.Error(c, http.StatusBadRequest, "缺少菜单编码", nil)
		return
	}

	hasPermission := h.permissionService.CheckPermission(userID, menuCode)

	response.Success(c, gin.H{
		"has_permission": hasPermission,
	}, "检查成功")
}

// ChangePassword 修改用户密码
func (h *SysAuthHandler) ChangePassword(c *gin.Context) {
	// TODO: 实现修改用户密码功能
}
