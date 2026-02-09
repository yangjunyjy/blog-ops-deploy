package models

import (
	"time"

	"my-blog-backend/internal/models"
)

// UserGroup 用户组表
type UserGroup struct {
	ID        uint          `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name      string        `gorm:"type:varchar(100);not null;uniqueIndex;comment:用户组名称"`
	Desc      string        `gorm:"type:varchar(255);comment:描述"`
	Sort      int           `gorm:"type:int;default:0;comment:排序"`
	Status    models.Status `gorm:"type:tinyint(1);not null;default:1;comment:状态(0:禁用,1:启用)"`
	CreatedBy uint          `gorm:"type:uint;not null;comment:创建人ID"`
	CreatedAt time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 设置表名
func (UserGroup) TableName() string {
	return "user_groups"
}
