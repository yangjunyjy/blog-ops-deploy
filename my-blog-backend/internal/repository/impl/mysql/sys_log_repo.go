package mysql

import (
	"gorm.io/gorm"
	"time"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

type SysLogRepositoryImpl struct {
	db *gorm.DB
}

func NewSysLogRepositoryImpl(db *gorm.DB) repository.SysLogRepository {
	return &SysLogRepositoryImpl{db: db}
}

// 操作日志
func (r *SysLogRepositoryImpl) CreateOperationLog(log *models.SysOperationLog) error {
	return r.db.Create(log).Error
}

func (r *SysLogRepositoryImpl) ListOperationLogs(page, pageSize int) ([]*models.SysOperationLog, int64, error) {
	var logs []*models.SysOperationLog
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&models.SysOperationLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *SysLogRepositoryImpl) DeleteOperationLogs(ids []uint64) error {
	return r.db.Delete(&models.SysOperationLog{}, ids).Error
}

// 登录日志
func (r *SysLogRepositoryImpl) CreateLoginLog(log *models.SysLoginLog) error {
	return r.db.Create(log).Error
}

func (r *SysLogRepositoryImpl) ListLoginLogs(page, pageSize int) ([]*models.SysLoginLog, int64, error) {
	var logs []*models.SysLoginLog
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&models.SysLoginLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *SysLogRepositoryImpl) DeleteLoginLogs(ids []uint64) error {
	return r.db.Delete(&models.SysLoginLog{}, ids).Error
}

// 清理日志
func (r *SysLogRepositoryImpl) ClearOldOperationLogs(days int) error {
	beforeDate := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", beforeDate).Delete(&models.SysOperationLog{}).Error
}

func (r *SysLogRepositoryImpl) ClearOldLoginLogs(days int) error {
	beforeDate := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", beforeDate).Delete(&models.SysLoginLog{}).Error
}
