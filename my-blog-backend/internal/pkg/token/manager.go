package token

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config Token配置
type Config struct {
	ExpireTime int // 过期时间（秒），默认5分钟
	Prefix     string
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		ExpireTime: 300, // 5分钟
		Prefix:     "once_token:",
	}
}

// TokenInfo Token信息
type TokenInfo struct {
	TokenID   string                 `json:"token_id"`
	UserID    uint64                 `json:"user_id"`
	Username  string                 `json:"username"`
	CreatedAt int64                  `json:"created_at"`
	ExpiresAt int64                  `json:"expires_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Manager Token管理器
type Manager interface {
	// Generate 生成一次性Token
	Generate(ctx context.Context, userID uint64, username string, metadata map[string]interface{}) (*TokenInfo, error)

	// Verify 验证Token（使用后自动删除）
	Verify(ctx context.Context, tokenID string) (*TokenInfo, error)

	// Delete 删除Token
	Delete(ctx context.Context, tokenID string) error

	// BatchDelete 批量删除用户的Token
	BatchDelete(ctx context.Context, userID uint64) error
}

// RedisTokenManager Redis Token管理器
type RedisTokenManager struct {
	client *redis.Client
	config *Config
}

// NewRedisTokenManager 创建Redis Token管理器
func NewRedisTokenManager(client *redis.Client, config *Config) *RedisTokenManager {
	if config == nil {
		config = DefaultConfig()
	}
	return &RedisTokenManager{
		client: client,
		config: config,
	}
}

// Generate 生成一次性Token
func (m *RedisTokenManager) Generate(ctx context.Context, userID uint64, username string, metadata map[string]interface{}) (*TokenInfo, error) {
	tokenID := generateTokenID()
	now := time.Now()

	token := &TokenInfo{
		TokenID:   tokenID,
		UserID:    userID,
		Username:  username,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(time.Duration(m.config.ExpireTime) * time.Second).Unix(),
		Metadata:  metadata,
	}

	// 序列化Token信息
	data, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	// 存储到Redis，设置过期时间
	key := m.config.Prefix + tokenID
	err = m.client.Set(ctx, key, data, time.Duration(m.config.ExpireTime)*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Verify 验证Token（使用后自动删除）
func (m *RedisTokenManager) Verify(ctx context.Context, tokenID string) (*TokenInfo, error) {
	key := m.config.Prefix + tokenID

	// 获取Token
	data, err := m.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// 立即删除Token（一次性使用）
	m.client.Del(ctx, key)

	// 反序列化
	var token TokenInfo
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().Unix() > token.ExpiresAt {
		return nil, ErrTokenExpired
	}

	return &token, nil
}

// Delete 删除Token
func (m *RedisTokenManager) Delete(ctx context.Context, tokenID string) error {
	key := m.config.Prefix + tokenID
	return m.client.Del(ctx, key).Err()
}

// BatchDelete 批量删除用户的Token
func (m *RedisTokenManager) BatchDelete(ctx context.Context, userID uint64) error {
	pattern := m.config.Prefix + "*"
	iter := m.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// 遍历所有Token，删除属于该用户的
	for _, key := range keys {
		data, err := m.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var token TokenInfo
		if err := json.Unmarshal(data, &token); err != nil {
			continue
		}

		if token.UserID == userID {
			m.client.Del(ctx, key)
		}
	}

	return nil
}

// generateTokenID 生成Token ID
func generateTokenID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(16)
}

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
)
