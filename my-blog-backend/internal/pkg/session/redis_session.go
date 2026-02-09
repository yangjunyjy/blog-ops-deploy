package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RedisSession Redis Session存储
type RedisSession struct {
	client *redis.Client
	config *Config
	prefix string
}

// NewRedisSession 创建Redis Session管理器
func NewRedisSession(client *redis.Client, config *Config) *RedisSession {
	return &RedisSession{
		client: client,
		config: config,
		prefix: "session:",
	}
}

// Create 创建Session
func (rs *RedisSession) Create(info *SessionInfo) (string, error) {
	ctx := context.Background()

	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	info.CreatedAt = time.Now()
	info.ExpiresAt = time.Now().Add(time.Duration(rs.config.MaxAge) * time.Second)

	// 序列化Session信息
	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	// 存储到Redis，设置过期时间
	key := rs.prefix + sessionID
	err = rs.client.Set(ctx, key, data, time.Duration(rs.config.MaxAge)*time.Second).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Get 获取Session信息
func (rs *RedisSession) Get(sessionID string) (*SessionInfo, error) {
	ctx := context.Background()
	key := rs.prefix + sessionID
	data, err := rs.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var info SessionInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().After(info.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return &info, nil
}

// Delete 删除Session
func (rs *RedisSession) Delete(sessionID string) error {
	ctx := context.Background()
	key := rs.prefix + sessionID
	return rs.client.Del(ctx, key).Err()
}

// Refresh 刷新Session过期时间
func (rs *RedisSession) Refresh(sessionID string) error {
	ctx := context.Background()
	key := rs.prefix + sessionID
	err := rs.client.Expire(ctx, key, time.Duration(rs.config.MaxAge)*time.Second).Err()
	if err != nil {
		return err
	}

	// 更新ExpiresAt
	data, err := rs.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	var info SessionInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}

	info.ExpiresAt = time.Now().Add(time.Duration(rs.config.MaxAge) * time.Second)
	updatedData, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return rs.client.Set(ctx, key, updatedData, time.Duration(rs.config.MaxAge)*time.Second).Err()
}

// Verify 验证Session
func (rs *RedisSession) Verify(sessionID string) (*SessionInfo, error) {
	return rs.Get(sessionID)
}

// SetCookie 设置Cookie
func (rs *RedisSession) SetCookie(c *gin.Context, sessionID string) {
	c.SetCookie(
		rs.config.CookieName,
		sessionID,
		rs.config.MaxAge,
		rs.config.Path,
		rs.config.Domain,
		rs.config.Secure,
		rs.config.HttpOnly,
	)
}

// GetCookie 从请求获取Cookie
func (rs *RedisSession) GetCookie(c *gin.Context) string {
	sessionID, _ := c.Cookie(rs.config.CookieName)
	return sessionID
}

// ClearCookie 清除Cookie
func (rs *RedisSession) ClearCookie(c *gin.Context) {
	c.SetCookie(
		rs.config.CookieName,
		"",
		-1,
		rs.config.Path,
		rs.config.Domain,
		rs.config.Secure,
		rs.config.HttpOnly,
	)
}

// CleanupExpired Redis不支持自动清理，依赖TTL
func (rs *RedisSession) CleanupExpired() error {
	return nil
}

// DeleteAllByUserID 删除指定用户的所有Session
func (rs *RedisSession) DeleteAllByUserID(userID uint) error {
	ctx := context.Background()

	// 使用SCAN遍历所有session键
	var cursor uint64
	for {
		keys, nextCursor, err := rs.client.Scan(ctx, cursor, rs.prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		// 检查每个Session是否属于该用户
		for _, key := range keys {
			data, err := rs.client.Get(ctx, key).Bytes()
			if err != nil {
				continue
			}

			var info SessionInfo
			if err := json.Unmarshal(data, &info); err != nil {
				continue
			}

			// 如果UserID匹配，删除该Session
			if info.UserID == userID {
				rs.client.Del(ctx, key)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
