package repository

import models "my-blog-backend/internal/models/frontendModel"

// ArticleLikeRepository 文章点赞仓储接口
type ArticleLikeRepository interface {
	// Create 创建点赞记录
	Create(like *models.ArticleLike) error
	// Delete 删除点赞记录
	Delete(articleID, userID uint) error
	// CheckUserLiked 检查用户是否已点赞某文章
	CheckUserLiked(articleID, userID uint) (bool, error)
	// GetLikeCount 获取文章点赞数
	GetLikeCount(articleID uint) (int64, error)
}
