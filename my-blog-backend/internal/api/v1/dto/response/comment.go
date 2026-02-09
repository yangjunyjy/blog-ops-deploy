package response

import "time"

// CommentAuthor 评论作者信息
type CommentAuthor struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}

// CommentItem 评论项
type CommentItem struct {
	ID        uint          `json:"id"`
	ArticleID uint          `json:"articleId"`
	Content   string        `json:"content"`
	Author    CommentAuthor `json:"author"`
	CreatedAt time.Time     `json:"createdAt"`
	ParentID  *uint         `json:"parentId"`
	Replies   []CommentItem `json:"replies"`
}

// CreateCommentResponse 创建评论响应
type CreateCommentResponse struct {
	ID        uint          `json:"id"`
	ArticleID uint          `json:"articleId"`
	Content   string        `json:"content"`
	Author    CommentAuthor `json:"author"`
	CreatedAt time.Time     `json:"createdAt"`
	ParentID  *uint         `json:"parentId"`
	Replies   []CommentItem `json:"replies"`
}

// DeleteCommentResponse 删除评论响应
type DeleteCommentResponse struct {
	Message string `json:"message"`
}

// CommentListItem 评论列表项（简化版，用于前端展示）
type CommentListItem struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	Author    CommentAuthor `json:"author"`
	CreatedAt time.Time     `json:"createdAt"`
	Likes     int           `json:"likes"` // 点赞数，目前固定为0
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Items    []CommentItemFull `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int64             `json:"total"`
}

// CommentAuthorWithInfo 评论作者信息（包含所有字段）
type CommentAuthorWithInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// CommentItemFull 完整的评论项
type CommentItemFull struct {
	ID         uint                  `json:"id"`
	ArticleID  uint                  `json:"article_id"`
	UserID     uint                  `json:"user_id"`
	ParentID   *uint                 `json:"parent_id"`
	Content    string                `json:"content"`
	Status     uint8                 `json:"status"`
	Likes      uint32                `json:"likes"`
	IsLiked    bool                  `json:"is_liked"`
	IPAddress  string                `json:"ip_address"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	Author     *CommentAuthorWithInfo `json:"author,omitempty"`
	ReplyTo    uint                  `json:"reply_to,omitempty"`       // 被回复的评论ID
	ReplyToUser *CommentAuthorWithInfo `json:"reply_to_user,omitempty"` // 被回复的用户信息
}
