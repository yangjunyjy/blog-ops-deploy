package services

import (
	"errors"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
	"time"
)

type CategoryService interface {
	CreateCategory(category *models.Category) error
	UpdateCategory(id uint, category *models.Category) error
	DeleteCategory(id uint) error
	GetCategory(id uint) (*models.Category, error)
	GetCategoryBySlug(slug string) (*models.Category, error)
	ListCategories(page, pageSize int) ([]*models.Category, int64, error)
	GetCategoryTree() ([]*models.Category, error)
	GetCategoryChildren(parentID uint) ([]*models.Category, error)
	GetCategoryArticles(categoryID uint, page, pageSize int) ([]*models.Article, int64, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(category *models.Category) error {
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	return s.categoryRepo.Create(category)
}

func (s *categoryService) UpdateCategory(id uint, category *models.Category) error {
	existing, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	category.ID = id
	category.CreatedAt = existing.CreatedAt
	category.UpdatedAt = time.Now()
	return s.categoryRepo.Update(category)
}

func (s *categoryService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) GetCategory(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	return s.categoryRepo.GetBySlug(slug)
}

func (s *categoryService) ListCategories(page, pageSize int) ([]*models.Category, int64, error) {
	return s.categoryRepo.List(page, pageSize)
}

func (s *categoryService) GetCategoryTree() ([]*models.Category, error) {
	return s.categoryRepo.GetTree()
}

func (s *categoryService) GetCategoryChildren(parentID uint) ([]*models.Category, error) {
	return s.categoryRepo.GetChildren(parentID)
}

func (s *categoryService) GetCategoryArticles(categoryID uint, page, pageSize int) ([]*models.Article, int64, error) {
	return s.categoryRepo.GetCategoryArticles(categoryID, page, pageSize)
}
