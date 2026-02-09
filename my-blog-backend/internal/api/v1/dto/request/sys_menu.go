package request

type CreateSysMenuRequest struct {
	Name        string `json:"name" binding:"required"`   // 必须
	Code        string `json:"code"`                      // 必须
	Type        uint8  `json:"type" binding:"required"`   // 必须
	Status      uint8  `json:"status" binding:"required"` // 必须
	ParentID    uint64 `json:"parent_id"`                 // 可选
	Sort        int    `json:"sort"`                      // 可选
	Component   string `json:"component"`                 // 可选
	Path        string `json:"path"`
	Description string `json:"description"` // 可选
	IsVisible   bool   `json:"isVisible"`   // 可选
	Icon        string `json:"icon"`        // 可选
}

type UpdateSysMenuRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`   // 必须
	Code        string `json:"code"`                      // 必须
	Type        uint8  `json:"type" binding:"required"`   // 必须
	Status      uint8  `json:"status" binding:"required"` // 必须
	ParentID    uint64 `json:"parent_id"`                 // 可选
	Sort        int    `json:"sort"`                      // 可选
	Component   string `json:"component"`                 // 可选
	Path        string `json:"path"`
	Description string `json:"description"` // 可选
	IsVisible   bool   `json:"isVisible"`   // 可选
	Icon        string `json:"icon"`        // 可选
}
