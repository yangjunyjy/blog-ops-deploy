package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	List(page, pageSize int) ([]*models.User, int64, error)
}
