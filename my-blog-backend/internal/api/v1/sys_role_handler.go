package api

import (
	"net/http"
	"strconv"
	"time"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SysRoleHandler 角色Handler
type SysRoleHandler struct {
	roleService services.SysRoleService
}

func NewSysRoleHandler(roleService services.SysRoleService) *SysRoleHandler {
	return &SysRoleHandler{roleService: roleService}
}

// List 角色列表
// @Summary 获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param roleName query string false "角色名称"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles [get]
func (h *SysRoleHandler) List(c *gin.Context) {
	var query request.RoleSearchQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误", err)
	}
	roles, total, err := h.roleService.List(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"list":  roles,
		"total": total,
	}, "获取成功")
}

// ListAll 所有角色
// @Summary 获取所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/all [get]
func (h *SysRoleHandler) ListAll(c *gin.Context) {
	roles, err := h.roleService.ListAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色列表失败", err)
		return
	}

	response.Success(c, roles, "获取成功")
}

// GetByID 获取角色详情
// @Summary 获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{id} [get]
func (h *SysRoleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	role, err := h.roleService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "角色不存在", err)
		return
	}

	response.Success(c, role, "获取成功")
}

// Create 创建角色
// @Summary 创建角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body request.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles [post]
func (h *SysRoleHandler) Create(c *gin.Context) {
	var role request.CreateRoleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	var now time.Time
	now = time.Now()
	if err := h.roleService.Create(&models.SysRole{
		RoleCode: role.RoleCode,
		RoleName: role.RoleName,
		RoleDesc: role.RoleDesc,
		Status:   role.Status,
		BaseModel: models.BaseModel{
			CreateTime: &now,
		},
	}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, role, "创建成功")
}

// Update 更新角色
// @Summary 更新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body request.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{id} [put]
func (h *SysRoleHandler) Update(c *gin.Context) {
	var role request.UpdateRoleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	var now time.Time
	now = time.Now()
	if err := h.roleService.Update(&models.SysRole{
		RoleCode: role.RoleCode,
		RoleName: role.RoleName,
		RoleDesc: role.RoleDesc,
		Status:   role.Status,
		BaseModel: models.BaseModel{
			ID:         uint64(role.ID),
			UpdateTime: &now,
		},
	}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, role, "更新成功")
}

// Delete 删除角色
// @Summary 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{id} [delete]
func (h *SysRoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	if err := h.roleService.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// AssignMenus 分配菜单
// @Summary 分配菜单给角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body request.RoleAssignMenusRequest true "菜单分配信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/assign-menus [post]
func (h *SysRoleHandler) AssignMenus(c *gin.Context) {
	var req request.RoleAssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	if err := h.roleService.AssignMenus(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "分配成功")
}

// GetMenus 获取角色菜单
// @Summary 获取角色的菜单
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{id}/menus [get]
func (h *SysRoleHandler) GetMenus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	menus, err := h.roleService.GetRoleMenus(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, menus, "获取成功")
}
