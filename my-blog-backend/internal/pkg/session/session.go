package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateSessionID 生成Session ID
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Config Session配置
type Config struct {
	CookieName string `json:"cookie_name"` // Cookie名称
	MaxAge     int    `json:"max_age"`     // 过期时间（秒）
	Path       string `json:"path"`        // Cookie路径
	Domain     string `json:"domain"`      // Cookie域名
	Secure     bool   `json:"secure"`      // 是否使用HTTPS
	HttpOnly   bool   `json:"http_only"`   // 是否仅HTTP访问
}

// SessionInfo Session信息
type SessionInfo struct {
	SessionID  string                 `json:"session_id"`
	UserID     uint                   `json:"user_id"`
	Email      string                 `json:"email"`
	Username   string                 `json:"username"`
	Nickname   string                 `json:"nickname"`
	Avatar     string                 `json:"avatar"`
	DeptID     uint                   `json:"dept_id"`  // 部门ID
	RoleIDs    []uint                 `json:"role_ids"` // 角色ID列表
	PermCodes  []string               `json:"perm_codes"`
	Data       map[string]interface{} `json:"data"` // 扩展数据
	CreatedAt  time.Time              `json:"created_at"`
	ExpiresAt  time.Time              `json:"expires_at"`
	RemoteAddr string                 `json:"remote_addr"` // 客户端IP
	UserAgent  string                 `json:"user_agent"`  // 用户代理
}

// Manager Session管理器接口
type Manager interface {
	Create(info *SessionInfo) (string, error)
	Get(sessionID string) (*SessionInfo, error)
	Delete(sessionID string) error
	Refresh(sessionID string) error
	CleanupExpired() error
	Verify(sessionID string) (*SessionInfo, error)
	SetCookie(c *gin.Context, sessionID string)
	GetCookie(c *gin.Context) string
	ClearCookie(c *gin.Context)
	DeleteAllByUserID(userID uint) error // 删除指定用户的所有Session
}

// RedisManager Redis Session管理器接口（带context）
type ContextManager interface {
	Create(ctx context.Context, info *SessionInfo) (string, error)
	Get(ctx context.Context, sessionID string) (*SessionInfo, error)
	Delete(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, sessionID string) error
	Verify(ctx context.Context, sessionID string) (*SessionInfo, error)
	SetCookie(c *gin.Context, sessionID string)
	GetCookie(c *gin.Context) string
	ClearCookie(c *gin.Context)
}

// Session相关错误
var (
	ErrSessionNotFound = &sessionError{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrSessionExpired  = &sessionError{Code: "SESSION_EXPIRED", Message: "Session expired"}
)

type sessionError struct {
	Code    string
	Message string
}

func (e *sessionError) Error() string {
	return e.Message
}
