package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// ArticleLikeRepositoryImpl 文章点赞仓储实现
type ArticleLikeRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleLikeRepositoryImpl 创建文章点赞仓储实例
func NewArticleLikeRepositoryImpl(db *gorm.DB) repository.ArticleLikeRepository {
	return &ArticleLikeRepositoryImpl{db: db}
}

// Create 创建点赞记录
func (r *ArticleLikeRepositoryImpl) Create(like *models.ArticleLike) error {
	// 检查是否已点赞
	var count int64
	if err := r.db.Model(&models.ArticleLike{}).
		Where("article_id = ? AND user_id = ?", like.ArticleID, like.UserID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	return r.db.Create(like).Error
}

// Delete 删除点赞记录
func (r *ArticleLikeRepositoryImpl) Delete(articleID, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&models.ArticleLike{}).Error
}

// CheckUserLiked 检查用户是否已点赞某文章
func (r *ArticleLikeRepositoryImpl) CheckUserLiked(articleID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ArticleLike{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetLikeCount 获取文章点赞数
func (r *ArticleLikeRepositoryImpl) GetLikeCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleLike{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
