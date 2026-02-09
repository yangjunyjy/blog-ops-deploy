package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// SeriesSectionRepositoryImpl 系列章节仓储实现
type SeriesSectionRepositoryImpl struct {
	db *gorm.DB
}

// NewSeriesSectionRepositoryImpl 创建系列章节仓储实例
func NewSeriesSectionRepositoryImpl(db *gorm.DB) repository.SeriesSectionRepository {
	return &SeriesSectionRepositoryImpl{db: db}
}

// Create 创建章节
func (r *SeriesSectionRepositoryImpl) Create(section *models.SeriesSection) error {
	return r.db.Create(section).Error
}

// Update 更新章节
func (r *SeriesSectionRepositoryImpl) Update(section *models.SeriesSection) error {
	return r.db.Save(section).Error
}

// Delete 删除章节
func (r *SeriesSectionRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.SeriesSection{}, id).Error
}

// GetByID 根据ID获取章节
func (r *SeriesSectionRepositoryImpl) GetByID(id uint) (*models.SeriesSection, error) {
	var section models.SeriesSection
	err := r.db.Preload("Subchapters").First(&section, id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

// GetBySeriesID 根据系列ID获取章节列表
func (r *SeriesSectionRepositoryImpl) GetBySeriesID(seriesID uint) ([]*models.SeriesSection, int64, error) {
	var sections []*models.SeriesSection
	var total int64

	r.db.Model(&models.SeriesSection{}).Where("series_id = ?", seriesID).Count(&total)
	err := r.db.Where("series_id = ?", seriesID).
		Order("sort_order ASC, created_at ASC").
		Find(&sections).Error

	return sections, total, err
}
