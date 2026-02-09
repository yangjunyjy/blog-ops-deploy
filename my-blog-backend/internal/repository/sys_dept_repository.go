package repository

import (
	"my-blog-backend/internal/models"
)

// SysDeptRepository 系统部门仓储接口
type SysDeptRepository interface {
	// FindByID 根据ID查找部门
	FindByID(id uint64) (*models.SysDept, error)
	// FindByName 根据部门名称查找部门
	FindByName(name string) (*models.SysDept, error)
	// List 查询部门列表
	List(page, pageSize int) ([]*models.SysDept, int64, error)
	// GetDeptTree 获取部门树
	GetDeptTree() ([]*models.SysDept, error)
	// Create 创建部门
	Create(dept *models.SysDept) error
	// Update 更新部门
	Update(dept *models.SysDept) error
	// Delete 删除部门
	Delete(id uint64) error
	// GetChildren 获取子部门
	GetChildren(parentID uint64) ([]*models.SysDept, error)
	// GetAncestors 获取祖级列表
	GetAncestors(id uint64) (string, error)
	// UpdateStatus 更新部门状态
	UpdateStatus(id uint64, status models.Status) error
}
