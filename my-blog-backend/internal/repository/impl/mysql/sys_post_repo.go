package mysql

import (
	"gorm.io/gorm"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

type SysPostRepositoryImpl struct {
	db *gorm.DB
}

func NewSysPostRepositoryImpl(db *gorm.DB) repository.SysPostRepository {
	return &SysPostRepositoryImpl{db: db}
}

func (r *SysPostRepositoryImpl) FindByID(id uint64) (*models.SysPost, error) {
	var post models.SysPost
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *SysPostRepositoryImpl) FindByName(name string) (*models.SysPost, error) {
	var post models.SysPost
	err := r.db.Where("name = ?", name).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *SysPostRepositoryImpl) FindByCode(code string) (*models.SysPost, error) {
	var post models.SysPost
	err := r.db.Where("code = ?", code).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *SysPostRepositoryImpl) List(page, pageSize int) ([]*models.SysPost, int64, error) {
	var posts []*models.SysPost
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&models.SysPost{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *SysPostRepositoryImpl) ListAll() ([]*models.SysPost, error) {
	var posts []*models.SysPost
	err := r.db.Where("status = ?", models.StatusEnabled).Order("sort ASC, id ASC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *SysPostRepositoryImpl) Create(post *models.SysPost) error {
	return r.db.Create(post).Error
}

func (r *SysPostRepositoryImpl) Update(post *models.SysPost) error {
	return r.db.Save(post).Error
}

func (r *SysPostRepositoryImpl) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除岗位
		if err := tx.Delete(&models.SysPost{}, id).Error; err != nil {
			return err
		}
		// 删除用户岗位关联
		if err := tx.Where("post_id = ?", id).Delete(&models.SysUserPost{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SysPostRepositoryImpl) UpdateStatus(id uint64, status models.Status) error {
	return r.db.Model(&models.SysPost{}).Where("id = ?", id).Update("status", status).Error
}
