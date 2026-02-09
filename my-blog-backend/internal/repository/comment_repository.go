package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
	DeleteWithChildren(id uint) error // 删除评论及其所有子评论
	GetByID(id uint) (*models.Comment, error)
	GetByIDWithUser(id uint) (*models.Comment, error) // 获取评论及用户信息
	GetByArticleID(articleID uint, page, pageSize int) ([]*models.Comment, int64, error)
	List(articleID uint, page, pageSize int) ([]*models.Comment, int64, error)
	ListReplies(parentID uint, page, pageSize int) ([]*models.Comment, int64, error) // 获取子评论列表
	ListByUserID(userID uint, page, pageSize int) ([]*models.Comment, int64, error)
	ListByStatus(status uint8, page, pageSize int) ([]*models.Comment, int64, error)
	ListPendingByUserAndArticle(userID, articleID uint) ([]*models.Comment, error) // 获取某用户在某文章下的待审核评论
	ListAll(page, pageSize int) ([]*models.Comment, int64, error)
	Approve(id uint) error
	Reject(id uint) error
	GetReplyCount(parentID uint) (int64, error) // 获取子评论数量

	// 新增方法：使用中间表
	ListUserCommentedArticles(userID uint, page, pageSize int) ([]*models.Article, int64, error)
	DeleteCommentByArticleID(articleId uint) error
	DecrementArticleCommentCount(articleID uint) error // 减少文章评论数
	IncrementArticleCommentCount(articleID uint) error // 增加文章评论数
}
