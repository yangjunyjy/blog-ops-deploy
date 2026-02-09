package models

import "time"

// SysRoleMenu 角色菜单关联表
type SysRoleMenu struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID    uint64    `gorm:"column:role_id;type:bigint;not null;index:idx_role_id,priority:1;comment:角色ID" json:"role_id"`
	MenuID    uint64    `gorm:"column:menu_id;type:bigint;not null;index:idx_menu_id,priority:2;comment:菜单ID" json:"menu_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
}

// TableName 指定表名
func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
