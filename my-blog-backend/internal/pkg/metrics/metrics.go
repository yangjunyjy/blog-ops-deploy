package metrics

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化指标
func Init(appName string) {
	// TODO: 实现指标收集
	// 这里可以使用 Prometheus 或其他指标收集工具
	_ = appName
}

// Metrics 中间件
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// TODO: 收集指标数据
	}
}
