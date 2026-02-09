package repository

import (
	"my-blog-backend/internal/models"
)

// SysMenuRepository 系统菜单仓储接口
type SysMenuRepository interface {
	// FindByID 根据ID查找菜单
	FindByID(id uint64) (*models.SysMenu, error)
	// FindByCode 根据菜单编码查找菜单
	FindByCode(code string) (*models.SysMenu, error)
	// FindByName 根据菜单名称查找菜单
	FindByName(name string) (*models.SysMenu, error)
	// FindByPath 根据菜单路径查找菜单
	FindByPath(path string) (*models.SysMenu, error)
	// List 查询菜单列表
	List(page, pageSize int) ([]*models.SysMenu, int64, error)
	// GetMenuTree 获取菜单树
	GetMenuTree() ([]*models.SysMenu, error)
	// GetUserMenus 获取用户的菜单列表
	GetUserMenus(userID uint64) ([]*models.SysMenu, error)
	// Create 创建菜单
	Create(menu *models.SysMenu) error
	// Update 更新菜单
	Update(menu *models.SysMenu) error
	// Delete 删除菜单
	Delete(id uint64) error
	// GetChildren 获取子菜单
	GetChildren(parentID uint64) ([]*models.SysMenu, error)
	// UpdateStatus 更新菜单状态
	UpdateStatus(id uint64, status models.Status) error
	// GetAllMenus 获取所有菜单
	GetAllMenus() ([]*models.SysMenu, error)
}
