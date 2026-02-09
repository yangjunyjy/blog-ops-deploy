package repository

import models "my-blog-backend/internal/models/frontendModel"

// SeriesRepository 系列仓储接口
type SeriesRepository interface {
	Create(series *models.Series) error
	Update(series *models.Series) error
	Delete(id uint) error
	GetByID(id uint) (*models.Series, error)
	GetBySlug(slug string) (*models.Series, error)
	List(page, pageSize int) ([]*models.Series, int64, error)
}

// SeriesSectionRepository 系列章节仓储接口
type SeriesSectionRepository interface {
	Create(section *models.SeriesSection) error
	Update(section *models.SeriesSection) error
	Delete(id uint) error
	GetByID(id uint) (*models.SeriesSection, error)
	GetBySeriesID(seriesID uint) ([]*models.SeriesSection, int64, error)
}

// SeriesSubchapterRepository 系列子章节仓储接口
type SeriesSubchapterRepository interface {
	Create(subchapter *models.SeriesSubchapter) error
	Update(subchapter *models.SeriesSubchapter) error
	Delete(id uint) error
	GetByID(id uint) (*models.SeriesSubchapter, error)
	GetBySectionID(sectionID uint) ([]*models.SeriesSubchapter, int64, error)
	AddArticle(subchapterID, articleID uint, sortOrder int) error
	RemoveArticle(subchapterID, articleID uint) error
	GetArticles(subchapterID uint) ([]*models.Article, error)
}
