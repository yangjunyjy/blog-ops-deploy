package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Create(article *models.Article) error
	Update(article *models.Article) error
	Delete(id uint) error
	GetByID(id uint) (*models.Article, error)
	List(page, pageSize int, status *int, categoryID *uint) ([]*models.Article, int64, error)
	Search(keyword string, page, pageSize int) ([]*models.Article, int64, error)
	IncrementViewCount(id uint) error
	IncrementLikeCount(id uint) error
	DecrementLikeCount(id uint) error
	GetHotArticles(limit int) ([]*models.Article, error)
	GetRecentArticles(limit int) ([]*models.Article, error)
}
