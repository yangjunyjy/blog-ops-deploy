package services

import (
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

// SysPermissionService 权限服务
type SysPermissionService interface {
	// CheckPermission 检查用户是否有指定权限
	CheckPermission(userID uint64, menuCode string) bool
	// GetUserPermissions 获取用户所有权限
	GetUserPermissions(userID uint64) ([]string, error)
	// GetUserMenus 获取用户菜单树
	GetUserMenus(userID uint64) ([]*models.SysMenu, error)
	// IsSuperAdmin 判断是否为超级管理员
	IsSuperAdmin(userID uint64) bool
}

type sysPermissionService struct {
	userRepo repository.SysUserRepository
	roleRepo repository.SysRoleRepository
	menuRepo repository.SysMenuRepository
}

func NewSysPermissionService(
	userRepo repository.SysUserRepository,
	roleRepo repository.SysRoleRepository,
	menuRepo repository.SysMenuRepository,
) SysPermissionService {
	return &sysPermissionService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		menuRepo: menuRepo,
	}
}

func (s *sysPermissionService) CheckPermission(userID uint64, menuCode string) bool {
	// 超级管理员拥有所有权限
	if s.IsSuperAdmin(userID) {
		return true
	}

	// 获取用户菜单
	menus, err := s.GetUserMenus(userID)
	if err != nil || len(menus) == 0 {
		return false
	}

	// 递归检查菜单树
	return s.checkMenuPermission(menus, menuCode)
}

func (s *sysPermissionService) GetUserPermissions(userID uint64) ([]string, error) {
	// 超级管理员拥有所有权限
	if s.IsSuperAdmin(userID) {
		return []string{"*"}, nil
	}

	menus, err := s.GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	var permissions []string
	for _, menu := range menus {
		permissions = append(permissions, s.extractPermissions(menu)...)
	}

	return permissions, nil
}

func (s *sysPermissionService) GetUserMenus(userID uint64) ([]*models.SysMenu, error) {
	return s.menuRepo.GetUserMenus(userID)
}

func (s *sysPermissionService) IsSuperAdmin(userID uint64) bool {
	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if role.ID == models.SuperAdminRoleID {
			return true
		}
	}

	return false
}

func (s *sysPermissionService) checkMenuPermission(menus []*models.SysMenu, menuCode string) bool {
	for _, menu := range menus {
		if menu.MenuCode == menuCode {
			return true
		}
		if len(menu.Children) > 0 {
			if s.checkMenuPermission(menu.Children, menuCode) {
				return true
			}
		}
	}
	return false
}

func (s *sysPermissionService) extractPermissions(menu *models.SysMenu) []string {
	var permissions []string
	if menu.MenuCode != "" {
		permissions = append(permissions, menu.MenuCode)
	}
	if menu.Permission != "" {
		permissions = append(permissions, menu.Permission)
	}
	for _, child := range menu.Children {
		permissions = append(permissions, s.extractPermissions(child)...)
	}
	return permissions
}
