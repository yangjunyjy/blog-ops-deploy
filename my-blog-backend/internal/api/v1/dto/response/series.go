package response

import "time"

// SubchapterInfo 子章节信息
type SubchapterInfo struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	SortOrder   int               `json:"sortOrder"`
	ArticleIDs  []uint            `json:"articleIds"`
	Articles    []ArticleListItem `json:"articles,omitempty"`
}

// SectionInfo 章节信息
type SectionInfo struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	SortOrder   int              `json:"sortOrder"`
	Subchapters []SubchapterInfo `json:"subchapters"`
}

// SeriesDetail 系列详情
type SeriesDetail struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Icon        string        `json:"icon"`
	Description string        `json:"description"`
	Cover       *string       `json:"cover"`
	SortOrder   int           `json:"sortOrder"`
	Status      uint8         `json:"status"`
	Sections    []SectionInfo `json:"sections"`
	Chapters    []ChapterInfo `json:"chapters,omitempty"` // 兼容前端的 chapters 字段
}

// ChapterInfo 章节信息（用于前端兼容）
type ChapterInfo struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	SortOrder   int               `json:"sortOrder"`
	Articles    []ArticleListItem `json:"articles,omitempty"`
}

// SeriesListItem 系列列表项
type SeriesListItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Cover       *string   `json:"cover"`
	SortOrder   int       `json:"sortOrder"`
	Status      uint8     `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SeriesListResponse 系列列表响应
type SeriesListResponse struct {
	List     []SeriesListItem `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}


