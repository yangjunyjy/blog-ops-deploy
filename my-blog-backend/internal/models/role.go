package models

import "time"

// Role 角色模型，定义用户的系统角色和权限
type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                              // 角色唯一标识ID
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`         // 角色名称，如"管理员"、"普通用户"
	Code      string    `gorm:"size:50;uniqueIndex" json:"code"`                 // 角色代码，用于程序判断，如"admin"、"user"
	Status    int       `gorm:"default:1" json:"status"`                          // 角色状态：0-禁用，1-启用
	CreatedAt time.Time `json:"created_at"`                                       // 角色创建时间
	UpdatedAt time.Time `json:"updated_at"`                                       // 角色最后更新时间
}

func (Role) TableName() string {
	return "roles"
}
