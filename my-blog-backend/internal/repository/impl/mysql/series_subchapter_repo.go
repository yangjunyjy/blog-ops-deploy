package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// SeriesSubchapterRepositoryImpl 系列子章节仓储实现
type SeriesSubchapterRepositoryImpl struct {
	db *gorm.DB
}

// NewSeriesSubchapterRepositoryImpl 创建系列子章节仓储实例
func NewSeriesSubchapterRepositoryImpl(db *gorm.DB) repository.SeriesSubchapterRepository {
	return &SeriesSubchapterRepositoryImpl{db: db}
}

// Create 创建子章节
func (r *SeriesSubchapterRepositoryImpl) Create(subchapter *models.SeriesSubchapter) error {
	return r.db.Create(subchapter).Error
}

// Update 更新子章节
func (r *SeriesSubchapterRepositoryImpl) Update(subchapter *models.SeriesSubchapter) error {
	return r.db.Save(subchapter).Error
}

// Delete 删除子章节
func (r *SeriesSubchapterRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.SeriesSubchapter{}, id).Error
}

// GetByID 根据ID获取子章节
func (r *SeriesSubchapterRepositoryImpl) GetByID(id uint) (*models.SeriesSubchapter, error) {
	var subchapter models.SeriesSubchapter
	err := r.db.Preload("Articles").First(&subchapter, id).Error
	if err != nil {
		return nil, err
	}
	return &subchapter, nil
}

// GetBySectionID 根据章节ID获取子章节列表
func (r *SeriesSubchapterRepositoryImpl) GetBySectionID(sectionID uint) ([]*models.SeriesSubchapter, int64, error) {
	var subchapters []*models.SeriesSubchapter
	var total int64

	r.db.Model(&models.SeriesSubchapter{}).Where("section_id = ?", sectionID).Count(&total)
	err := r.db.Where("section_id = ?", sectionID).
		Order("sort_order ASC, created_at ASC").
		Find(&subchapters).Error

	return subchapters, total, err
}

// AddArticle 添加文章到子章节
func (r *SeriesSubchapterRepositoryImpl) AddArticle(subchapterID, articleID uint, sortOrder int) error {
	// 检查是否已存在
	var count int64
	r.db.Table("subchapter_articles").
		Where("subchapter_id = ? AND article_id = ?", subchapterID, articleID).
		Count(&count)

	if count > 0 {
		// 已存在，更新排序
		return r.db.Table("subchapter_articles").
			Where("subchapter_id = ? AND article_id = ?", subchapterID, articleID).
			Update("sort_order", sortOrder).Error
	}

	// 不存在，创建新记录
	return r.db.Exec(
		"INSERT INTO subchapter_articles (subchapter_id, article_id, sort_order, created_at) VALUES (?, ?, ?, NOW())",
		subchapterID, articleID, sortOrder,
	).Error
}

// RemoveArticle 从子章节移除文章
func (r *SeriesSubchapterRepositoryImpl) RemoveArticle(subchapterID, articleID uint) error {
	return r.db.Exec(
		"DELETE FROM subchapter_articles WHERE subchapter_id = ? AND article_id = ?",
		subchapterID, articleID,
	).Error
}

// GetArticles 获取子章节的文章列表
func (r *SeriesSubchapterRepositoryImpl) GetArticles(subchapterID uint) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Table("articles").
		Joins("JOIN subchapter_articles ON articles.id = subchapter_articles.article_id").
		Where("subchapter_articles.subchapter_id = ?", subchapterID).
		Order("subchapter_articles.sort_order ASC").
		Find(&articles).Error
	return articles, err
}
