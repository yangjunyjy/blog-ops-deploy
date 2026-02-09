package services

import (
	"context"
	"errors"
	"fmt"
	dtoResponse "my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/infrastructure/smtp"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/claims"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	captchaExpireDuration     = 3 * time.Minute
	redisCaptchaPrefix        = "captcha:"
	smtpCaptchaExpireDuration = 5 * time.Minute
	smtpCaptchaPrefix         = "smtpCaptcha:" // 存储邮箱验证码的前缀
	rateEmailPrefix           = "rate:email:"
	rateIpPrefix              = "rate:ip:"
)

type AuthService interface {
	Login(username, password string) (string, error)
	LoginByEmail(email string, code string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GenerateToken(userID uint, username string, email string, role uint8) (string, error)
	RefreshToken(tokenString string) (string, error)
	GenerateSmtpCaptcha(ctx context.Context, email string, username string) error
	VerifySmtpCaptcha(ctx context.Context, email string, code string, clear bool) (bool, error)
	Register(ctx context.Context, username string, password string, email string) error
	GetUserByUsername(username string) (*dtoResponse.UserProfile, error)
	GetUserByEmail(email string) (*dtoResponse.UserProfile, error)
}

type authService struct {
	userService     UserService
	favoriteService FavoriteService
	cacheRepo       repository.CacheRepo
	jwtSecret       string
	jwtExpire       time.Duration
	logger          logger.Logger
	emailServer     *smtp.Sender
}

func NewAuthService(userService UserService, favoriteService FavoriteService, jwtSecret string, jwtExpire time.Duration) AuthService {
	return &authService{
		userService:     userService,
		favoriteService: favoriteService,
		jwtSecret:       jwtSecret,
		jwtExpire:       jwtExpire,
		logger:          logger.L(),
		emailServer:     smtp.GetEmailServer(),
	}
}

func NewAuthServiceWithCache(userService UserService, favoriteService FavoriteService, cacheRepo repository.CacheRepo, jwtSecret string, jwtExpire time.Duration) AuthService {
	return &authService{
		userService:     userService,
		favoriteService: favoriteService,
		cacheRepo:       cacheRepo,
		jwtSecret:       jwtSecret,
		jwtExpire:       jwtExpire,
		logger:          logger.L(),
		emailServer:     smtp.GetEmailServer(),
	}
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userService.ValidatePassword(username, password)
	if err != nil {
		return "", err
	}

	if user.Status != 1 {
		return "", errors.New("user is disabled")
	}

	return s.GenerateToken(user.ID, user.Username, user.Email, user.Role)
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
}

func (s *authService) GenerateToken(userID uint, username string, email string, role uint8) (string, error) {
	now := time.Now()

	// 如果jwtExpire为0或小于1小时，强制设置为24小时
	expireTime := s.jwtExpire
	if expireTime < time.Hour {
		fmt.Printf("⚠️ JWT过期时间异常: %v, 强制设置为24小时\n", expireTime)
		expireTime = 24 * time.Hour
	}

	userClaims := claims.UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	// 打印token生成信息用于调试
	s.logger.Info("JWT Token生成成功",
		logger.Uint("user_id", userID),
		logger.String("username", username),
		logger.String("email", email),
		logger.Int("role", int(role)),
		logger.Time("issued_at", now),
		logger.Time("expires_at", now.Add(expireTime)),
		logger.Duration("expires_in", expireTime),
		logger.String("token_prefix", tokenString[:min(20, len(tokenString))]+"..."),
	)

	return tokenString, nil
}

func (s *authService) RefreshToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	if userClaims, ok := token.Claims.(*claims.UserClaims); ok && token.Valid {
		return s.GenerateToken(userClaims.UserID, userClaims.Username, userClaims.Email, userClaims.Role)
	}

	return "", errors.New("invalid token")
}

// LoginByEmail 邮箱验证码登录
func (s *authService) LoginByEmail(email string, code string) (string, error) {
	// 根据邮箱查找用户
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		s.logger.Error("邮箱登录失败：用户不存在", logger.Err("error", err))
		return "", errors.New("邮箱未注册")
	}

	if user.Status != 1 {
		return "", errors.New("用户已被禁用")
	}

	// 验证邮箱验证码
	verified, err := s.VerifySmtpCaptcha(context.Background(), email, code, true)
	if err != nil {
		s.logger.Error("验证邮箱验证码失败", logger.Err("error", err))
		return "", errors.New("验证验证码失败")
	}

	if !verified {
		return "", errors.New("验证码错误或已过期")
	}

	// 验证码正确，生成 JWT token
	token, err := s.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		s.logger.Error("生成token失败", logger.Err("error", err))
		return "", err
	}

	s.logger.Info("邮箱登录成功", logger.String("email", email))
	return token, nil
}

func (r *authService) GenerateSmtpCaptcha(ctx context.Context, email string, username string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// --- 限流逻辑 -http.StatusBadRequest, "参数有误"--
	if !r.allowRequest(ctx, rateEmailPrefix+email, 1, 1*time.Minute) {
		return fmt.Errorf("该邮箱请求过于频繁，请稍后再试")
	}
	code, err := r.emailServer.SendVerificationCode(email, username)
	if err != nil {
		return err
	}
	fmt.Println(username, email)
	if err := r.cacheRepo.SetCaptcha(ctx, smtpCaptchaPrefix+email, code, smtpCaptchaExpireDuration); err != nil {
		r.logger.Error("存储邮箱验证码失败", logger.Err("redis", err))
		return fmt.Errorf("验证码生成失败")
	}

	r.logger.Info("邮箱验证码已发送", logger.Any("email", email))
	return nil
}

func (r *authService) allowRequest(ctx context.Context, key string, limit int, window time.Duration) bool {
	countStr, _ := r.cacheRepo.GetCapetcha(ctx, key)
	if countStr == "" {
		// 初始化计数
		_ = r.cacheRepo.SetCaptcha(ctx, key, "1", window)
		return true
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		_ = r.cacheRepo.SetCaptcha(ctx, key, "1", window)
		return true
	}

	if count >= limit {
		return false
	}

	// 自增计数，保持 TTL
	_ = r.cacheRepo.SetCaptcha(ctx, key, strconv.Itoa(count+1), window)
	return true
}

// VerifySmtpCaptcha 验证邮箱验证码
func (r *authService) VerifySmtpCaptcha(ctx context.Context, email string, code string, clear bool) (bool, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(code) == "" {
		return false, nil
	}
	// 从 Redis 获取存储的验证码
	storedCode, err := r.cacheRepo.GetCaptcha(ctx, smtpCaptchaPrefix+email)
	if err != nil {
		return false, fmt.Errorf("获取验证码失败: %v", err)
	}
	// 验证码匹配
	if storedCode == code {
		if clear {
			// 验证成功后清除验证码
			r.cacheRepo.DelCaptcha(ctx, smtpCaptchaPrefix+email)
		}
		return true, nil
	}
	return false, nil
}

func (r *authService) Register(ctx context.Context, username string, password string, email string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(email) == "" {
		return fmt.Errorf("注册参数中有空字符串")
	}
	var user models.User
	user.Email = email
	user.Username = username
	user.Password = password
	if err := r.userService.CreateUser(&user); err != nil {
		return fmt.Errorf("系统创建用户时报错 %v", err)
	}

	// 创建默认收藏文件夹
	defaultFolder := &models.FavoriteFolder{
		UserID:      user.ID,
		Name:        "默认收藏夹",
		Description: "系统默认收藏夹，用于存放未分类的收藏",
		SortOrder:   0,
	}
	if err := r.favoriteService.CreateFolder(defaultFolder); err != nil {
		r.logger.Warn("创建默认收藏文件夹失败", logger.Err("error", err))
		// 不影响用户注册，只记录警告
	}

	return nil
}

func (s *authService) GetUserByUsername(username string) (*dtoResponse.UserProfile, error) {
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userProfile := &dtoResponse.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    stringToPtr(user.Avatar),
		Bio:       stringToPtr(user.Bio),
		Website:   nil,
		Github:    nil,
		Role:      int(user.Role),
		Status:    int(user.Status),
		CreatedAt: user.CreatedAt,
	}

	return userProfile, nil
}

func (s *authService) GetUserByEmail(email string) (*dtoResponse.UserProfile, error) {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	userProfile := &dtoResponse.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    stringToPtr(user.Avatar),
		Bio:       stringToPtr(user.Bio),
		Website:   nil,
		Github:    nil,
		Role:      int(user.Role),
		Status:    int(user.Status),
		CreatedAt: user.CreatedAt,
	}

	return userProfile, nil
}

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
