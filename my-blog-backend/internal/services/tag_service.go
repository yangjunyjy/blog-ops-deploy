package services

import (
	"errors"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
	"time"
)

type TagService interface {
	CreateTag(tag *models.Tag) error
	UpdateTag(id uint, tag *models.Tag) error
	DeleteTag(id uint) error
	GetTag(id uint) (*models.Tag, error)
	GetTagBySlug(slug string) (*models.Tag, error)
	ListTags(page, pageSize int) ([]*models.Tag, int64, error)
	GetArticleTags(articleID uint) ([]*models.Tag, error)
	GetTagArticles(tagID uint, page, pageSize int) ([]*models.Article, int64, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (s *tagService) CreateTag(tag *models.Tag) error {
	now := time.Now()
	tag.CreatedAt = now
	return s.tagRepo.Create(tag)
}

func (s *tagService) UpdateTag(id uint, tag *models.Tag) error {
	existing, err := s.tagRepo.GetByID(id)
	if err != nil {
		return errors.New("tag not found")
	}

	// Tag的ID是uint64，不需要转换
	tag.CreatedAt = existing.CreatedAt
	return s.tagRepo.Update(tag)
}

func (s *tagService) DeleteTag(id uint) error {
	return s.tagRepo.Delete(id)
}

func (s *tagService) GetTag(id uint) (*models.Tag, error) {
	return s.tagRepo.GetByID(id)
}

func (s *tagService) GetTagBySlug(slug string) (*models.Tag, error) {
	return s.tagRepo.GetBySlug(slug)
}

func (s *tagService) ListTags(page, pageSize int) ([]*models.Tag, int64, error) {
	return s.tagRepo.List(page, pageSize)
}

func (s *tagService) GetArticleTags(articleID uint) ([]*models.Tag, error) {
	return s.tagRepo.GetByArticleID(articleID)
}

func (s *tagService) GetTagArticles(tagID uint, page, pageSize int) ([]*models.Article, int64, error) {
	return s.tagRepo.GetTagArticles(tagID, page, pageSize)
}
