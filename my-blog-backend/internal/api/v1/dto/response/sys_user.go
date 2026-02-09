package response

import (
	"my-blog-backend/internal/models"
	"time"
)

type UserListResponse struct {
	ID         uint64            `json:"id"`
	Username   string            `json:"username"`
	Password   string            `json:"password"`
	Nickname   string            `json:"nickName"`
	Avatar     string            `json:"avatar"`
	Email      string            `json:"email"`
	Gender     int8              `json:"gender"`
	Status     int8              `json:"status"`
	IsAdmin    int8              `json:"isAdmin"`
	CreateTime time.Time         `json:"createTime"`
	Roles      []*models.SysRole `json:"roles"`
}
