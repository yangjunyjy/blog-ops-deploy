package api

import (
	"fmt"
	"net/http"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	dtoResponse "my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/pkg/captcha"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	// 初始化验证码存储，默认5分钟过期
	captcha.InitCaptchaStore(300)
	return &AuthHandler{
		authService: authService,
	}
}

// GetCaptcha 获取验证码
// @Summary 获取验证码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param type query int false "验证码类型" default(1)
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/captcha [get]
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	var req dtoRequest.CaptchaRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Type = 1 // 默认使用数字验证码
	}

	// 验证码参数配置
	captchaType := req.Type
	length := 4  // 验证码长度
	width := 140 // 图片宽度
	height := 44 // 图片高度

	id, b64s, err := captcha.GenerateCaptcha(captchaType, length, width, height)
	if err != nil {
		dtoResponse.Error(c, http.StatusInternalServerError, "生成验证码失败", err)
		return
	}

	resp := &dtoResponse.CaptchaResponse{
		CaptchaID: id,
		Image:     b64s,
	}
	dtoResponse.Success(c, resp, "获取成功")
}

// Login 登录
// @Summary 用户名密码登录
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.LoginRequest true "登录信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) LoginWithCode(c *gin.Context) {
	var req dtoRequest.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 验证验证码
	if !captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaAnswer) {
		dtoResponse.Error(c, http.StatusBadRequest, "验证码错误或已过期", nil)
		return
	}

	// 打印调试信息
	fmt.Printf("🔍 登录请求: 用户名=%s, 密码长度=%d\n", req.Username, len(req.Password))

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		dtoResponse.Error(c, http.StatusUnauthorized, "用户名或密码错误", err)
		return
	}

	fmt.Printf("登录成功: 用户名=%s\n", req.Username)

	// 获取用户信息
	userInfo, err := h.authService.GetUserByUsername(req.Username)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		dtoResponse.Error(c, http.StatusInternalServerError, "获取用户信息失败", err)
		return
	}

	resp := &dtoResponse.AuthResponse{
		Token: token,
		User:  *userInfo,
	}

	// 打印token前10个字符和用户ID用于调试
	fmt.Printf("🔑 登录响应: 用户ID=%d, Token前10字符=%.10s...\n", userInfo.ID, token)

	dtoResponse.Success(c, resp, "登录成功")
}

// LoginWithEmail 邮箱验证码登录
// @Summary 邮箱验证码登录
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.LoginByEmailRequest true "登录信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/login/email [post]
func (h *AuthHandler) LoginWithEmail(c *gin.Context) {
	var req dtoRequest.LoginByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 调用服务进行邮箱登录
	token, err := h.authService.LoginByEmail(req.Email, req.Code)
	if err != nil {
		dtoResponse.Error(c, http.StatusUnauthorized, err.Error(), err)
		return
	}

	// 根据邮箱获取用户信息
	userInfo, err := h.authService.GetUserByEmail(req.Email)
	if err != nil {
		fmt.Printf("❌ 获取用户信息失败: %v\n", err)
		dtoResponse.Error(c, http.StatusInternalServerError, "获取用户信息失败", err)
		return
	}

	resp := &dtoResponse.AuthResponse{
		Token: token,
		User:  *userInfo,
	}
	dtoResponse.Success(c, resp, "登录成功")
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body object true "令牌信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshRequest struct {
		Token string `json:"token" binding:"required"`
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		dtoResponse.Error(c, http.StatusUnauthorized, "令牌无效或已过期", err)
		return
	}

	resp := &dtoResponse.LoginResponse{
		Token: newToken,
	}
	dtoResponse.Success(c, resp, "刷新成功")
}

// SendEmailCaptcha 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.EmailCaptchaRequest true "邮箱信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/email/code [post]
func (h *AuthHandler) SendEmailCaptcha(c *gin.Context) {

	var req dtoRequest.EmailCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 如果没有提供用户名，使用邮箱作为默认值
	username := req.Username
	if username == "" {
		username = req.Email
	}

	// 调用服务发送验证码
	err := h.authService.GenerateSmtpCaptcha(c.Request.Context(), req.Email, username)
	if err != nil {
		dtoResponse.Error(c, http.StatusInternalServerError, "发送验证码失败: "+err.Error(), err)
		return
	}

	dtoResponse.Success(c, nil, "验证码已发送")
}

// VerifyEmailCaptcha 验证邮箱验证码
func (h *AuthHandler) VerifyEmailCaptcha(c *gin.Context) {
	var req dtoRequest.VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 调用服务验证验证码
	verified, err := h.authService.VerifySmtpCaptcha(c.Request.Context(), req.Email, req.Code, false)
	if err != nil {
		dtoResponse.Error(c, http.StatusInternalServerError, "验证验证码失败", err)
		return
	}

	if !verified {
		dtoResponse.Error(c, http.StatusBadRequest, "验证码错误或已过期", nil)
		return
	}

	dtoResponse.Success(c, gin.H{"verified": true}, "验证成功")
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.RegisterRequest true "注册信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dtoRequest.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	right, err := h.authService.VerifySmtpCaptcha(c, req.Email, req.Code, false)
	if err != nil {
		dtoResponse.Error(c, http.StatusInternalServerError, "服务器内部错误", err)
		return
	}
	if !right {
		dtoResponse.Error(c, http.StatusBadRequest, "验证码错误或已过期", nil)
		return
	}
	if err := h.authService.Register(c, req.Username, req.Password, req.Email); err != nil {
		dtoResponse.Error(c, http.StatusInternalServerError, "服务器内部错误", err)
		return
	}
	dtoResponse.Success(c, gin.H{"register": true}, "注册成功")
}
