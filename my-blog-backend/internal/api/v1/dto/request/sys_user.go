package request

type CreateSysUserRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   uint8  `json:"status" binding:"required,min=0,max=1"`
	Gender   uint8  `json:"gender"`
}

type UpdateSysUserRequest struct {
	Id       uint64 `json:"id" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Nickname string `json:"nickname"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   uint8  `json:"status" binding:"required,min=0,max=1"`
	Gender   uint8  `json:"gender"`
}

type AssignRolesRequest struct {
	UserID  uint64   `json:"userId" binding:"required"`
	RoleIDs []uint64 `json:"roleIds" binding:"required"`
}

type SearchSysUserQueryRequest struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Status   uint8  `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}
