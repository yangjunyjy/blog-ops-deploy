package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/infrastructure/smtp"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
)

const (
	sysCaptchaExpireDuration = 5 * time.Minute
	sysCaptchaPrefix         = "sysCaptcha:"
	sysRateEmailPrefix       = "rate:sys:email:"
)

// SysUserService 系统用户服务
type SysUserService interface {
	// Login 用户登录（用户名+密码）
	Login(username, password string) (*models.SysUser, string, error)
	// GenerateEmailCode 生成邮箱验证码
	GenerateEmailCode(ctx context.Context, email, username string) error
	// VerifyEmailCode 验证邮箱验证码
	VerifyEmailCode(ctx context.Context, email, code string, clear bool) (bool, error)
	// Logout 用户登出
	Logout(userID uint64) error
	// GetUserInfo 获取用户信息
	GetUserInfo(userID uint64) (*models.SysUser, error)
	// GetUserInfoByEmail 根据邮箱获取用户信息
	GetUserInfoByEmail(email string) (*models.SysUser, error)
	// CreateUser 创建用户
	CreateUser(user *models.SysUser, creatorID uint64) error
	// UpdateUser 更新用户
	UpdateUser(user *models.SysUser, updaterID uint64) error
	// DeleteUser 删除用户
	DeleteUser(id, operatorID uint64) error
	// ResetPassword 重置密码
	ResetPassword(id uint64, newPassword string) error
	// AssignRoles 分配角色
	AssignRoles(userID uint64, roleIDs []uint64) error
	// AssignPosts 分配岗位
	AssignPosts(userID uint64, postIDs []uint64) error
	// List 获取用户列表（分页）
	List(query *request.SearchSysUserQueryRequest) ([]*response.UserListResponse, int64, error)
}

type sysUserService struct {
	userRepo    repository.SysUserRepository
	logRepo     repository.SysLogRepository
	cacheRepo   repository.CacheRepo
	emailServer *smtp.Sender
}

func NewSysUserService(
	userRepo repository.SysUserRepository,
	logRepo repository.SysLogRepository,
	cacheRepo repository.CacheRepo,
) SysUserService {
	return &sysUserService{
		userRepo:    userRepo,
		logRepo:     logRepo,
		cacheRepo:   cacheRepo,
		emailServer: smtp.GetEmailServer(),
	}
}

func (s *sysUserService) Login(username, password string) (*models.SysUser, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, "", errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginTime = &now
	s.userRepo.Update(user)

	// 记录登录日志
	s.logRepo.CreateLoginLog(&models.SysLoginLog{
		Username:     username,
		LoginIP:      "",
		Browser:      "",
		OS:           "",
		LoginStatus:  0,
		LoginMessage: "登录成功",
		LoginTime:    time.Now(),
	})

	return user, "", nil
}

func (s *sysUserService) Logout(userID uint64) error {
	// TODO: 清除 Session
	return nil
}

func (s *sysUserService) GetUserInfo(userID uint64) (*models.SysUser, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 加载角色和岗位信息
	roles, _ := s.userRepo.GetUserRoles(userID)
	// logger.Info("角色信息", logger.Any("roles", roles))
	posts, _ := s.userRepo.GetUserPosts(userID)
	user.Roles = roles
	user.Posts = posts

	return user, nil
}

func (s *sysUserService) CreateUser(user *models.SysUser, creatorID uint64) error {
	// 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(user.Username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.FindByEmail(user.Email); err == nil {
		return errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(user)
}

func (s *sysUserService) UpdateUser(user *models.SysUser, updaterID uint64) error {
	// 如果提供了新密码，则加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		// 如果没有提供新密码，查询旧密码并保留
		oldUser, err := s.userRepo.FindByID(user.ID)
		if err != nil {
			return err
		}
		user.Password = oldUser.Password
	}
	return s.userRepo.Update(user)
}

func (s *sysUserService) DeleteUser(id, operatorID uint64) error {
	// 不能删除自己
	if id == operatorID {
		return errors.New("不能删除自己")
	}

	return s.userRepo.Delete(id)
}

func (s *sysUserService) ResetPassword(id uint64, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.ResetPassword(id, string(hashedPassword))
}

func (s *sysUserService) AssignRoles(userID uint64, roleIDs []uint64) error {
	return s.userRepo.AssignRoles(userID, roleIDs)
}

func (s *sysUserService) AssignPosts(userID uint64, postIDs []uint64) error {
	return s.userRepo.AssignPosts(userID, postIDs)
}

// List 获取用户列表（分页）
func (s *sysUserService) List(query *request.SearchSysUserQueryRequest) ([]*response.UserListResponse, int64, error) {
	return s.userRepo.List(query)
}

// GetUserInfoByEmail 根据邮箱获取用户信息
func (s *sysUserService) GetUserInfoByEmail(email string) (*models.SysUser, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	// 加载角色和岗位信息
	roles, _ := s.userRepo.GetUserRoles(user.ID)
	posts, _ := s.userRepo.GetUserPosts(user.ID)
	user.Roles = roles
	user.Posts = posts

	return user, nil
}

// GenerateEmailCode 生成邮箱验证码
func (s *sysUserService) GenerateEmailCode(ctx context.Context, email, username string) error {
	logger.Info("GenerateEmailCode 开始",
		logger.String("email", email),
		logger.String("username", username))

	if strings.TrimSpace(email) == "" {
		logger.Error("GenerateEmailCode: 邮箱为空")
		return fmt.Errorf("邮箱不能为空")
	}

	// 限流检查
	if !s.allowRequest(ctx, sysRateEmailPrefix+email, 1, 1*time.Minute) {
		logger.Error("GenerateEmailCode: 请求过于频繁", logger.String("email", email))
		return fmt.Errorf("该邮箱请求过于频繁，请稍后再试")
	}
	logger.Info("GenerateEmailCode: 限流检查通过")

	// 生成并发送验证码
	code, err := s.emailServer.SendVerificationCode(email, username)
	if err != nil {
		logger.Error("GenerateEmailCode: 发送验证码失败", logger.Err("error", err))
		return err
	}
	logger.Info("GenerateEmailCode: 验证码已生成并发送", logger.String("code", code))

	// 存储验证码到 Redis
	cacheKey := sysCaptchaPrefix + email
	logger.Info("GenerateEmailCode: 存储验证码到Redis",
		logger.String("cache_key", cacheKey),
		logger.String("code", code),
		logger.String("duration", sysCaptchaExpireDuration.String()))
	if err := s.cacheRepo.SetCaptcha(ctx, cacheKey, code, sysCaptchaExpireDuration); err != nil {
		logger.Error("GenerateEmailCode: 存储验证码到Redis失败", logger.Err("error", err))
		return fmt.Errorf("验证码生成失败")
	}
	logger.Info("GenerateEmailCode: 验证码存储成功")

	return nil
}

// VerifyEmailCode 验证邮箱验证码
func (s *sysUserService) VerifyEmailCode(ctx context.Context, email, code string, clear bool) (bool, error) {
	logger.Info("VerifyEmailCode 开始",
		logger.String("email", email),
		logger.String("code", code),
		logger.String("prefix", sysCaptchaPrefix))

	if strings.TrimSpace(email) == "" || strings.TrimSpace(code) == "" {
		logger.Error("VerifyEmailCode: 邮箱或验证码为空")
		return false, nil
	}
	logger.Info("VerifyEmailCode: 参数校验通过")

	// 从 Redis 获取存储的验证码
	cacheKey := sysCaptchaPrefix + email
	logger.Info("VerifyEmailCode: 从Redis获取验证码", logger.String("cache_key", cacheKey))
	storedCode, err := s.cacheRepo.GetCaptcha(ctx, cacheKey)
	if err != nil {
		logger.Error("VerifyEmailCode: 获取验证码失败", logger.Err("error", err))
		return false, fmt.Errorf("获取验证码失败: %v", err)
	}
	logger.Info("VerifyEmailCode: Redis中存储的验证码", logger.String("stored_code", storedCode))

	// 验证码匹配
	if storedCode == code {
		logger.Info("VerifyEmailCode: 验证码匹配成功")
		if clear {
			logger.Info("VerifyEmailCode: 清除验证码")
			s.cacheRepo.DelCaptcha(ctx, cacheKey)
		}
		return true, nil
	}
	logger.Error("VerifyEmailCode: 验证码不匹配",
		logger.String("stored_code", storedCode),
		logger.String("input_code", code))
	return false, nil
}

// allowRequest 限流检查
func (s *sysUserService) allowRequest(ctx context.Context, key string, limit int, window time.Duration) bool {
	countStr, _ := s.cacheRepo.GetCapetcha(ctx, key)
	if countStr == "" {
		// 初始化计数
		_ = s.cacheRepo.SetCaptcha(ctx, key, "1", window)
		return true
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		_ = s.cacheRepo.SetCaptcha(ctx, key, "1", window)
		return true
	}

	if count >= limit {
		return false
	}

	// 自增计数，保持 TTL
	_ = s.cacheRepo.SetCaptcha(ctx, key, strconv.Itoa(count+1), window)
	return true
}
