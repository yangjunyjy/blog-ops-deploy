package services

import (
	"errors"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/repository"
)

// SysDeptService 部门服务
type SysDeptService interface {
	// List 列表
	List(page, pageSize int) ([]*models.SysDept, int64, error)
	// GetTree 获取部门树
	GetTree() ([]*models.SysDept, error)
	// GetByID 获取详情
	GetByID(id uint64) (*models.SysDept, error)
	// Create 创建
	Create(dept *models.SysDept) error
	// Update 更新
	Update(dept *models.SysDept) error
	// Delete 删除
	Delete(id uint64) error
}

type sysDeptService struct {
	deptRepo repository.SysDeptRepository
}

func NewSysDeptService(deptRepo repository.SysDeptRepository) SysDeptService {
	return &sysDeptService{deptRepo: deptRepo}
}

func (s *sysDeptService) List(page, pageSize int) ([]*models.SysDept, int64, error) {
	return s.deptRepo.List(page, pageSize)
}

func (s *sysDeptService) GetTree() ([]*models.SysDept, error) {
	return s.deptRepo.GetDeptTree()
}

func (s *sysDeptService) GetByID(id uint64) (*models.SysDept, error) {
	return s.deptRepo.FindByID(id)
}

func (s *sysDeptService) Create(dept *models.SysDept) error {
	// 检查部门名称是否已存在
	if _, err := s.deptRepo.FindByName(dept.Name); err == nil {
		return errors.New("部门名称已存在")
	}
	return s.deptRepo.Create(dept)
}

func (s *sysDeptService) Update(dept *models.SysDept) error {
	// 检查部门名称是否被其他部门使用
	existing, err := s.deptRepo.FindByName(dept.Name)
	if err == nil && existing.ID != dept.ID {
		return errors.New("部门名称已存在")
	}
	return s.deptRepo.Update(dept)
}

func (s *sysDeptService) Delete(id uint64) error {
	return s.deptRepo.Delete(id)
}
