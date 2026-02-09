package repository

import (
	"my-blog-backend/internal/models"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
	GetByID(id uint) (*models.Role, error)
	GetByCode(code string) (*models.Role, error)
	List(page, pageSize int) ([]*models.Role, int64, error)
}
