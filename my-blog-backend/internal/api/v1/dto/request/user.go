package request

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Password string `json:"password" binding:"omitempty,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Bio      string `json:"bio" binding:"omitempty"`
	Website  string `json:"website" binding:"omitempty"`
	GitHub   string `json:"github" binding:"omitempty"`
	Avatar   string `json:"avatar" binding:"omitempty"`
	Status   *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}
