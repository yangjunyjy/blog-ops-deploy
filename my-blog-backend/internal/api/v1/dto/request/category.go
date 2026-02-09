package request

// CategoryDetailRequest 获取分类详情请求
type CategoryDetailRequest struct {
	ID       uint `uri:"id" binding:"required"`
	Page     int  `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int  `form:"pageSize" binding:"omitempty,min=1" default:"12"`
}
