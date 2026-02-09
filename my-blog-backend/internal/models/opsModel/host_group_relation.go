package models

import "time"

// HostGroupRelation 主机组关联表（主机组 <-> 主机）
type HostGroupRelation struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	HostGroupID uint      `gorm:"type:uint;not null;comment:主机组ID"`
	HostID      uint      `gorm:"type:uint;not null;comment:主机ID"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (HostGroupRelation) TableName() string {
	return "host_group_relations"
}
