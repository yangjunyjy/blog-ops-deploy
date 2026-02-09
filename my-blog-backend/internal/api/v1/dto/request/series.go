package request

// CreateSeriesRequest 创建系列请求
type CreateSeriesRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Slug        string  `json:"slug" binding:"required,max=100"`
	Icon        *string `json:"icon" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Cover       *string `json:"cover" binding:"omitempty"`
	SortOrder   int     `json:"sortOrder" binding:"omitempty"`
}

// UpdateSeriesRequest 更新系列请求
type UpdateSeriesRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Slug        *string `json:"slug" binding:"omitempty,max=100"`
	Icon        *string `json:"icon" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Cover       *string `json:"cover" binding:"omitempty"`
	SortOrder   *int    `json:"sortOrder" binding:"omitempty"`
	Status      *uint8  `json:"status" binding:"omitempty,oneof=0 1"`
}

// SeriesIDRequest 系列ID请求
type SeriesIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// SeriesListRequest 系列列表请求
type SeriesListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100" default:"10"`
}

// CreateSectionRequest 创建章节请求
type CreateSectionRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   int     `json:"sortOrder" binding:"omitempty"`
}

// UpdateSectionRequest 更新章节请求
type UpdateSectionRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   *int    `json:"sortOrder" binding:"omitempty"`
}

// SectionIDRequest 章节ID请求
type SectionIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// SectionListRequest 章节列表请求
type SectionListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100" default:"10"`
}

// CreateSubchapterRequest 创建子章节请求
type CreateSubchapterRequest struct {
	SectionID   uint64  `json:"sectionId" binding:"required"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   int     `json:"sortOrder" binding:"omitempty"`
}

// UpdateSubchapterRequest 更新子章节请求
type UpdateSubchapterRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SortOrder   *int    `json:"sortOrder" binding:"omitempty"`
}

// SubchapterIDRequest 子章节ID请求
type SubchapterIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// SubchapterListRequest 子章节列表请求
type SubchapterListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100" default:"10"`
}

// AddArticleToSubchapterRequest 添加文章到子章节请求
type AddArticleToSubchapterRequest struct {
	ArticleID uint `json:"articleId" binding:"required"`
	SortOrder int  `json:"sortOrder" binding:"omitempty"`
}

// RemoveArticleFromSubchapterRequest 从子章节移除文章请求
type RemoveArticleFromSubchapterRequest struct {
	SubchapterID uint `uri:"id" binding:"required"`
	ArticleID    uint `uri:"articleId" binding:"required"`
}

// SeriesDetailRequest 获取专栏详情请求
type SeriesDetailRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// SubchapterArticlesRequest 获取子章节文章请求
type SubchapterArticlesRequest struct {
	SeriesID  uint `uri:"seriesId" binding:"required"`
	ChapterID uint `uri:"chapterId" binding:"required"`
}
