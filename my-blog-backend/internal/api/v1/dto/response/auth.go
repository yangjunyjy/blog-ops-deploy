package response

import "time"

// UserProfile 用户信息
type UserProfile struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    *string   `json:"avatar"`
	Bio       *string   `json:"bio"`
	Website   *string   `json:"website"`
	Github    *string   `json:"github"`
	Role      int       `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User  UserProfile `json:"user"`
	Token string      `json:"token"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID string `json:"captchaId"`
	Image     string `json:"image"` // base64 编码的图片
}
