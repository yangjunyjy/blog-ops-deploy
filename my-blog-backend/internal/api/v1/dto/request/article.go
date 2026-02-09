package request

// ArticleListRequest 获取文章列表请求
type ArticleListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize   int    `form:"pageSize" binding:"omitempty,min=1,max=50" default:"10"`
	CategoryID uint   `form:"categoryId" binding:"omitempty"`
	Tag        string `form:"tag" binding:"omitempty"`
	Keyword    string `form:"keyword" binding:"omitempty"`
	Sort       string `form:"sort" binding:"omitempty,oneof=latest hot"`
	Status     uint8  `form:"status" binding:"omitempty,oneof=0 1"`
}

// LikeArticleRequest 点赞文章请求
type LikeArticleRequest struct {
	ArticleID uint `uri:"id" binding:"required"`
}

// FavoriteArticleRequest 收藏文章请求
type FavoriteArticleRequest struct {
	ArticleID uint `uri:"id" binding:"required"`
}
