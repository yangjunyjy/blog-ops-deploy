package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(corsConfig interface{}) gin.HandlerFunc {
	// 允许的源列表，支持配置
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否在允许的源列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.EqualFold(origin, allowedOrigin) {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


