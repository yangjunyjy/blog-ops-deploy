package session

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// MemorySession 内存Session存储
type MemorySession struct {
	sessions map[string]*SessionInfo
	mu       sync.RWMutex
	config   *Config
}

// NewMemorySession 创建内存Session管理器
func NewMemorySession(config *Config) *MemorySession {
	ms := &MemorySession{
		sessions: make(map[string]*SessionInfo),
		config:   config,
	}
	// 启动清理过期Session的协程
	go ms.cleanupExpiredSessions()
	return ms
}

// Create 创建Session
func (ms *MemorySession) Create(info *SessionInfo) (string, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	info.CreatedAt = time.Now()
	info.ExpiresAt = time.Now().Add(time.Duration(ms.config.MaxAge) * time.Second)
	ms.sessions[sessionID] = info

	return sessionID, nil
}

// Get 获取Session信息
func (ms *MemorySession) Get(sessionID string) (*SessionInfo, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	info, exists := ms.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	return info, nil
}

// Delete 删除Session
func (ms *MemorySession) Delete(sessionID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.sessions, sessionID)
	return nil
}

// Verify 验证Session
func (ms *MemorySession) Verify(sessionID string) (*SessionInfo, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	info, exists := ms.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	// 检查是否过期
	if time.Now().After(info.ExpiresAt) {
		delete(ms.sessions, sessionID)
		return nil, ErrSessionExpired
	}

	return info, nil
}

// Refresh 刷新Session过期时间
func (ms *MemorySession) Refresh(sessionID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	info, exists := ms.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	info.ExpiresAt = time.Now().Add(time.Duration(ms.config.MaxAge) * time.Second)
	return nil
}

// SetCookie 设置Cookie
func (ms *MemorySession) SetCookie(c *gin.Context, sessionID string) {
	c.SetCookie(
		ms.config.CookieName,
		sessionID,
		ms.config.MaxAge,
		ms.config.Path,
		ms.config.Domain,
		ms.config.Secure,
		ms.config.HttpOnly,
	)
}

// GetCookie 从请求获取Cookie
func (ms *MemorySession) GetCookie(c *gin.Context) string {
	sessionID, _ := c.Cookie(ms.config.CookieName)
	return sessionID
}

// ClearCookie 清除Cookie
func (ms *MemorySession) ClearCookie(c *gin.Context) {
	c.SetCookie(
		ms.config.CookieName,
		"",
		-1,
		ms.config.Path,
		ms.config.Domain,
		ms.config.Secure,
		ms.config.HttpOnly,
	)
}

// cleanupExpiredSessions 清理过期的Session
func (ms *MemorySession) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ms.mu.Lock()
		now := time.Now()
		for sessionID, info := range ms.sessions {
			if now.After(info.ExpiresAt) {
				delete(ms.sessions, sessionID)
			}
		}
		ms.mu.Unlock()
	}
}

// CleanupExpired 立即清理过期的Session（公开方法）
func (ms *MemorySession) CleanupExpired() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	for sessionID, info := range ms.sessions {
		if now.After(info.ExpiresAt) {
			delete(ms.sessions, sessionID)
		}
	}
	return nil
}

// MarshalJSON 序列化
func (ms *MemorySession) MarshalJSON() ([]byte, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return json.Marshal(ms.sessions)
}

// GetSessionCount 获取当前Session数量
func (ms *MemorySession) GetSessionCount() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.sessions)
}

// DeleteAllByUserID 删除指定用户的所有Session
func (ms *MemorySession) DeleteAllByUserID(userID uint) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for sessionID, info := range ms.sessions {
		if info.UserID == userID {
			delete(ms.sessions, sessionID)
		}
	}
	return nil
}
