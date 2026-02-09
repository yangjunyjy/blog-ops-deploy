package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Code     string `json:"code" binding:"required,omitempty,max=50"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CaptchaID     string `json:"captchaId" binding:"required"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required,len=4"`
}

type LoginByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,omitempty,max=50"`
}

// CaptchaRequest 验证码请求
type CaptchaRequest struct {
	Type int `form:"type" binding:"omitempty,min=1,max=2"` // 1: 数字验证码, 2: 字母数字验证码
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=50"`
	Bio      *string `json:"bio" binding:"omitempty,max=500"`
	Website  *string `json:"website" binding:"omitempty,url"`
	Github   *string `json:"github" binding:"omitempty,url"`
}

// ChangePasswordRequest 更改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=32"`
}

// UserListRequest 用户列表请求（分页）
type UserListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=50" default:"12"`
}

// CheckFavoriteRequest 检查收藏请求
type CheckFavoriteRequest struct {
	ArticleID uint `uri:"articleId" binding:"required"`
}

// CheckLikeRequest 检查点赞请求
type CheckLikeRequest struct {
	ArticleID uint `uri:"articleId" binding:"required"`
}

// EmailCaptchaRequest 邮箱验证码请求
type EmailCaptchaRequest struct {
	Username string `json:"username" binding:"omitempty" example:"admin"` // 用户名可选，如果未提供则使用邮箱
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"`
}

type VerifyCaptchaRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
