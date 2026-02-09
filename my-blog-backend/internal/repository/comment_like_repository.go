package repository

import models "my-blog-backend/internal/models/frontendModel"

// CommentLikeRepository 评论点赞仓储接口
type CommentLikeRepository interface {
	// Create 创建点赞记录
	Create(like *models.CommentLike) error

	// Delete 删除点赞记录
	Delete(commentID, userID uint) error

	// CheckUserLiked 检查用户是否已点赞某评论
	CheckUserLiked(commentID, userID uint) (bool, error)

	// GetLikeCount 获取评论的点赞数
	GetLikeCount(commentID uint) (int64, error)

	// DeleteByCommentID 删除评论的所有点赞记录（用于删除评论时清理）
	DeleteByCommentID(commentID uint) error

	// DeleteByUserID 删除用户的所有点赞记录（用于删除用户时清理）
	DeleteByUserID(userID uint) error

	// GetCommentIDsByUserID 获取用户点赞的评论ID列表
	GetCommentIDsByUserID(userID uint, page, pageSize int) ([]uint, int64, error)
}
