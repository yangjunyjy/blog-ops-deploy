package mysql

import (
	"gorm.io/gorm"
	"strings"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

type SysDeptRepositoryImpl struct {
	db *gorm.DB
}

func NewSysDeptRepositoryImpl(db *gorm.DB) repository.SysDeptRepository {
	return &SysDeptRepositoryImpl{db: db}
}

func (r *SysDeptRepositoryImpl) FindByID(id uint64) (*models.SysDept, error) {
	var dept models.SysDept
	err := r.db.Preload("Children").First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *SysDeptRepositoryImpl) FindByName(name string) (*models.SysDept, error) {
	var dept models.SysDept
	err := r.db.Where("name = ?", name).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *SysDeptRepositoryImpl) List(page, pageSize int) ([]*models.SysDept, int64, error) {
	var depts []*models.SysDept
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&models.SysDept{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&depts).Error
	if err != nil {
		return nil, 0, err
	}

	return depts, total, nil
}

func (r *SysDeptRepositoryImpl) GetDeptTree() ([]*models.SysDept, error) {
	var depts []*models.SysDept
	err := r.db.Where("parent_id = 0").Order("sort ASC, id ASC").Find(&depts).Error
	if err != nil {
		return nil, err
	}

	// 递归加载子部门
	for _, dept := range depts {
		if err := r.loadChildren(dept); err != nil {
			return nil, err
		}
	}

	return depts, nil
}

func (r *SysDeptRepositoryImpl) Create(dept *models.SysDept) error {
	return r.db.Create(dept).Error
}

func (r *SysDeptRepositoryImpl) Update(dept *models.SysDept) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新部门
		if err := tx.Save(dept).Error; err != nil {
			return err
		}

		// 如果修改了父部门，更新 ancestors 字段
		if dept.ParentID > 0 {
			var parent models.SysDept
			if err := tx.First(&parent, dept.ParentID).Error; err != nil {
				return err
			}
			dept.Ancestors = parent.Ancestors + "," + strings.TrimPrefix(parent.Ancestors, "0")
		}

		return nil
	})
}

func (r *SysDeptRepositoryImpl) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 递归删除子部门
		if err := r.deleteChildren(tx, id); err != nil {
			return err
		}
		// 删除部门
		return tx.Delete(&models.SysDept{}, id).Error
	})
}

func (r *SysDeptRepositoryImpl) GetChildren(parentID uint64) ([]*models.SysDept, error) {
	var depts []*models.SysDept
	err := r.db.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *SysDeptRepositoryImpl) GetAncestors(id uint64) (string, error) {
	var dept models.SysDept
	err := r.db.First(&dept, id).Error
	if err != nil {
		return "", err
	}
	return dept.Ancestors, nil
}

func (r *SysDeptRepositoryImpl) UpdateStatus(id uint64, status models.Status) error {
	return r.db.Model(&models.SysDept{}).Where("id = ?", id).Update("status", status).Error
}

func (r *SysDeptRepositoryImpl) loadChildren(dept *models.SysDept) error {
	children, err := r.GetChildren(dept.ID)
	if err != nil {
		return err
	}
	dept.Children = children

	for _, child := range children {
		if err := r.loadChildren(child); err != nil {
			return err
		}
	}

	return nil
}

func (r *SysDeptRepositoryImpl) deleteChildren(tx *gorm.DB, parentID uint64) error {
	var children []*models.SysDept
	if err := tx.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := r.deleteChildren(tx, child.ID); err != nil {
			return err
		}
		if err := tx.Delete(&models.SysDept{}, child.ID).Error; err != nil {
			return err
		}
	}

	return nil
}
