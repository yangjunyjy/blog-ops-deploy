package services

import (
	"errors"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/e"
	"my-blog-backend/internal/repository"
)

// SysRoleService 角色服务
type SysRoleService interface {
	// List 列表
	List(query *request.RoleSearchQueryRequest) ([]*models.SysRole, int64, error)
	// ListAll 所有角色
	ListAll() ([]*models.SysRole, error)
	// GetByID 获取详情
	GetByID(id uint64) (*models.SysRole, error)
	// Create 创建
	Create(role *models.SysRole) error
	// Update 更新
	Update(role *models.SysRole) error
	// Delete 删除
	Delete(id uint64) error
	// AssignMenus 分配菜单
	AssignMenus(req *request.RoleAssignMenusRequest) error
	// GetRoleMenus 获取角色菜单
	GetRoleMenus(roleID uint64) ([]*models.SysMenu, error)
}

type sysRoleService struct {
	roleRepo repository.SysRoleRepository
}

func NewSysRoleService(roleRepo repository.SysRoleRepository) SysRoleService {
	return &sysRoleService{roleRepo: roleRepo}
}

func (s *sysRoleService) List(query *request.RoleSearchQueryRequest) ([]*models.SysRole, int64, error) {
	if query.Page <= 0 || query.PageSize <= 0 {
		return nil, 0, e.ERROR_PARAMS
	}
	return s.roleRepo.List(query)
}

func (s *sysRoleService) ListAll() ([]*models.SysRole, error) {
	return s.roleRepo.ListAll()
}

func (s *sysRoleService) GetByID(id uint64) (*models.SysRole, error) {
	return s.roleRepo.FindByID(id)
}

func (s *sysRoleService) Create(role *models.SysRole) error {
	return s.roleRepo.Create(role)
}

func (s *sysRoleService) Update(role *models.SysRole) error {
	return s.roleRepo.Update(role)
}

func (s *sysRoleService) Delete(id uint64) error {
	// 不能删除超级管理员角色
	if id == models.SuperAdminRoleID {
		return errors.New("不能删除超级管理员角色")
	}
	return s.roleRepo.Delete(id)
}

func (s *sysRoleService) AssignMenus(req *request.RoleAssignMenusRequest) error {
	return s.roleRepo.AssignMenus(req.RoleID, req.MenuIDs)
}

func (s *sysRoleService) GetRoleMenus(roleID uint64) ([]*models.SysMenu, error) {
	return s.roleRepo.GetRoleMenus(roleID)
}
