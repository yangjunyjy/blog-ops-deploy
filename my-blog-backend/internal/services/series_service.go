package services

import (
	"errors"
	"fmt"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

var (
	ErrSeriesNotFound        = errors.New("系列不存在")
	ErrSeriesSectionNotFound = errors.New("章节不存在")
)

// SeriesService 系列服务
type SeriesService struct {
	seriesRepo     repository.SeriesRepository
	sectionRepo    repository.SeriesSectionRepository
	subchapterRepo repository.SeriesSubchapterRepository
}

// NewSeriesService 创建系列服务
func NewSeriesService(
	seriesRepo repository.SeriesRepository,
	sectionRepo repository.SeriesSectionRepository,
	subchapterRepo repository.SeriesSubchapterRepository,
) *SeriesService {
	return &SeriesService{
		seriesRepo:     seriesRepo,
		sectionRepo:    sectionRepo,
		subchapterRepo: subchapterRepo,
	}
}

// CreateSeries 创建系列
func (s *SeriesService) CreateSeries(series *models.Series) error {
	return s.seriesRepo.Create(series)
}

// UpdateSeries 更新系列
func (s *SeriesService) UpdateSeries(series *models.Series) error {
	return s.seriesRepo.Update(series)
}

// DeleteSeries 删除系列
func (s *SeriesService) DeleteSeries(id uint) error {
	return s.seriesRepo.Delete(id)
}

// GetSeriesByID 根据ID获取系列
func (s *SeriesService) GetSeriesByID(id uint) (*models.Series, error) {
	series, err := s.seriesRepo.GetByID(id)
	if err != nil {
		return nil, ErrSeriesNotFound
	}
	return series, nil
}

// GetSeriesBySlug 根据slug获取系列
func (s *SeriesService) GetSeriesBySlug(slug string) (*models.Series, error) {
	series, err := s.seriesRepo.GetBySlug(slug)
	if err != nil {
		return nil, ErrSeriesNotFound
	}
	return series, nil
}

// ListSeries 分页获取系列列表
func (s *SeriesService) ListSeries(page, pageSize int) ([]*models.Series, int64, error) {
	return s.seriesRepo.List(page, pageSize)
}

// GetSeriesDetail 获取系列详情（包含章节）
func (s *SeriesService) GetSeriesDetail(id uint) (*models.Series, error) {
	fmt.Printf("🔍 查询系列详情: id=%d\n", id)
	series, err := s.seriesRepo.GetByID(id)
	if err != nil {
		fmt.Printf("❌ 查询系列失败: id=%d, err=%v\n", id, err)
		return nil, ErrSeriesNotFound
	}
	fmt.Printf("✅ 查询系列成功: id=%d, name=%s\n", series.ID, series.Name)

	// 已经在 repository 层通过 Preload 加载了所有关联数据
	return series, nil
}

// ==================== 章节相关 ====================

// CreateSection 创建章节
func (s *SeriesService) CreateSection(section *models.SeriesSection) error {
	// 检查系列是否存在
	_, err := s.seriesRepo.GetByID(uint(section.SeriesID))
	if err != nil {
		return ErrSeriesNotFound
	}
	return s.sectionRepo.Create(section)
}

// UpdateSection 更新章节
func (s *SeriesService) UpdateSection(section *models.SeriesSection) error {
	return s.sectionRepo.Update(section)
}

// DeleteSection 删除章节
func (s *SeriesService) DeleteSection(id uint) error {
	return s.sectionRepo.Delete(id)
}

// GetSectionByID 根据ID获取章节
func (s *SeriesService) GetSectionByID(id uint) (*models.SeriesSection, error) {
	section, err := s.sectionRepo.GetByID(id)
	if err != nil {
		return nil, ErrSeriesSectionNotFound
	}
	return section, nil
}

// ListSectionsBySeriesID 根据系列ID获取章节列表
func (s *SeriesService) ListSectionsBySeriesID(seriesID uint, page, pageSize int) ([]*models.SeriesSection, int64, error) {
	// 检查系列是否存在
	_, err := s.seriesRepo.GetByID(seriesID)
	if err != nil {
		return nil, 0, ErrSeriesNotFound
	}
	sections, total, err := s.sectionRepo.GetBySeriesID(seriesID)
	return sections, total, err
}

// GetSectionDetail 获取章节详情（包含子章节）
func (s *SeriesService) GetSectionDetail(id uint) (*models.SeriesSection, error) {
	section, err := s.sectionRepo.GetByID(id)
	if err != nil {
		return nil, ErrSeriesSectionNotFound
	}

	// 加载子章节 - 已经在repository中通过Preload加载
	return section, nil
}

// ==================== 子章节相关 ====================

// CreateSubchapter 创建子章节
func (s *SeriesService) CreateSubchapter(subchapter *models.SeriesSubchapter) error {
	// 检查章节是否存在
	_, err := s.sectionRepo.GetByID(uint(subchapter.SectionID))
	if err != nil {
		return ErrSeriesSectionNotFound
	}
	return s.subchapterRepo.Create(subchapter)
}

// UpdateSubchapter 更新子章节
func (s *SeriesService) UpdateSubchapter(subchapter *models.SeriesSubchapter) error {
	return s.subchapterRepo.Update(subchapter)
}

// DeleteSubchapter 删除子章节
func (s *SeriesService) DeleteSubchapter(id uint) error {
	return s.subchapterRepo.Delete(id)
}

// GetSubchapterByID 根据ID获取子章节
func (s *SeriesService) GetSubchapterByID(id uint) (*models.SeriesSubchapter, error) {
	subchapter, err := s.subchapterRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("子章节不存在")
	}
	return subchapter, nil
}

// ListSubchaptersBySectionID 根据章节ID获取子章节列表
func (s *SeriesService) ListSubchaptersBySectionID(sectionID uint, page, pageSize int) ([]*models.SeriesSubchapter, int64, error) {
	// 检查章节是否存在
	_, err := s.sectionRepo.GetByID(sectionID)
	if err != nil {
		return nil, 0, ErrSeriesSectionNotFound
	}
	subchapters, total, err := s.subchapterRepo.GetBySectionID(sectionID)
	return subchapters, total, err
}

// ==================== 子章节文章关联 ====================

// AddArticleToSubchapter 添加文章到子章节
func (s *SeriesService) AddArticleToSubchapter(subchapterID, articleID uint, sortOrder int) error {
	return s.subchapterRepo.AddArticle(subchapterID, articleID, sortOrder)
}

// RemoveArticleFromSubchapter 从子章节移除文章
func (s *SeriesService) RemoveArticleFromSubchapter(subchapterID, articleID uint) error {
	return s.subchapterRepo.RemoveArticle(subchapterID, articleID)
}

// GetArticlesBySubchapterID 根据子章节ID获取文章列表
func (s *SeriesService) GetArticlesBySubchapterID(subchapterID uint) ([]*models.Article, error) {
	return s.subchapterRepo.GetArticles(subchapterID)
}
