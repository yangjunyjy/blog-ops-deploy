package services

import (
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

// SysLogService 日志服务
type SysLogService interface {
	// 操作日志
	LogOperation(module, business, method, request, params, ip, location, agent, operator string,
		status int, errorMsg string, costTime int64) error
	ListOperationLogs(page, pageSize int) ([]*models.SysOperationLog, int64, error)
	DeleteOperationLogs(ids []uint64) error

	// 登录日志
	ListLoginLogs(page, pageSize int) ([]*models.SysLoginLog, int64, error)
	DeleteLoginLogs(ids []uint64) error

	// 清理日志
	ClearOldLogs(days int) error
}

type sysLogService struct {
	logRepo repository.SysLogRepository
}

func NewSysLogService(logRepo repository.SysLogRepository) SysLogService {
	return &sysLogService{logRepo: logRepo}
}

func (s *sysLogService) LogOperation(module, business, method, request, params, ip, location, agent, operator string,
	status int, errorMsg string, costTime int64) error {

	log := &models.SysOperationLog{
		Module:   module,
		Business: business,
		Method:   method,
		Request:  request,
		Params:   params,
		IP:       ip,
		Location: location,
		Agent:    agent,
		Operator: operator,
		Status:   status,
		ErrorMsg: errorMsg,
		CostTime: costTime,
	}

	return s.logRepo.CreateOperationLog(log)
}

func (s *sysLogService) ListOperationLogs(page, pageSize int) ([]*models.SysOperationLog, int64, error) {
	return s.logRepo.ListOperationLogs(page, pageSize)
}

func (s *sysLogService) DeleteOperationLogs(ids []uint64) error {
	return s.logRepo.DeleteOperationLogs(ids)
}

func (s *sysLogService) ListLoginLogs(page, pageSize int) ([]*models.SysLoginLog, int64, error) {
	return s.logRepo.ListLoginLogs(page, pageSize)
}

func (s *sysLogService) DeleteLoginLogs(ids []uint64) error {
	return s.logRepo.DeleteLoginLogs(ids)
}

func (s *sysLogService) ClearOldLogs(days int) error {
	s.logRepo.ClearOldOperationLogs(days)
	s.logRepo.ClearOldLoginLogs(days)
	return nil
}
