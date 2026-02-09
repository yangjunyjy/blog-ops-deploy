package models

import (
	"time"

	"my-blog-backend/internal/models"
)

type AccountType uint

const (
	RootAccount  AccountType = iota + 1 // 1: root 账号
	NormalAccount                      // 2: 普通账号
)

// HostAccount 主机账号表
type HostAccount struct {
	ID          uint          `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name        string        `gorm:"type:varchar(50);not null;comment:账号名称"`
	Username    string        `gorm:"type:varchar(50);not null;comment:用户名"`
	Password    string        `gorm:"type:varchar(255);comment:密码（加密存储）"`
	SecretKey   string        `gorm:"type:text;comment:私钥内容"`
	Type        AccountType    `gorm:"type:tinyint(1);not null;default:2;comment:账号类型(1:root,2:普通)"`
	HostID      uint          `gorm:"type:uint;not null;comment:关联主机ID"`
	Status      models.Status `gorm:"type:tinyint(1);not null;default:1;comment:状态(0:禁用,1:启用)"`
	Remark      string        `gorm:"type:varchar(255);comment:备注"`
	CreatedBy   uint          `gorm:"type:uint;not null;comment:创建人ID"`
	CreatedAt   time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt   time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 设置表名
func (HostAccount) TableName() string {
	return "host_accounts"
}

// 枚举值的字符串表示
func (t AccountType) String() string {
	switch t {
	case RootAccount:
		return "root"
	case NormalAccount:
		return "普通账号"
	default:
		return "未知"
	}
}

// MarshalJSON 自定义JSON序列化
func (t AccountType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}
