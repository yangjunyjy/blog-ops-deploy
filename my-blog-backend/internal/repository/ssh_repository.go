package repository

import (
	models "my-blog-backend/internal/models/opsModel"
)

type HostRepository interface {
	Create(host *models.RemoteHost) error
	Update(host *models.RemoteHost) error
	Delete(id uint) error
	GetByID(id uint) (*models.RemoteHost, error)
	List(page, pageSize int, name, address, hostType, status string) ([]*models.RemoteHost, int64, error)
	GetAll() ([]*models.RemoteHost, error)
}
