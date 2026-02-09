package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"my-blog-backend/internal/pkg/session"
)

var sessionManager session.Manager

// InitSessionManager 初始化Session管理器
func InitSessionManager(manager session.Manager) {
	sessionManager = manager
}

// GetSessionManager 获取Session管理器
func GetSessionManager() session.Manager {
	return sessionManager
}

// SessionMiddleware Session中间件
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := sessionManager.GetCookie(c)

		if sessionID == "" {
			c.Next()
			return
		}

		// 验证Session
		info, err := sessionManager.Verify(sessionID)
		if err != nil {
			// Session无效或过期，清除Cookie
			sessionManager.ClearCookie(c)
			c.Next()
			return
		}

		// 刷新Session过期时间
		if err := sessionManager.Refresh(sessionID); err != nil {
			// 刷新失败，清除Cookie
			sessionManager.ClearCookie(c)
			c.Next()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", info.UserID)
		c.Set("username", info.Username)
		c.Set("nickname", info.Nickname)
		c.Set("role_ids", info.RoleIDs)
		c.Set("dept_id", info.DeptID)
		c.Set("session_info", info)

		c.Next()
	}
}

// RequireAuth 需要认证的中间件
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists || userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RBACAuth RBAC管理后台认证中间件（基于cookie session）
func RBACAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := sessionManager.GetCookie(c)

		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// 验证Session
		info, err := sessionManager.Verify(sessionID)
		if err != nil {
			// Session无效或过期，清除Cookie
			sessionManager.ClearCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// 刷新Session过期时间
		if err := sessionManager.Refresh(sessionID); err != nil {
			// 刷新失败，清除Cookie
			sessionManager.ClearCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", info.UserID)
		c.Set("username", info.Username)
		c.Set("nickname", info.Nickname)
		c.Set("role_ids", info.RoleIDs)
		c.Set("dept_id", info.DeptID)
		c.Set("session_info", info)

		c.Next()
	}
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		return 0, false
	}
	// 支持两种类型：uint 和 uint64
	switch v := userID.(type) {
	case uint:
		return uint64(v), true
	case uint64:
		return v, true
	default:
		return 0, false
	}
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetCurrentRoleIDs 获取当前用户角色ID列表
func GetCurrentRoleIDs(c *gin.Context) []uint {
	roleIDs, exists := c.Get("role_ids")
	if !exists {
		return []uint{}
	}
	return roleIDs.([]uint)
}
