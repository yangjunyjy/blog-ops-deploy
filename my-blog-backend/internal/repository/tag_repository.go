package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// TagRepository 标签仓储接口
type TagRepository interface {
	Create(tag *models.Tag) error
	Update(tag *models.Tag) error
	Delete(id uint) error
	GetByID(id uint) (*models.Tag, error)
	GetBySlug(slug string) (*models.Tag, error)
	List(page, pageSize int) ([]*models.Tag, int64, error)
	GetByArticleID(articleID uint) ([]*models.Tag, error)
	GetTagArticles(tagID uint, page, pageSize int) ([]*models.Article, int64, error)
}
