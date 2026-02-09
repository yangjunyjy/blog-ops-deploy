package middleware

import (
	"context"
	"net/http"
	"strings"

	"my-blog-backend/internal/config"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/token"

	"github.com/gin-gonic/gin"
)

var (
	tokenManager token.Manager
	appConfig    *config.Config
	totpInstance *token.TOTP
)

// InitTokenManager 初始化Token管理器
func InitTokenManager(manager token.Manager) {
	tokenManager = manager
}

// InitConfig 初始化配置
func InitConfig(cfg *config.Config) {
	appConfig = cfg

	// 初始化TOTP实例（如果启用）
	if cfg.Auth.TOTP.Enabled {
		var err error
		totpInstance, err = token.NewTOTP(&cfg.Auth)
		if err != nil {
			// 记录错误但不阻止启动，TOTP只是可选项
			totpInstance = nil
		}
	}
}

// GetTokenManager 获取Token管理器
func GetTokenManager() token.Manager {
	return tokenManager
}

// OnceTokenMiddleware 一次性Token中间件
func OnceTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenID := c.Query("once_token")
		if tokenID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少一次性Token",
			})
			c.Abort()
			return
		}
		logger.Info("token为", logger.String("token", tokenID))

		// 根据Token格式和配置选择验证方式
		var verified bool
		var tokenInfo *token.TokenInfo
		var err error

		// 检查是否是TOTP Token（6位数字）
		if totpInstance != nil && len(tokenID) == 6 && strings.TrimSpace(tokenID) == tokenID {
			// 尝试验证TOTP Token
			verified, err = totpInstance.Verify(tokenID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "TOTP验证系统错误",
				})
				c.Abort()
				return
			}

			if !verified {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Token无效或已过期",
				})
				c.Abort()
				return
			}

			// TOTP验证成功，需要从Session中获取用户信息
			userID, exists := GetCurrentUserID(c)
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "未登录或Session已过期",
				})
				c.Abort()
				return
			}

			username := GetCurrentUsername(c)

			// 创建虚拟的TokenInfo用于存储用户信息
			tokenInfo = &token.TokenInfo{
				UserID:   userID,
				Username: username,
			}

			// 将用户信息存入context（与Redis Token保持一致）
			c.Set("token_user_id", userID)
			c.Set("token_username", username)
			c.Set("token_info", tokenInfo)
			c.Next()
			return
		}

		// 否则使用传统的Redis Token验证
		if tokenManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Token管理器未初始化",
			})
			c.Abort()
			return
		}

		// 验证Token（使用后自动删除）
		tokenInfo, err = tokenManager.Verify(context.Background(), tokenID)
		if err != nil {
			code := http.StatusUnauthorized
			message := "Token无效或已过期"

			if err == token.ErrTokenExpired {
				message = "Token已过期"
			}

			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
			c.Abort()
			return
		}

		// 将Token信息存入context
		c.Set("token_user_id", tokenInfo.UserID)
		c.Set("token_username", tokenInfo.Username)
		c.Set("token_info", tokenInfo)

		c.Next()
	}
}

// RequireTokenUserID 验证Token用户ID与当前用户ID一致
func RequireTokenUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录",
			})
			c.Abort()
			return
		}

		tokenUserID, exists := c.Get("token_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token无效",
			})
			c.Abort()
			return
		}

		if currentUserID != tokenUserID.(uint64) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Token用户与当前用户不匹配",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
