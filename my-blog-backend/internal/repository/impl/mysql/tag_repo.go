package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// TagRepositoryImpl 标签仓储实现
type TagRepositoryImpl struct {
	db *gorm.DB
}

// NewTagRepositoryImpl 创建标签仓储实例
func NewTagRepositoryImpl(db *gorm.DB) repository.TagRepository {
	return &TagRepositoryImpl{db: db}
}

// Create 创建标签
func (r *TagRepositoryImpl) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *TagRepositoryImpl) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *TagRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

// GetByID 根据ID获取标签
func (r *TagRepositoryImpl) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetBySlug 根据slug获取标签
func (r *TagRepositoryImpl) GetBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// List 分页获取标签列表
func (r *TagRepositoryImpl) List(page, pageSize int) ([]*models.Tag, int64, error) {
	var tags []*models.Tag
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.Tag{}).Count(&total)
	err := r.db.Offset(offset).Limit(pageSize).Find(&tags).Error

	// 为每个标签填充文章数量
	for _, tag := range tags {
		var count int64
		r.db.Table("article_tag").
			Joins("JOIN articles ON article_tag.article_id = articles.id").
			Where("article_tag.tag_id = ? AND articles.status = ?", tag.ID, 1).
			Count(&count)
		tag.ArticleCount = count
	}

	return tags, total, err
}

// GetByArticleID 根据文章ID获取标签列表
func (r *TagRepositoryImpl) GetByArticleID(articleID uint) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Joins("JOIN article_tag ON tag.id = article_tag.tag_id").
		Where("article_tag.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// GetTagArticles 获取标签下的文章
func (r *TagRepositoryImpl) GetTagArticles(tagID uint, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize

	// 获取文章列表和总数，排除 content 字段以提高查询性能
	query := r.db.Model(&models.Article{}).
		Select("articles.id, articles.title, articles.slug, articles.summary, articles.cover, articles.category_id, articles.author_id, articles.views, articles.likes, articles.favorites, articles.comment_count, articles.status, articles.is_top, articles.sort_order, articles.published_at, articles.created_at, articles.updated_at").
		Joins("JOIN article_tag ON articles.id = article_tag.article_id").
		Where("article_tag.tag_id = ? AND articles.status = ?", tagID, 1)

	query.Count(&total)
	err := query.Preload("Tags").
		Order("articles.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}
