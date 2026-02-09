package request

type RoleSearchQueryRequest struct {
	RoleName string `form:"roleName"`
	RoleCode string `form:"roleCode"`
	Status   int    `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type RoleAssignMenusRequest struct {
	RoleID  uint64   `json:"role_id" binding:"required"`
	MenuIDs []uint64 `json:"menu_ids" binding:"required"`
}

type CreateRoleRequest struct {
	RoleCode string `json:"role_code" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
	RoleDesc string `json:"role_desc"`
	Status   int8   `json:"status" binding:"required"`
}

type UpdateRoleRequest struct {
	ID       uint   `json:"id" binding:"required"`
	RoleCode string `json:"role_code" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
	RoleDesc string `json:"role_desc"`
	Status   int8   `json:"status" binding:"required"`
}
