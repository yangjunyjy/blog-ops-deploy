package services

import (
	"errors"
	"fmt"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

// SysMenuService 菜单服务
type SysMenuService interface {
	// List 列表
	List(page, pageSize int) ([]*models.SysMenu, int64, error)
	// GetTree 获取菜单树
	GetTree() ([]*models.SysMenu, error)
	// GetByID 获取详情
	GetByID(id uint64) (*models.SysMenu, error)
	// Create 创建
	Create(menu *models.SysMenu) error
	// Update 更新
	Update(menu *models.SysMenu) error
	// Delete 删除
	Delete(id uint64) error
	// GetUserMenus 获取用户菜单
	GetUserMenus(userID uint64) ([]*models.SysMenu, error)
	// 获取所有菜单
	GetAllMenus() ([]*models.SysMenu, error)
}

type sysMenuService struct {
	menuRepo repository.SysMenuRepository
}

func NewSysMenuService(menuRepo repository.SysMenuRepository) SysMenuService {
	return &sysMenuService{menuRepo: menuRepo}
}

func (s *sysMenuService) List(page, pageSize int) ([]*models.SysMenu, int64, error) {
	return s.menuRepo.List(page, pageSize)
}

func (s *sysMenuService) GetTree() ([]*models.SysMenu, error) {
	return s.menuRepo.GetMenuTree()
}

func (s *sysMenuService) GetByID(id uint64) (*models.SysMenu, error) {
	return s.menuRepo.FindByID(id)
}

func (s *sysMenuService) Create(menu *models.SysMenu) error {
	// 如果是按钮，检查菜单编码是否已存在
	if menu.MenuType == 3 {
		if _, errCode := s.menuRepo.FindByCode(menu.MenuCode); errCode == nil {
			return errors.New("菜单编码已存在")
		}
	}
	// 检测菜单名称是否存在
	if _, errName := s.menuRepo.FindByName(menu.MenuName); errName == nil {
		return errors.New("菜单名称已存在")
	}
	// 检查菜单路径是否已经被使用,不处理按钮
	if menu.MenuType != 3 {
		existing, errPath := s.menuRepo.FindByPath(menu.Path)
		if errPath == nil && existing.ID != menu.ID {
			return fmt.Errorf("菜单路径已被使用")
		}
	}
	return s.menuRepo.Create(menu)
}

func (s *sysMenuService) Update(menu *models.SysMenu) error {
	// 检查菜单编码是否被其他菜单使用
	if menu.MenuType == 3 {
		existing, errCode := s.menuRepo.FindByCode(menu.MenuCode)
		if errCode == nil && existing.ID != menu.ID {
			return errors.New("菜单编码已存在")
		}
	}
	// 检查菜单路径是否已经被使用,不处理按钮
	if menu.MenuType != 3 {
		existing, errPath := s.menuRepo.FindByPath(menu.Path)
		if errPath == nil && existing.ID != menu.ID {
			return fmt.Errorf("菜单路径已被使用")
		}
	}

	// 检查菜单名称是否已被使用
	existing, errName := s.menuRepo.FindByName(menu.MenuName)
	if errName == nil && existing.ID != menu.ID {
		return fmt.Errorf("菜单名称已被使用")
	}
	return s.menuRepo.Update(menu)
}

func (s *sysMenuService) Delete(id uint64) error {
	return s.menuRepo.Delete(id)
}

func (s *sysMenuService) GetUserMenus(userID uint64) ([]*models.SysMenu, error) {
	return s.menuRepo.GetUserMenus(userID)
}

func (s *sysMenuService) GetAllMenus() ([]*models.SysMenu, error) {
	return s.menuRepo.GetAllMenus()
}
