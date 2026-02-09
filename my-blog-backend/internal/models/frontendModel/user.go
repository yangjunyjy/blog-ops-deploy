package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型，存储系统用户的基本信息和账户状态
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`                                                   // 用户唯一标识ID
	Username  string         `gorm:"size:50;not null;index:idx_username" json:"username" validate:"required,min=3,max=50"` // 用户登录用户名，全局唯一
	Email     string         `gorm:"size:100;not null;index:idx_email" json:"email" validate:"required,email"`             // 用户邮箱地址，全局唯一，可用于找回密码
	Password  string         `gorm:"size:255;not null" json:"-" validate:"required,min=6"`                                 // 用户密码哈希值，使用bcrypt加密，-表示不序列化到JSON
	Nickname  string         `gorm:"size:50" json:"nickname"`                                                              // 用户昵称，用于展示
	Avatar    string         `gorm:"size:500" json:"avatar"`                                                               // 用户头像URL地址
	Bio       string         `gorm:"type:text" json:"bio"`                                                                 // 用户个人简介
	Role      uint8          `gorm:"not null;default:0" json:"role"`                                                       // 用户角色权限：0-普通用户，1-管理员
	Status    uint8          `gorm:"not null;default:1" json:"status"`                                                     // 账户状态：0-禁用，1-启用
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`                                                     // 账户创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                                                     // 账户最后更新时间
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`                                                        // 软删除时间，用于记录用户何时被删除

	// 关联关系（不使用外键约束）
	Comments []Comment `gorm:"-" json:"comments,omitempty"` // 用户发表的评论列表，通过手动查询获取
}

func (User) TableName() string {
	return "user"
}
