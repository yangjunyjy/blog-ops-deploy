package request

// SearchRequest 搜索文章请求
type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int    `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1" default:"12"`
}
