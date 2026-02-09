package api

import (
	"net/http"
	"strconv"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dtoRequest.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	}

	if err := h.userService.CreateUser(user); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建用户失败", err)
		return
	}

	// 返回时不包含密码
	user.Password = ""
	response.Success(c, user, "创建成功")
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dtoRequest.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	var req dtoRequest.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	user := &models.User{
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := h.userService.UpdateUser(uint(id), user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body object true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	var req struct {
		Status uint8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	user := &models.User{
		Status: req.Status,
	}

	if err := h.userService.UpdateUser(uint(id), user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户状态失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除用户失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在", err)
		return
	}

	// 不返回密码
	user.Password = ""
	response.Success(c, user, "")
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/users/current [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	user, err := h.userService.GetUser(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在", err)
		return
	}

	// 不返回密码
	user.Password = ""
	response.Success(c, user, "获取成功")
}

// UpdateCurrentUser 更新当前用户信息
// @Summary 更新当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users/current [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	var req dtoRequest.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	user := &models.User{
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
		Bio:      req.Bio,
		Avatar:   req.Avatar,
	}

	if err := h.userService.UpdateUser(userID.(uint), user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户信息失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败", err)
		return
	}

	// 清空密码
	for _, user := range users {
		user.Password = ""
	}

	response.Success(c, gin.H{
		"items":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}
