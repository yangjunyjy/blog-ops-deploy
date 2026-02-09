package response

// TagItem 标签项
type TagItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int64  `json:"articleCount"`
}

// TagDetailResponse 标签详情响应
type TagDetailResponse struct {
	Tag      TagInfo              `json:"tag"`
	Articles *ArticleListResponse `json:"articles"`
}
