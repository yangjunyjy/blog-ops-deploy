package claims

import (
	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 用户JWT声明
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     uint8  `json:"role"` // 0: 访客, 1: 普通用户, 2: 管理员
	jwt.RegisteredClaims
}

func ParseToken(token string, secret string) (*UserClaims, error) {
	var claims UserClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
