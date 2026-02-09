package models

import (
	"time"

)

type ScheduleType uint

const (
	OnceSchedule   ScheduleType = 1 // 执行一次
	DailySchedule  ScheduleType = 2 // 每天
	WeeklySchedule ScheduleType = 3 // 每周
	MonthlySchedule ScheduleType = 4 // 每月
	CronSchedule   ScheduleType = 5 // Cron 表达式
)

type ScheduleStatus uint

const (
	ScheduleActive   ScheduleStatus = 1 // 激活
	SchedulePaused  ScheduleStatus = 2 // 暂停
	ScheduleExpired ScheduleStatus = 3 // 已过期
)

// SchedulePlan 定时计划表
type SchedulePlan struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name           string         `gorm:"type:varchar(100);not null;comment:计划名称"`
	Type           ScheduleType   `gorm:"type:tinyint(1);not null;comment:计划类型(1:一次,2:每天,3:每周,4:每月,5:Cron)"`
	Status         ScheduleStatus `gorm:"type:tinyint(1);not null;default:1;index;comment:状态(1:激活,2:暂停,3:已过期)"`
	CronExpression string         `gorm:"type:varchar(100);comment:Cron表达式"`
	ExecuteDate   string         `gorm:"type:date;comment:执行日期(一次性执行)"`
	ExecuteTime   string         `gorm:"type:time;comment:执行时间"`
	WeekDays      string         `gorm:"type:varchar(20);comment:周几(1-7,逗号分隔)"`
	MonthDay      int            `gorm:"type:int;comment:每月几号(1-31)"`
	TaskType      uint           `gorm:"type:tinyint(1);not null;comment:关联任务类型(1:命令,2:文件上传,3:文件下载,4:脚本)"`
	Command       string         `gorm:"type:text;comment:命令或脚本内容"`
	SourcePath    string         `gorm:"type:varchar(500);comment:源路径"`
	TargetPath    string         `gorm:"type:varchar(500);comment:目标路径"`
	ScriptType    string         `gorm:"type:varchar(20);comment:脚本类型"`
	Timeout       int            `gorm:"type:int;default:300;comment:超时时间(秒)"`
	Remark        string         `gorm:"type:varchar(255);comment:备注"`
	CreatedBy     uint           `gorm:"type:uint;not null;comment:创建人ID"`
	CreatedAt     time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index;comment:创建时间"`
	UpdatedAt     time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	LastExecutedAt *time.Time    `gorm:"type:datetime;comment:最后执行时间"`
	NextExecutedAt *time.Time    `gorm:"type:datetime;index;comment:下次执行时间"`
}

// TableName 设置表名
func (SchedulePlan) TableName() string {
	return "schedule_plans"
}

// ScheduleHostRelation 定时计划主机关联表
type ScheduleHostRelation struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	SchedulePlanID uint      `gorm:"type:uint;not null;index;comment:定时计划ID"`
	HostID         uint      `gorm:"type:uint;not null;comment:主机ID"`
	AccountID      uint      `gorm:"type:uint;comment:使用的账号ID"`
	CreatedAt      time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (ScheduleHostRelation) TableName() string {
	return "schedule_host_relations"
}

// 枚举值的字符串表示
func (t ScheduleType) String() string {
	switch t {
	case OnceSchedule:
		return "执行一次"
	case DailySchedule:
		return "每天"
	case WeeklySchedule:
		return "每周"
	case MonthlySchedule:
		return "每月"
	case CronSchedule:
		return "Cron表达式"
	default:
		return "未知"
	}
}

func (s ScheduleStatus) String() string {
	switch s {
	case ScheduleActive:
		return "激活"
	case SchedulePaused:
		return "暂停"
	case ScheduleExpired:
		return "已过期"
	default:
		return "未知"
	}
}
