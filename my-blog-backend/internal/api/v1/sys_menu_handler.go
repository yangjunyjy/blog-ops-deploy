package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SysMenuHandler 菜单Handler
type SysMenuHandler struct {
	menuService services.SysMenuService
}

func NewSysMenuHandler(menuService services.SysMenuService) *SysMenuHandler {
	return &SysMenuHandler{menuService: menuService}
}

// List 菜单列表
// @Summary 获取菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/menus [get]
func (h *SysMenuHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	menus, total, err := h.menuService.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取菜单列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"list":  menus,
		"total": total,
	}, "获取成功")
}

// GetTree 获取菜单树
// @Summary 获取菜单树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/menus/tree [get]
func (h *SysMenuHandler) GetTree(c *gin.Context) {
	menus, err := h.menuService.GetTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取菜单树失败", err)
		return
	}

	response.Success(c, menus, "获取成功")
}

// GetByID 获取菜单详情
// @Summary 获取菜单详情
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/menus/{id} [get]
func (h *SysMenuHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID", err)
		return
	}

	menu, err := h.menuService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "菜单不存在", err)
		return
	}

	response.Success(c, menu, "获取成功")
}

// Create 创建菜单
// @Summary 创建菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body request.CreateSysMenuRequest true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/menus/create [post]
func (h *SysMenuHandler) Create(c *gin.Context) {
	var menu request.CreateSysMenuRequest
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	value, ok := c.Get("user_id")
	if !ok {
		logger.Error("创建人信息获取出错", logger.Err("error", errors.New("上下文获取用户id出错")))
		response.Error(c, http.StatusInternalServerError, "参数错误", errors.New("上下文获取用户id出错"))
		return
	}
	user_id, ok2 := value.(uint)
	if !ok2 {
		logger.Error("类型断言用户id出错")
		response.Error(c, http.StatusInternalServerError, "参数错误", errors.New("类型断言用户id出错"))
		return
	}
	var create_time time.Time
	create_time = time.Now()
	if err := h.menuService.Create(&models.SysMenu{
		BaseModel: models.BaseModel{
			CreateBy:   uint64(user_id),
			CreateTime: &create_time,
		},
		ParentID:  menu.ParentID,
		Component: menu.Component,
		MenuName:  menu.Name,
		MenuCode:  menu.Code,
		MenuType:  models.MenuType(menu.Type),
		Sort:      menu.Sort,
		Status:    models.Status(menu.Status),
		Icon:      menu.Icon,
		Path:      menu.Path,
		IsVisible: menu.IsVisible,
		Remark:    menu.Description,
	}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, menu, "创建成功")
}

// Update 更新菜单
// @Summary 更新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param request body models.SysMenu true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/v1/rbac/menus/update [put]
func (h *SysMenuHandler) Update(c *gin.Context) {
	var menu request.UpdateSysMenuRequest
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}
	value, ok := c.Get("user_id")
	if !ok {
		logger.Error("创建人信息获取出错", logger.Err("error", errors.New("上下文获取用户id出错")))
		response.Error(c, http.StatusInternalServerError, "参数错误", errors.New("上下文获取用户id出错"))
		return
	}
	user_id, ok2 := value.(uint)
	if !ok2 {
		logger.Error("类型断言用户id出错")
		response.Error(c, http.StatusInternalServerError, "参数错误", errors.New("类型断言用户id出错"))
		return
	}
	var update_time time.Time
	update_time = time.Now()
	if err := h.menuService.Update(&models.SysMenu{
		BaseModel: models.BaseModel{
			ID:         menu.ID,
			UpdateTime: &update_time,
			UpdateBy:   uint64(user_id),
		},
		ParentID:  menu.ParentID,
		Component: menu.Component,
		MenuName:  menu.Name,
		MenuCode:  menu.Code,
		MenuType:  models.MenuType(menu.Type),
		Sort:      menu.Sort,
		Status:    models.Status(menu.Status),
		Icon:      menu.Icon,
		Path:      menu.Path,
		IsVisible: menu.IsVisible,
		Remark:    menu.Description,
	}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, menu, "更新成功")
}

// Delete 删除菜单
// @Summary 删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/menus/{id} [delete]
func (h *SysMenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID", err)
		return
	}

	if err := h.menuService.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetAllMenus 获取所有菜单，不分页，不返回树结构
// @Summary 获取所有菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/menus/all [get]
func (h *SysMenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.menuService.GetAllMenus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取所有菜单出错", err)
		return
	}
	response.Success(c, menus, "获取所有菜单成功")
}
