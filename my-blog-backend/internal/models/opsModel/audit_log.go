package models

import "time"

type AuditStatus uint

const (
	AuditSuccess AuditStatus = iota + 1 // 1: 成功
	AuditFailed                      // 2: 失败
	AuditWarning                     // 3: 警告
)

type AuditAction uint

const (
	LoginAction   AuditAction = iota + 1 // 1: 登录
	ExecuteAction                       // 2: 执行命令
	FileUploadAction                    // 3: 文件上传
	FileDownloadAction                  // 4: 文件下载
	SessionAction                      // 5: 会话管理
)

type RiskLevel uint

const (
	LowRisk    RiskLevel = 1 // 低风险
	MediumRisk RiskLevel = 2 // 中风险
	HighRisk   RiskLevel = 3 // 高风险
	CriticalRisk RiskLevel = 4 // 严重风险
)

// AuditLog 用户审计表（审计分析用户的命令和异常行为）
type AuditLog struct {
	ID           uint        `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID       uint        `gorm:"type:uint;not null;index:idx_user_time;comment:用户ID"`
	UserName     string      `gorm:"type:varchar(50);not null;comment:用户名"`
	HostID       uint        `gorm:"type:uint;not null;index;comment:主机ID"`
	HostName     string      `gorm:"type:varchar(100);not null;comment:主机名称"`
	HostAddress   string      `gorm:"type:varchar(100);not null;comment:主机地址"`
	SessionID    string      `gorm:"type:varchar(100);not null;index;comment:会话ID"`
	Action       AuditAction `gorm:"type:tinyint(1);not null;comment:操作类型(1:登录,2:执行命令,3:文件上传,4:文件下载,5:会话管理)"`
	Command      string      `gorm:"type:text;comment:执行的命令"`
	Status       AuditStatus `gorm:"type:tinyint(1);not null;index:idx_user_time;comment:状态(1:成功,2:失败,3:警告)"`
	RiskLevel    RiskLevel   `gorm:"type:tinyint(1);not null;index;comment:风险等级(1:低,2:中,3:高,4:严重)"`
	ClientIP     string      `gorm:"type:varchar(50);comment:客户端IP"`
	ClientAgent  string      `gorm:"type:varchar(255);comment:客户端User-Agent"`
	ErrorMessage string      `gorm:"type:text;comment:错误信息"`
	Duration     int64       `gorm:"type:bigint;comment:执行时长(毫秒)"`
	StartTime    time.Time   `gorm:"type:datetime;not null;index:idx_user_time;comment:开始时间"`
	EndTime      *time.Time  `gorm:"type:datetime;comment:结束时间"`
	CreatedAt    time.Time   `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 设置表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// 枚举值的字符串表示
func (s AuditStatus) String() string {
	switch s {
	case AuditSuccess:
		return "成功"
	case AuditFailed:
		return "失败"
	case AuditWarning:
		return "警告"
	default:
		return "未知"
	}
}

func (a AuditAction) String() string {
	switch a {
	case LoginAction:
		return "登录"
	case ExecuteAction:
		return "执行命令"
	case FileUploadAction:
		return "文件上传"
	case FileDownloadAction:
		return "文件下载"
	case SessionAction:
		return "会话管理"
	default:
		return "未知"
	}
}

func (r RiskLevel) String() string {
	switch r {
	case LowRisk:
		return "低风险"
	case MediumRisk:
		return "中风险"
	case HighRisk:
		return "高风险"
	case CriticalRisk:
		return "严重风险"
	default:
		return "未知"
	}
}
