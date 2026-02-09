package mysql

import (
	"time"

	"gorm.io/gorm"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

type SysRoleRepositoryImpl struct {
	db *gorm.DB
}

func NewSysRoleRepositoryImpl(db *gorm.DB) repository.SysRoleRepository {
	return &SysRoleRepositoryImpl{db: db}
}

func (r *SysRoleRepositoryImpl) FindByID(id uint64) (*models.SysRole, error) {
	var role models.SysRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *SysRoleRepositoryImpl) FindByName(name string) (*models.SysRole, error) {
	var role models.SysRole
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *SysRoleRepositoryImpl) List(query *request.RoleSearchQueryRequest) ([]*models.SysRole, int64, error) {
	var roles []*models.SysRole
	var total int64

	// 构建查询条件
	db := r.db.Model(&models.SysRole{})

	// 应用查询条件
	if query.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+query.RoleName+"%") // 通常支持模糊查询
	}
	if query.RoleCode != "" {
		db = db.Where("role_code like ?", "%"+query.RoleCode+"%")
	}
	if query.Status != 0 { // 假设0表示查询所有状态
		db = db.Where("status = ?", query.Status)
	}

	// 先统计总数（必须在分页之前）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("id ASC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *SysRoleRepositoryImpl) ListAll() ([]*models.SysRole, error) {
	var roles []*models.SysRole
	err := r.db.Where("status = ?", models.StatusEnabled).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *SysRoleRepositoryImpl) Create(role *models.SysRole) error {
	return r.db.Create(role).Error
}

func (r *SysRoleRepositoryImpl) Update(role *models.SysRole) error {
	return r.db.Save(role).Error
}

func (r *SysRoleRepositoryImpl) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Delete(&models.SysRole{}, id).Error; err != nil {
			return err
		}
		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&models.SysRoleMenu{}).Error; err != nil {
			return err
		}
		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&models.SysUserRole{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SysRoleRepositoryImpl) AssignMenus(roleID uint64, menuIDs []uint64) error {
	// 先删除现有菜单
	if err := r.db.Where("role_id = ?", roleID).Delete(&models.SysRoleMenu{}).Error; err != nil {
		return err
	}

	// 批量插入新菜单
	if len(menuIDs) == 0 {
		return nil
	}

	now := time.Now()
	var roleMenus []models.SysRoleMenu
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, models.SysRoleMenu{
			RoleID:    roleID,
			MenuID:    menuID,
			CreatedAt: now,
		})
	}

	return r.db.CreateInBatches(roleMenus, 100).Error
}

func (r *SysRoleRepositoryImpl) GetRoleMenus(roleID uint64) ([]*models.SysMenu, error) {
	var role models.SysRole
	err := r.db.Preload("Menus").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Menus, nil
}

func (r *SysRoleRepositoryImpl) GetRoleUsers(roleID uint64) ([]*models.SysUser, error) {
	var users []*models.SysUser
	err := r.db.Table("sys_user").
		Joins("INNER JOIN sys_user_role ON sys_user.id = sys_user_role.user_id").
		Where("sys_user_role.role_id = ?", roleID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *SysRoleRepositoryImpl) UpdateStatus(id uint64, status models.Status) error {
	return r.db.Model(&models.SysRole{}).Where("id = ?", id).Update("status", status).Error
}
