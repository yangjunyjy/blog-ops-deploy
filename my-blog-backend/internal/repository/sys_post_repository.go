package repository

import (
	"my-blog-backend/internal/models"
)

// SysPostRepository 系统岗位仓储接口
type SysPostRepository interface {
	// FindByID 根据ID查找岗位
	FindByID(id uint64) (*models.SysPost, error)
	// FindByName 根据岗位名称查找岗位
	FindByName(name string) (*models.SysPost, error)
	// FindByCode 根据岗位编码查找岗位
	FindByCode(code string) (*models.SysPost, error)
	// List 查询岗位列表
	List(page, pageSize int) ([]*models.SysPost, int64, error)
	// ListAll 查询所有岗位
	ListAll() ([]*models.SysPost, error)
	// Create 创建岗位
	Create(post *models.SysPost) error
	// Update 更新岗位
	Update(post *models.SysPost) error
	// Delete 删除岗位
	Delete(id uint64) error
	// UpdateStatus 更新岗位状态
	UpdateStatus(id uint64, status models.Status) error
}
