package mysql

import (
	opsModel "my-blog-backend/internal/models/opsModel"
	"my-blog-backend/internal/repository"

	"gorm.io/gorm"
)

type HostRepository struct {
	db *gorm.DB
}

func NewHostRepository(db *gorm.DB) repository.HostRepository {
	return &HostRepository{db: db}
}

func (r *HostRepository) Create(host *opsModel.RemoteHost) error {
	return r.db.Create(host).Error
}

func (r *HostRepository) Update(host *opsModel.RemoteHost) error {
	return r.db.Save(host).Error
}

func (r *HostRepository) Delete(id uint) error {
	return r.db.Delete(&opsModel.RemoteHost{}, id).Error
}

func (r *HostRepository) GetByID(id uint) (*opsModel.RemoteHost, error) {
	var host opsModel.RemoteHost
	err := r.db.First(&host, id).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *HostRepository) List(page, pageSize int, name, address, hostType, status string) ([]*opsModel.RemoteHost, int64, error) {
	var hosts []*opsModel.RemoteHost
	var total int64

	query := r.db.Model(&opsModel.RemoteHost{})

	// 添加过滤条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if address != "" {
		query = query.Where("address LIKE ?", "%"+address+"%")
	}
	if hostType != "" {
		var sshType opsModel.SshType
		switch hostType {
		case "password":
			sshType = opsModel.Pwd
		case "key":
			sshType = opsModel.Key
		}
		query = query.Where("type = ?", sshType)
	}
	var st int8 = 1
	if status != "active" {
		st = 0
	}
	if status != "" {
		query = query.Where("status = ?", st)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&hosts).Error; err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

func (r *HostRepository) GetAll() ([]*opsModel.RemoteHost, error) {
	var hosts []*opsModel.RemoteHost
	err := r.db.Order("id DESC").Find(&hosts).Error
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
