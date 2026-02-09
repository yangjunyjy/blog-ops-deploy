package models

import (
	"time"
)

// SysLoginLog 登录日志表
type SysLoginLog struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username      string `gorm:"column:username;type:varchar(50);comment:用户名" json:"username"`
	LoginIP       string `gorm:"column:login_ip;type:varchar(50);comment:登录IP" json:"login_ip"`
	LoginLocation string `gorm:"column:login_location;type:varchar(255);comment:登录地点" json:"login_location"`
	Browser       string `gorm:"column:browser;type:varchar(100);comment:浏览器类型" json:"browser"`
	OS            string `gorm:"column:os;type:varchar(100);comment:操作系统" json:"os"`
	LoginStatus   int8   `gorm:"column:login_status;type:tinyint;default:0;comment:登录状态 0:成功 1:失败" json:"login_status"`
	LoginMessage  string `gorm:"column:login_message;type:varchar(500);comment:登录消息" json:"login_message"`
	LoginTime     time.Time `gorm:"column:login_time;type:datetime;comment:登录时间" json:"login_time"`
}

// TableName 指定表名
func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
