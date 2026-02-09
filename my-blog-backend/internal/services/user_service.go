package services

import (
	"errors"
	"fmt"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *models.User) error
	UpdateUser(id uint, user *models.User) error
	DeleteUser(id uint) error
	GetUser(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	ListUsers(page, pageSize int) ([]*models.User, int64, error)
	ValidatePassword(username, password string) (*models.User, error)
	ChangePassword(id uint, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(user *models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Status = 1

	return s.userRepo.Create(user)
}

func (s *userService) UpdateUser(id uint, user *models.User) error {
	existing, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.ID = id
	user.CreatedAt = existing.CreatedAt
	user.UpdatedAt = time.Now()

	// 如果密码不为空，则加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		user.Password = existing.Password
	}

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) GetUser(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) ListUsers(page, pageSize int) ([]*models.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

func (s *userService) ValidatePassword(username, password string) (*models.User, error) {
	fmt.Printf("🔍 验证密码: 用户名=%s\n", username)

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		fmt.Printf("❌ 用户不存在: %v\n", err)
		return nil, errors.New("invalid username or password")
	}

	fmt.Printf("✅ 找到用户: ID=%d, 用户名=%s, 密码哈希=%s\n", user.ID, user.Username, user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("❌ 密码不匹配: %v\n", err)
		return nil, errors.New("invalid username or password")
	}

	fmt.Printf("✅ 密码验证成功\n")
	return user, nil
}

func (s *userService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}
