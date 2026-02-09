package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int8

type MenuType int8

// BaseModel 基础模型
type BaseModel struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CreateTime *time.Time `gorm:"autoCreateTime;column:create_time;comment:创建时间" json:"createTime"`
	UpdateTime *time.Time `gorm:"autoUpdateTime;column:update_time;comment:更新时间" json:"updateTime"`
	CreateBy   uint64     `gorm:"column:create_by;comment:创建人" json:"createBy"`
	UpdateBy   uint64     `gorm:"column:update_by;comment:更新人" json:"updateBy"`
	Remark     string     `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`
	Deleted    int8       `gorm:"column:deleted;default:0;comment:删除标记 0:正常 1:删除" json:"deleted"`
}

// BeforeCreate 创建前钩子
func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreateTime = &now
	m.UpdateTime = &now
	return
}

// BeforeUpdate 更新前钩子
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateTime = &now
	return
}

// DeleteModel 删除标记模型
type DeleteModel struct {
	Deleted int8 `gorm:"column:deleted;default:0;comment:删除标记 0:正常 1:删除" json:"deleted"`
}

// StatusModel 状态模型
type StatusModel struct {
	Status int8 `gorm:"column:status;default:1;comment:状态 0:禁用 1:启用" json:"status"`
}

// SortModel 排序模型
type SortModel struct {
	Sort int `gorm:"column:sort;default:0;comment:排序" json:"sort"`
}

// ScopeDeleted 未删除的作用域
func ScopeDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted = ?", 0)
}

// ScopeEnabled 启用状态的作用域
func ScopeEnabled(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", 1)
}
