package models

import (
	"time"
)

type TaskStatus uint

const (
	TaskPending  TaskStatus = 1 // 待执行
	TaskRunning  TaskStatus = 2 // 执行中
	TaskSuccess  TaskStatus = 3 // 成功
	TaskFailed   TaskStatus = 4 // 失败
	TaskCanceled TaskStatus = 5 // 已取消
)

type TaskType uint

const (
	CommandTask      TaskType = 1 // 命令任务
	FileUploadTask   TaskType = 2 // 文件上传
	FileDownloadTask TaskType = 3 // 文件下载
	ScriptTask       TaskType = 4 // 脚本任务
)

// BatchTask 任务表（用于执行多个主机批量化任务）
type BatchTask struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name         string     `gorm:"type:varchar(100);not null;comment:任务名称"`
	Type         TaskType   `gorm:"type:tinyint(1);not null;comment:任务类型(1:命令,2:文件上传,3:文件下载,4:脚本)"`
	Status       TaskStatus `gorm:"type:tinyint(1);not null;default:1;index;comment:状态(1:待执行,2:执行中,3:成功,4:失败,5:已取消)"`
	Command      string     `gorm:"type:text;comment:执行的命令或脚本内容"`
	SourcePath   string     `gorm:"type:varchar(500);comment:源文件路径（文件上传/下载）"`
	TargetPath   string     `gorm:"type:varchar(500);comment:目标路径"`
	ScriptType   string     `gorm:"type:varchar(20);comment:脚本类型(如:bash,python,shell)"`
	Timeout      int        `gorm:"type:int;default:300;comment:超时时间(秒)"`
	SuccessCount int        `gorm:"type:int;default:0;comment:成功数量"`
	FailedCount  int        `gorm:"type:int;default:0;comment:失败数量"`
	TotalHosts   int        `gorm:"type:int;default:0;comment:总主机数"`
	Progress     float64    `gorm:"type:decimal(5,2);default:0;comment:进度(0-100)"`
	Remark       string     `gorm:"type:varchar(255);comment:备注"`
	CreatedBy    uint       `gorm:"type:uint;not null;comment:创建人ID"`
	CreatedAt    time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index;comment:创建时间"`
	StartedAt    *time.Time `gorm:"type:datetime;comment:开始时间"`
	FinishedAt   *time.Time `gorm:"type:datetime;comment:完成时间"`
}

// TableName 设置表名
func (BatchTask) TableName() string {
	return "batch_tasks"
}

// 枚举值的字符串表示
func (s TaskStatus) String() string {
	switch s {
	case TaskPending:
		return "待执行"
	case TaskRunning:
		return "执行中"
	case TaskSuccess:
		return "成功"
	case TaskFailed:
		return "失败"
	case TaskCanceled:
		return "已取消"
	default:
		return "未知"
	}
}

func (t TaskType) String() string {
	switch t {
	case CommandTask:
		return "命令任务"
	case FileUploadTask:
		return "文件上传"
	case FileDownloadTask:
		return "文件下载"
	case ScriptTask:
		return "脚本任务"
	default:
		return "未知"
	}
}

// TaskHostRelation 任务主机关联表（任务 <-> 主机）
type TaskHostRelation struct {
	ID         uint       `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	TaskID     uint       `gorm:"type:uint;not null;index;comment:任务ID"`
	HostID     uint       `gorm:"type:uint;not null;comment:主机ID"`
	HostName   string     `gorm:"type:varchar(100);not null;comment:主机名称"`
	HostAddr   string     `gorm:"type:varchar(100);comment:主机地址"`
	AccountID  uint       `gorm:"type:uint;comment:使用的账号ID"`
	Status     TaskStatus `gorm:"type:tinyint(1);not null;default:1;comment:执行状态"`
	Output     string     `gorm:"type:text;comment:执行输出"`
	Error      string     `gorm:"type:text;comment:错误信息"`
	Duration   int64      `gorm:"type:bigint;comment:执行时长(毫秒)"`
	StartedAt  *time.Time `gorm:"type:datetime;comment:开始执行时间"`
	FinishedAt *time.Time `gorm:"type:datetime;comment:完成时间"`
	CreatedAt  time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (TaskHostRelation) TableName() string {
	return "task_host_relations"
}
