package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id uint) error
	GetByID(id uint) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	List(page, pageSize int) ([]*models.Category, int64, error)
	GetTree() ([]*models.Category, error)
	GetChildren(parentID uint) ([]*models.Category, error)
	GetCategoryArticles(categoryID uint, page, pageSize int) ([]*models.Article, int64, error)
}
