package request

type CreateHostRequest struct {
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Port      int    `json:"port" binding:"required,min=1,max=65535"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password"`
	SecretKey string `json:"secret_key"`
	Type      string `json:"type" binding:"required,oneof=password key both"` // password, key, both
	Status    string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateHostRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Port      int    `json:"port" binding:"required,min=1,max=65535"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password"`
	SecretKey string `json:"secret_key"`
	Type      string `json:"type" binding:"required,oneof=password key both"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type DeleteHostRequest struct {
	ID uint `json:"id" binding:"required"`
}

type ListHostRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Name     string `form:"name"`
	Address  string `form:"address"`
	Type     string `form:"type" binding:"omitempty,oneof=password key both"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive"`
}

type GetHostRequest struct {
	ID uint `uri:"id" binding:"required"`
}
