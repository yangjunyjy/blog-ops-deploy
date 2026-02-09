package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/claims"
	"my-blog-backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	permissionService PermissionService
)

var JWTSecret []byte

// PermissionService 权限服务接口
type PermissionService interface {
	CheckPermission(userID uint64, menuCode string) bool
	GetUserMenus(userID uint64) ([]*models.SysMenu, error)
}

// InitPermissionService 初始化权限服务
func InitPermissionService(service PermissionService) {
	permissionService = service
}

// IsSuperAdmin 判断是否为超级管理员
func IsSuperAdmin(roleIDs []uint) bool {
	for _, roleID := range roleIDs {
		if uint64(roleID) == models.SuperAdminRoleID {
			return true
		}
	}
	return false
}

// PermissionMiddleware 权限中间件
func PermissionMiddleware(menuCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户信息
		userID, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// 获取角色ID列表
		roleIDs := GetCurrentRoleIDs(c)

		// 超级管理员跳过权限检查
		if IsSuperAdmin(roleIDs) {
			c.Next()
			return
		}

		// 检查权限
		if permissionService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "权限服务未初始化",
			})
			c.Abort()
			return
		}

		if !permissionService.CheckPermission(userID, menuCode) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware 角色中间件
func RoleMiddleware(allowedRoleIDs ...uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := GetCurrentUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// 获取角色ID列表
		roleIDs := GetCurrentRoleIDs(c)

		// 超级管理员跳过检查
		if IsSuperAdmin(roleIDs) {
			c.Next()
			return
		}

		// 检查是否有允许的角色
		hasPermission := false
		for _, roleID := range roleIDs {
			for _, allowedID := range allowedRoleIDs {
				if uint64(roleID) == allowedID {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "角色权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateToken(tokenString string) (*claims.UserClaims, error) {
	return validateToken(tokenString)
}

// 博客前台用户认证
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 检查 Authorization header 是否存在
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少有效的token",
			})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少有效的token",
			})
			c.Abort()
			return
		}

		token := parts[1]
		userClaims, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息设置到上下文（兼容 user_claims 和 user_id 两种方式）
		c.Set("user_claims", userClaims)
		c.Set("user_id", userClaims.UserID)
		c.Set("username", userClaims.Username)
		c.Set("email", userClaims.Email)
		c.Set("role", userClaims.Role)

		c.Next()
	}
}

// 前台用户验证token
func validateToken(tokenString string) (*claims.UserClaims, error) {
	// 使用ParseWithClaims来解析到UserClaims结构
	token, err := jwt.ParseWithClaims(tokenString, &claims.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		logger.Error("JWT解析失败", logger.Err("error", err))
		return nil, err
	}

	// 提取claims
	if userClaims, ok := token.Claims.(*claims.UserClaims); ok && token.Valid {
		// 打印token的过期时间信息
		now := time.Now()
		logger.Info("JWT验证成功",
			logger.Uint("user_id", userClaims.UserID),
			logger.Time("expires_at", userClaims.ExpiresAt.Time),
			logger.Time("current_time", now),
			logger.String("time_until_expiry", userClaims.ExpiresAt.Time.Sub(now).String()),
		)
		return userClaims, nil
	}

	return nil, errors.New("invalid token claims")
}
