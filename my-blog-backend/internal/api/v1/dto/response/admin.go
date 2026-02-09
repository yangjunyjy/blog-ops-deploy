package response

import "time"

// CreateArticleResponse 创建文章响应
type CreateArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}

// UpdateArticleResponse 更新文章响应
type UpdateArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateCategoryResponse 创建分类响应
type CreateCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreateTagResponse 创建标签响应
type CreateTagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
