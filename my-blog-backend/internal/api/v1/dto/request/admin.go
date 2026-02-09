package request

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title       string  `json:"title" binding:"required,max=200"`
	Summary     string  `json:"summary" binding:"max=500"`
	Content     string  `json:"content" binding:"required"`
	Cover       *string `json:"cover" binding:"omitempty,url"`
	CategoryID  *uint64 `json:"categoryId" binding:"omitempty"`
	TagIDs      []uint  `json:"tags" binding:"omitempty"`
	AuthorID    uint64  `json:"author_id" binding:"omitempty"`
	Status      uint8   `json:"status" binding:"omitempty,oneof=0 1" default:"1"`
	PublishedAt *string `json:"publishedAt" binding:"omitempty"`
	IsTop       bool    `json:"is_top" binding:"omitempty"`
	Slug        string  `json:"slug" binding:"required,max=200"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	CreateArticleRequest
}

// ArticleIDRequest 文章ID请求
type ArticleIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Slug        string  `json:"slug" binding:"required,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Icon        *string `json:"icon" binding:"omitempty"`
}

// UpdateCommentStatusRequest 更新评论状态请求
type UpdateCommentStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=1 2"`
}

// CommentIDRequest 评论ID请求
type CommentIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}
