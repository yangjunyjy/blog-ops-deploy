package api

import (
	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/middleware"
	"my-blog-backend/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SysUserHandler struct {
	sysUserService services.SysUserService
}

func NewSysUserHandler(sysUserService services.SysUserService) *SysUserHandler {
	return &SysUserHandler{
		sysUserService: sysUserService,
	}
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users/{id} [delete]
func (h *SysUserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	if err := h.sysUserService.DeleteUser(id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// Update 更新用户信息
// @Summary 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.UpdateSysUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users/update [put]
func (h *SysUserHandler) Update(c *gin.Context) {
	var req request.UpdateSysUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 获取当前操作用户ID
	operatorID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	// 查询用户是否存在
	user, err := h.sysUserService.GetUserInfo(req.Id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在", err)
		return
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Username = req.UserName
	user.Status = int8(req.Status)
	user.Gender = int8(req.Gender)

	// 如果提供了密码，则更新密码
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := h.sysUserService.UpdateUser(user, operatorID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, user, "更新成功")
}

// List 获取带分页的用户列表
// @Summary 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param username query string false "用户名"
// @Param email query string false "邮箱"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users [get]
func (h *SysUserHandler) List(c *gin.Context) {
	var query request.SearchSysUserQueryRequest

	// 设置默认值
	query.Page = 1
	query.PageSize = 10

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("参数错误", logger.Err("error", err))
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 获取分页数据
	users, total, err := h.sysUserService.List(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": total,
		"page":  query.Page,
		"size":  query.PageSize,
	}, "获取成功")
}

// Assign 分配角色给用户
// @Summary 分配角色给用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.AssignRolesRequest true "角色分配信息"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users/assign-roles [post]
func (h *SysUserHandler) Assign(c *gin.Context) {
	var req request.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := h.sysUserService.AssignRoles(req.UserID, req.RoleIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "分配角色成功")
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body request.CreateSysUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users [post]
func (h *SysUserHandler) Create(c *gin.Context) {
	var req request.CreateSysUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("请求参数有误", logger.Err("error", err), logger.Any("req", req))
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	logger.Info("解析后的请求参数", logger.Any("req", req))

	// 获取当前操作用户ID
	creatorID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	now := time.Now()
	user := &models.SysUser{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Status:   int8(req.Status),
		Gender:   int8(req.Gender),
		BaseModel: models.BaseModel{
			CreateTime: &now,
		},
	}

	if err := h.sysUserService.CreateUser(user, creatorID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, user, "创建成功")
}

// GetAllUser 获取所有用户,不带分页
// @Summary 获取所有用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/users/all [get]
func (h *SysUserHandler) GetAllUser(c *gin.Context) {
	users, total, err := h.sysUserService.List(&request.SearchSysUserQueryRequest{
		Page:     1,
		PageSize: 10000,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"list":  users,
		"total": total,
	}, "获取成功")
}

// GetUserByID 通过ID查询用户信息
func (h *SysUserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	user, err := h.sysUserService.GetUserInfo(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在", err)
		return
	}

	response.Success(c, user, "获取成功")
}
