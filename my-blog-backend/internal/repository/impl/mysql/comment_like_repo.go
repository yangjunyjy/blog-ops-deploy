package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
)

// CommentLikeRepositoryImpl 评论点赞仓储实现
type CommentLikeRepositoryImpl struct {
	db *gorm.DB
}

// NewCommentLikeRepositoryImpl 创建评论点赞仓储实例
func NewCommentLikeRepositoryImpl(db *gorm.DB) repository.CommentLikeRepository {
	return &CommentLikeRepositoryImpl{db: db}
}

// Create 创建点赞记录
func (r *CommentLikeRepositoryImpl) Create(like *models.CommentLike) error {
	return r.db.Create(like).Error
}

// Delete 删除点赞记录
func (r *CommentLikeRepositoryImpl) Delete(commentID, userID uint) error {
	return r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&models.CommentLike{}).Error
}

// CheckUserLiked 检查用户是否已点赞某评论
func (r *CommentLikeRepositoryImpl) CheckUserLiked(commentID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	if count > 0 {
		logger.Info("用户已点赞该评论", logger.Any("comment_id", commentID), logger.Any("user_id", userID))
	}
	return count > 0, err
}

// GetLikeCount 获取评论的点赞数
func (r *CommentLikeRepositoryImpl) GetLikeCount(commentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}

// DeleteByCommentID 删除评论的所有点赞记录
func (r *CommentLikeRepositoryImpl) DeleteByCommentID(commentID uint) error {
	return r.db.Where("comment_id = ?", commentID).Delete(&models.CommentLike{}).Error
}

// DeleteByUserID 删除用户的所有点赞记录
func (r *CommentLikeRepositoryImpl) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CommentLike{}).Error
}

// GetCommentIDsByUserID 获取用户点赞的评论ID列表
func (r *CommentLikeRepositoryImpl) GetCommentIDsByUserID(userID uint, page, pageSize int) ([]uint, int64, error) {
	var commentIDs []uint
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.CommentLike{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Model(&models.CommentLike{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Pluck("comment_id", &commentIDs).Error

	return commentIDs, total, err
}
