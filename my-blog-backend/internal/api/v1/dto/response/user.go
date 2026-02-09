package response

import "time"

// FavoriteArticleItem 收藏文章项
type FavoriteArticleItem struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Summary   string     `json:"summary"`
	Cover     *string    `json:"cover"`
	Author    AuthorInfo `json:"author"`
	CreatedAt time.Time  `json:"createdAt"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	List     []FavoriteArticleItem `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// FolderWithCount 带文章数量的文件夹
type FolderWithCount struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	SortOrder     int       `json:"sortOrder"`
	ArticleCount  int64     `json:"articleCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// FolderListResponse 文件夹列表响应
type FolderListResponse struct {
	List     []FolderWithCount `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// LikeListResponse 点赞列表响应
type LikeListResponse struct {
	List     []FavoriteArticleItem `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// UserCommentItem 用户评论项
type UserCommentItem struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	Article struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"article"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserCommentListResponse 用户评论列表响应
type UserCommentListResponse struct {
	List     []UserCommentItem `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// CheckFavoriteResponse 检查收藏响应
type CheckFavoriteResponse struct {
	IsFavorite bool `json:"isFavorite"`
}

// CheckLikeResponse 检查点赞响应
type CheckLikeResponse struct {
	IsLiked bool `json:"isLiked"`
}
