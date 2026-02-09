package response

// CategoryItem 分类项
type CategoryItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"articleCount"`
}

// CategoryDetailResponse 分类详情响应
type CategoryDetailResponse struct {
	Category CategoryInfo         `json:"category"`
	Articles *ArticleListResponse `json:"articles"`
}
