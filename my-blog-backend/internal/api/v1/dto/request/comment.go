package request

// CommentListRequest 获取文章评论请求
type CommentListRequest struct {
	ArticleID uint `uri:"articleId" binding:"required"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ArticleID uint   `json:"articleId" binding:"required"`
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	ParentID  *uint  `json:"parentId" binding:"omitempty"`
}

// DeleteCommentRequest 删除评论请求
type DeleteCommentRequest struct {
	ID uint `uri:"id" binding:"required"`
}
