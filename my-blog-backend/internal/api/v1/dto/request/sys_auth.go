package request

// SysLoginRequest 管理后台登录请求
type SysLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"`
	Code     string `json:"code" binding:"required" example:"123456"`
}

// SysSendEmailCodeRequest 发送邮箱验证码请求
type SysSendEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@example.com"`
}
