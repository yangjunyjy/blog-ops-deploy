package repository

import (
	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/models"
)

// SysRoleRepository 系统角色仓储接口
type SysRoleRepository interface {
	// FindByID 根据ID查找角色
	FindByID(id uint64) (*models.SysRole, error)
	// FindByName 根据角色名称查找角色
	FindByName(name string) (*models.SysRole, error)
	// List 查询角色列表
	List(query *request.RoleSearchQueryRequest) ([]*models.SysRole, int64, error)
	// ListAll 查询所有角色
	ListAll() ([]*models.SysRole, error)
	// Create 创建角色
	Create(role *models.SysRole) error
	// Update 更新角色
	Update(role *models.SysRole) error
	// Delete 删除角色
	Delete(id uint64) error
	// AssignMenus 分配菜单给角色
	AssignMenus(roleID uint64, menuIDs []uint64) error
	// GetRoleMenus 获取角色的菜单列表
	GetRoleMenus(roleID uint64) ([]*models.SysMenu, error)
	// GetRoleUsers 获取角色的用户列表
	GetRoleUsers(roleID uint64) ([]*models.SysUser, error)
	// UpdateStatus 更新角色状态
	UpdateStatus(id uint64, status models.Status) error
}
