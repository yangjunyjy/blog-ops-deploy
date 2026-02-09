package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// SeriesRepositoryImpl 系列仓储实现
type SeriesRepositoryImpl struct {
	db *gorm.DB
}

// NewSeriesRepositoryImpl 创建系列仓储实例
func NewSeriesRepositoryImpl(db *gorm.DB) repository.SeriesRepository {
	return &SeriesRepositoryImpl{db: db}
}

// Create 创建系列
func (r *SeriesRepositoryImpl) Create(series *models.Series) error {
	return r.db.Create(series).Error
}

// Update 更新系列
func (r *SeriesRepositoryImpl) Update(series *models.Series) error {
	return r.db.Save(series).Error
}

// Delete 删除系列
func (r *SeriesRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Series{}, id).Error
}

// GetByID 根据ID获取系列
func (r *SeriesRepositoryImpl) GetByID(id uint) (*models.Series, error) {
	var series models.Series
	err := r.db.Preload("Sections.Subchapters.Articles").
		Preload("Sections.Subchapters.Articles.Tags").
		First(&series, id).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

// GetBySlug 根据slug获取系列
func (r *SeriesRepositoryImpl) GetBySlug(slug string) (*models.Series, error) {
	var series models.Series
	err := r.db.Where("slug = ?", slug).First(&series).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

// List 分页获取系列列表
func (r *SeriesRepositoryImpl) List(page, pageSize int) ([]*models.Series, int64, error) {
	var series []*models.Series
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.Series{}).Count(&total)
	err := r.db.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&series).Error

	return series, total, err
}
