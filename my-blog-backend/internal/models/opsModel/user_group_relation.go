package models

import "time"

// UserGroupRelation 用户组关联表（用户组 <-> 用户）
type UserGroupRelation struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserGroupID uint      `gorm:"type:uint;not null;comment:用户组ID"`
	UserID      uint      `gorm:"type:uint;not null;comment:用户ID"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (UserGroupRelation) TableName() string {
	return "user_group_relations"
}
