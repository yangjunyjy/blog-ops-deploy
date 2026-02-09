package models

import (
	"time"

	"my-blog-backend/internal/models"
)

type SshType uint

const (
	Key SshType = iota // 0: 密钥登录
	Pwd                // 1: 密码登录
)

type RemoteHost struct {
	ID        uint          `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name      string        `gorm:"type:varchar(100);not null;comment:主机名称"`
	Username  string        `gorm:"type:varchar(50);not null;comment:用户名"`
	Password  string        `gorm:"type:varchar(255);comment:密码（加密存储）"`
	SecretKey string        `gorm:"type:text;comment:私钥内容"`
	Port      int64         `gorm:"type:int;not null;default:22;comment:SSH端口"`
	Address   string        `gorm:"type:varchar(100);not null;comment:主机地址或IP"`
	Type      SshType       `gorm:"type:tinyint(1);not null;default:1;comment:登录类型(0:密钥,1:密码)"`
	Status    models.Status `gorm:"type:tinyint(1);not null;default:1;comment:状态(0:禁用,1:启用)"`
	CreatedAt time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 设置表名
func (RemoteHost) TableName() string {
	return "remote_hosts"
}

// 枚举值的字符串表示（可选）
func (t SshType) String() string {
	switch t {
	case Key:
		return "密钥登录"
	case Pwd:
		return "密码登录"
	default:
		return "未知"
	}
}

// 状态枚举方法（如果需要）
func (t SshType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}
