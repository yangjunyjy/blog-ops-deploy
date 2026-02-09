package mysql

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
)

type SysMenuRepositoryImpl struct {
	db *gorm.DB
}

func NewSysMenuRepositoryImpl(db *gorm.DB) repository.SysMenuRepository {
	return &SysMenuRepositoryImpl{db: db}
}

func (r *SysMenuRepositoryImpl) FindByID(id uint64) (*models.SysMenu, error) {
	var menu models.SysMenu
	err := r.db.Preload("Children").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *SysMenuRepositoryImpl) FindByCode(code string) (*models.SysMenu, error) {
	var menu models.SysMenu
	err := r.db.Where("menu_code = ?", code).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *SysMenuRepositoryImpl) FindByName(name string) (*models.SysMenu, error) {
	var menu models.SysMenu
	if name == "" {
		return nil, fmt.Errorf("查询名称为空")
	}
	err := r.db.Where("menu_name = ?", name).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *SysMenuRepositoryImpl) FindByPath(path string) (*models.SysMenu, error) {
	var menu models.SysMenu
	if path == "" {
		return nil, fmt.Errorf("查询路径为空")
	}
	err := r.db.Where("path = ?", path).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *SysMenuRepositoryImpl) List(page, pageSize int) ([]*models.SysMenu, int64, error) {
	var menus []*models.SysMenu
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&models.SysMenu{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

func (r *SysMenuRepositoryImpl) GetMenuTree() ([]*models.SysMenu, error) {
	var menus []*models.SysMenu
	err := r.db.Where("parent_id = 0").Order("sort ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 递归加载子菜单
	for _, menu := range menus {
		if err := r.loadChildren(menu); err != nil {
			return nil, err
		}
	}

	return menus, nil
}

func (r *SysMenuRepositoryImpl) GetUserMenus(userID uint64) ([]*models.SysMenu, error) {
	// 查询用户的所有角色
	var roleIDs []uint64
	if err := r.db.Table("sys_role").
		Joins("INNER JOIN sys_user_role ON sys_role.id = sys_user_role.role_id").
		Where("sys_user_role.user_id = ? AND sys_role.status = ?", userID, models.StatusEnabled).
		Pluck("sys_role.id", &roleIDs).Error; err != nil {
		return nil, err
	}

	if len(roleIDs) == 0 {
		return nil, nil
	}

	// 查询角色的所有菜单
	var menuIDs []uint64
	if err := r.db.Table("sys_menu").
		Joins("INNER JOIN sys_role_menu ON sys_menu.id = sys_role_menu.menu_id").
		Where("sys_role_menu.role_id IN ? AND sys_menu.status = ?", roleIDs, models.StatusEnabled).
		Pluck("sys_menu.id", &menuIDs).Error; err != nil {
		return nil, err
	}

	if len(menuIDs) == 0 {
		return nil, nil
	}

	// 包含所有有权限菜单的父级目录
	allMenuIDs := make([]uint64, 0, len(menuIDs))
	allMenuIDs = append(allMenuIDs, menuIDs...)

	// 递归查找所有父级目录
	for _, menuID := range menuIDs {
		parentIDs, err := r.getAllParentIDs(menuID)
		if err != nil {
			return nil, err
		}
		allMenuIDs = append(allMenuIDs, parentIDs...)
	}

	// 去重
	uniqueMenuIDs := make([]uint64, 0, len(allMenuIDs))
	seen := make(map[uint64]bool)
	for _, id := range allMenuIDs {
		if !seen[id] {
			seen[id] = true
			uniqueMenuIDs = append(uniqueMenuIDs, id)
		}
	}

	// 查询菜单详情并构建树
	var menus []*models.SysMenu
	err := r.db.Where("id IN ?", uniqueMenuIDs).Order("sort ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return r.buildMenuTree(menus, 0), nil
}

func (r *SysMenuRepositoryImpl) Create(menu *models.SysMenu) error {
	return r.db.Create(menu).Error
}

func (r *SysMenuRepositoryImpl) Update(menu *models.SysMenu) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		logger.Error("开启事务失败", logger.Err("error", tx.Error))
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Save(menu).Error; err != nil {
		tx.Rollback()
		logger.Error("更新时出现错误", logger.Err("error", err))
		return fmt.Errorf("仓库层更新失败 %v", err)
	}
	return tx.Commit().Error
}

func (r *SysMenuRepositoryImpl) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 递归删除子菜单
		if err := r.deleteChildren(tx, id); err != nil {
			return err
		}
		// 删除菜单
		if err := tx.Delete(&models.SysMenu{}, id).Error; err != nil {
			return err
		}
		// 删除角色菜单关联
		if err := tx.Where("menu_id = ?", id).Delete(&models.SysRoleMenu{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SysMenuRepositoryImpl) GetChildren(parentID uint64) ([]*models.SysMenu, error) {
	var menus []*models.SysMenu
	err := r.db.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *SysMenuRepositoryImpl) UpdateStatus(id uint64, status models.Status) error {
	return r.db.Model(&models.SysMenu{}).Where("id = ?", id).Update("status", status).Error
}

// loadChildren 递归加载子菜单
func (r *SysMenuRepositoryImpl) loadChildren(menu *models.SysMenu) error {
	children, err := r.GetChildren(menu.ID)
	if err != nil {
		return err
	}
	menu.Children = children

	for _, child := range children {
		if err := r.loadChildren(child); err != nil {
			return err
		}
	}

	return nil
}

// deleteChildren 递归删除子菜单
func (r *SysMenuRepositoryImpl) deleteChildren(tx *gorm.DB, parentID uint64) error {
	var children []*models.SysMenu
	if err := tx.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := r.deleteChildren(tx, child.ID); err != nil {
			return err
		}
		if err := tx.Delete(&models.SysMenu{}, child.ID).Error; err != nil {
			return err
		}
	}

	return nil
}

// getAllParentIDs 递归查找菜单的所有父级目录ID
func (r *SysMenuRepositoryImpl) getAllParentIDs(menuID uint64) ([]uint64, error) {
	var parentIDs []uint64
	var menu models.SysMenu

	// 查询当前菜单的父级
	if err := r.db.Select("parent_id").First(&menu, menuID).Error; err != nil {
		// 如果菜单不存在，返回空
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return parentIDs, nil
		}
		return nil, err
	}

	// 如果没有父级（根菜单），返回
	if menu.ParentID == 0 {
		return parentIDs, nil
	}

	// 添加父级ID
	parentIDs = append(parentIDs, menu.ParentID)

	// 递归查找父级的父级
	grandParentIDs, err := r.getAllParentIDs(menu.ParentID)
	if err != nil {
		return nil, err
	}
	parentIDs = append(parentIDs, grandParentIDs...)

	return parentIDs, nil
}

// buildMenuTree 构建菜单树
func (r *SysMenuRepositoryImpl) buildMenuTree(menus []*models.SysMenu, parentID uint64) []*models.SysMenu {
	var tree []*models.SysMenu
	var menuMap = make(map[uint64]*models.SysMenu)

	// 构建映射
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	// 构建树
	for _, menu := range menus {
		if menu.ParentID == parentID {
			tree = append(tree, menu)
		} else if parent, exists := menuMap[menu.ParentID]; exists {
			if parent.Children == nil {
				parent.Children = []*models.SysMenu{}
			}
			parent.Children = append(parent.Children, menu)
		}
	}

	return tree
}

// 获取所有菜单，不带分页
func (r *SysMenuRepositoryImpl) GetAllMenus() (menus []*models.SysMenu, err error) {
	err = r.db.Find(&menus).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("暂无菜单数据", logger.Err("error", err))
			return nil, nil
		}
		return nil, fmt.Errorf("获取菜单出错: %v", err)
	}
	return menus, nil
}
