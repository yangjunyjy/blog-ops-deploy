package response

import (
	"time"
)

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TagInfo 标签信息
type TagInfo struct {
	Name         string `json:"name"`
	ArticleCount int64  `json:"articleCount"`
}

// AuthorInfo 作者信息
type AuthorInfo struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
	Bio    *string `json:"bio"`
}

// ArticleListItem 文章列表项
type ArticleListItem struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Summary   string       `json:"summary"`
	Cover     *string      `json:"cover"`
	Category  CategoryInfo `json:"category"`
	Tags      []TagInfo    `json:"tags"`
	Author    AuthorInfo   `json:"author"`
	Views     int64        `json:"views"`
	Likes     int64        `json:"likes"`
	Favorites int64        `json:"favorites"`
	Comments  int64        `json:"comments"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	List     []ArticleListItem `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// ArticleDetail 文章详情
type ArticleDetail struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Content     string       `json:"content"`
	Cover       *string      `json:"cover"`
	Category    CategoryInfo `json:"category"`
	Tags        []TagInfo    `json:"tags"`
	Author      AuthorInfo   `json:"author"`
	Views       int64        `json:"views"`
	Likes       int64        `json:"likes"`
	Favorites   int64        `json:"favorites"`
	Comments    int64        `json:"comments"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	PrevArticle *ArticleNav  `json:"prevArticle,omitempty"`
	NextArticle *ArticleNav  `json:"nextArticle,omitempty"`
}

// ArticleNav 文章导航（上一篇/下一篇）
type ArticleNav struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// HotArticleItem 热门文章项
type HotArticleItem struct {
	ID    uint    `json:"id"`
	Title string  `json:"title"`
	Views int64   `json:"views"`
	Cover *string `json:"cover"`
}

// RecentArticleItem 最新文章项
type RecentArticleItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Cover     *string   `json:"cover"`
}

// LikeResponse 点赞响应
type LikeResponse struct {
	Likes int64 `json:"likes"`
}

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	Favorites int64 `json:"favorites"`
}
