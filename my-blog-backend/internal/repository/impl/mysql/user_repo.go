package mysql

import (
	"gorm.io/gorm"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepositoryImpl 创建用户仓储实例
func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 分页获取用户列表
func (r *UserRepositoryImpl) List(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	offset := (page - 1) * pageSize
	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}
