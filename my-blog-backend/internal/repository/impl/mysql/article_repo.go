package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// ArticleRepositoryImpl 文章仓储实现
type ArticleRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleRepositoryImpl 创建文章仓储实例
func NewArticleRepositoryImpl(db *gorm.DB) repository.ArticleRepository {
	return &ArticleRepositoryImpl{db: db}
}

// Create 创建文章
func (r *ArticleRepositoryImpl) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *ArticleRepositoryImpl) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *ArticleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error

}

// GetByID 根据ID获取文章
func (r *ArticleRepositoryImpl) GetByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// List 分页获取文章列表
func (r *ArticleRepositoryImpl) List(page, pageSize int, status *int, categoryID *uint) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Article{}).
		Select("id, title, slug, summary, cover, category_id, author_id, views, likes, favorites, comment_count, status, is_top, sort_order, published_at, created_at, updated_at")

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error

	return articles, total, err
}

// GetByCategoryID 根据分类ID获取文章列表
func (r *ArticleRepositoryImpl) GetByCategoryID(categoryID uint, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Article{}).
		Select("id, title, slug, summary, cover, category_id, author_id, views, likes, favorites, comment_count, status, is_top, sort_order, published_at, created_at, updated_at").
		Where("category_id = ?", categoryID)
	query.Count(&total)
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error

	return articles, total, err
}

// GetByTagID 根据标签ID获取文章列表
func (r *ArticleRepositoryImpl) GetByTagID(tagID uint, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Article{}).
		Select("articles.id, articles.title, articles.slug, articles.summary, articles.cover, articles.category_id, articles.author_id, articles.views, articles.likes, articles.favorites, articles.comment_count, articles.status, articles.is_top, articles.sort_order, articles.published_at, articles.created_at, articles.updated_at").
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", tagID)
	query.Count(&total)
	err := query.Offset(offset).Limit(pageSize).Order("articles.created_at DESC").Find(&articles).Error

	return articles, total, err
}

// Search 搜索文章
func (r *ArticleRepositoryImpl) Search(keyword string, page, pageSize int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Article{}).
		Select("id, title, slug, summary, cover, category_id, author_id, views, likes, favorites, comment_count, status, is_top, sort_order, published_at, created_at, updated_at").
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	query.Count(&total)
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error

	return articles, total, err
}

// IncrementViewCount 增加浏览次数
func (r *ArticleRepositoryImpl) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *ArticleRepositoryImpl) IncrementLikeCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (r *ArticleRepositoryImpl) DecrementLikeCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes - 1")).Error
}

// GetHotArticles 获取热门文章
func (r *ArticleRepositoryImpl) GetHotArticles(limit int) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Select("id, title, slug, summary, cover, category_id, author_id, views, likes, favorites, comment_count, status, is_top, sort_order, published_at, created_at, updated_at").Order("views DESC").Limit(limit).Find(&articles).Error
	return articles, err
}

// GetRecentArticles 获取最新文章
func (r *ArticleRepositoryImpl) GetRecentArticles(limit int) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Select("id, title, slug, summary, cover, category_id, author_id, views, likes, favorites, comment_count, status, is_top, sort_order, published_at, created_at, updated_at").Order("created_at DESC").Limit(limit).Find(&articles).Error
	return articles, err
}
