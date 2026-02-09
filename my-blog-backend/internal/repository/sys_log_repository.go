package repository

import (
	"my-blog-backend/internal/models"
)

// SysLogRepository 系统日志仓储接口
type SysLogRepository interface {
	// 操作日志
	CreateOperationLog(log *models.SysOperationLog) error
	ListOperationLogs(page, pageSize int) ([]*models.SysOperationLog, int64, error)
	DeleteOperationLogs(ids []uint64) error

	// 登录日志
	CreateLoginLog(log *models.SysLoginLog) error
	ListLoginLogs(page, pageSize int) ([]*models.SysLoginLog, int64, error)
	DeleteLoginLogs(ids []uint64) error

	// 清理日志
	ClearOldOperationLogs(days int) error
	ClearOldLoginLogs(days int) error
}
