package models

import "time"

// HostUserPermission 主机用户权限关联表（用户组 <-> 主机组）
// 实现RBAC：用户组对主机组拥有访问权限
type HostUserPermission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserGroupID uint      `gorm:"type:uint;not null;comment:用户组ID"`
	HostGroupID uint      `gorm:"type:uint;not null;comment:主机组ID"`
	CreatedBy   uint      `gorm:"type:uint;not null;comment:创建人ID"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (HostUserPermission) TableName() string {
	return "host_user_permissions"
}
