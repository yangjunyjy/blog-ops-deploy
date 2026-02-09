package middleware

import (
	"net/http"
	"time"

	"my-blog-backend/internal/pkg/token"

	"github.com/gin-gonic/gin"
)

// TOTPConfig TOTP配置信息
type TOTPConfig struct {
	Enabled    bool `json:"enabled"`
	TimeStep   int  `json:"timeStep"`   // 时间步长（秒）
	WindowSize int  `json:"windowSize"` // 验证窗口大小
}

// GetTOTPConfig 获取TOTP配置信息
// 这个接口可以让前端知道是否启用TOTP，以及相关参数
func GetTOTPConfig(c *gin.Context) {
	if appConfig == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "配置未初始化",
		})
		return
	}

	config := TOTPConfig{
		Enabled:    appConfig.Auth.TOTP.Enabled,
		TimeStep:   appConfig.Auth.TOTP.TimeStep,
		WindowSize: appConfig.Auth.TOTP.WindowSize,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    config,
		"message": "获取成功",
	})
}

// GenerateTOTPToken 生成TOTP Token（供前端调用）
// 注意：这个接口应该受到保护，只有已认证的用户可以调用
func GenerateTOTPToken(c *gin.Context) {
	_, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 检查TOTP是否启用
	if totpInstance == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "TOTP未启用",
		})
		return
	}

	// 生成TOTP Token
	token, err := totpInstance.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Token失败: " + err.Error(),
		})
		return
	}

	// 获取当前时间窗口信息
	now := time.Now()
	windowStart := now.Unix() - (now.Unix() % int64(appConfig.Auth.TOTP.TimeStep))
	windowEnd := windowStart + int64(appConfig.Auth.TOTP.TimeStep)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"token":        token,
			"window_start": windowStart,
			"window_end":   windowEnd,
			"time_step":    appConfig.Auth.TOTP.TimeStep,
			"expires_in":   windowEnd - now.Unix(),
		},
		"message": "生成成功",
	})
}

// TOTPStatusResponse TOTP状态响应
type TOTPStatusResponse struct {
	Enabled     bool  `json:"enabled"`
	CurrentTime int64 `json:"current_time"`
	TimeStep    int   `json:"time_step"`
	WindowSize  int   `json:"window_size"`
}

// CheckTOTPStatus 检查TOTP状态（供前端检测配置）
func CheckTOTPStatus(c *gin.Context) {
	response := TOTPStatusResponse{
		Enabled:     appConfig != nil && appConfig.Auth.TOTP.Enabled,
		CurrentTime: time.Now().Unix(),
		TimeStep:    30, // 默认值
		WindowSize:  1,  // 默认值
	}

	if appConfig != nil {
		response.TimeStep = appConfig.Auth.TOTP.TimeStep
		response.WindowSize = appConfig.Auth.TOTP.WindowSize
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    response,
		"message": "获取成功",
	})
}

// IsTOTPEnabled 检查TOTP是否启用
func IsTOTPEnabled() bool {
	return appConfig != nil && appConfig.Auth.TOTP.Enabled && totpInstance != nil
}

// GetTOTPInstance 获取TOTP实例
func GetTOTPInstance() *token.TOTP {
	return totpInstance
}
