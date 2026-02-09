package request

// TagDetailRequest 获取标签详情请求
type TagDetailRequest struct {
	Name     string `uri:"name" binding:"required"`
	Page     int    `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1" default:"12"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Slug        string  `json:"slug" binding:"required,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}
