package response

// SearchResultItem 搜索结果项
type SearchResultItem struct {
	ID      uint    `json:"id"`
	Title   string  `json:"title"`
	Summary string  `json:"summary"`
	Cover   *string `json:"cover"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	List     []SearchResultItem `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}
