package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationLogService 操作日志服务接口
type OperationLogService interface {
	LogOperation(module, business, method, request, params, ip, location, agent, operator string,
		status int, errorMsg string, costTime int64) error
}

var operationLogService OperationLogService

// InitOperationLogService 初始化操作日志服务
func InitOperationLogService(service OperationLogService) {
	operationLogService = service
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 如果操作日志服务未初始化，跳过
		if operationLogService == nil {
			return
		}

		// 计算耗时
		costTime := time.Since(startTime).Milliseconds()

		// 获取用户信息
		_, exists := GetCurrentUserID(c)
		if !exists {
			return
		}
		username := GetCurrentUsername(c)

		// 请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 跳过某些路径的日志记录
		skipPaths := []string{
			"/api/v1/auth/login",
			"/api/v1/auth/logout",
			"/api/v1/auth/info",
			"/api/v1/auth/menu",
			"/api/v1/auth/token",
		}

		for _, skipPath := range skipPaths {
			if strings.HasPrefix(path, skipPath) {
				return
			}
		}

		// 状态码
		statusCode := c.Writer.Status()
		status := 0
		errorMsg := ""

		if statusCode >= http.StatusBadRequest {
			status = 1
			// 尝试从响应中获取错误信息
			if len(c.Errors) > 0 {
				errorMsg = c.Errors.String()
			}
		}

		// 请求参数
		params := ""
		if c.Request.Method != "GET" {
			// 解析JSON参数
			var data interface{}
			if err := json.Unmarshal(requestBody, &data); err == nil {
				paramsBytes, _ := json.Marshal(data)
				params = string(paramsBytes)
			} else {
				params = string(requestBody)
			}
		} else {
			params = c.Request.URL.RawQuery
		}

		// 限制参数长度
		if len(params) > 2000 {
			params = params[:2000] + "..."
		}

		// IP地址
		ip := c.ClientIP()
		location := "" // 可以使用IP库解析地理位置
		agent := c.Request.UserAgent()

		// 根据路径确定模块和业务类型
		module := "系统管理"
		business := getBusinessType(path)

		// 限制用户代理长度
		if len(agent) > 500 {
			agent = agent[:500]
		}

		// 限制错误信息长度
		if len(errorMsg) > 2000 {
			errorMsg = errorMsg[:2000] + "..."
		}

		// 记录日志（异步）
		go func() {
			operationLogService.LogOperation(
				module,
				business,
				method,
				c.Request.RequestURI,
				params,
				ip,
				location,
				agent,
				username,
				status,
				errorMsg,
				costTime,
			)
		}()
	}
}

// getBusinessType 根据路径获取业务类型
func getBusinessType(path string) string {
	if strings.Contains(path, "/users") {
		return "用户管理"
	} else if strings.Contains(path, "/roles") {
		return "角色管理"
	} else if strings.Contains(path, "/menus") {
		return "菜单管理"
	} else if strings.Contains(path, "/depts") {
		return "部门管理"
	} else if strings.Contains(path, "/posts") {
		return "岗位管理"
	} else if strings.Contains(path, "/dict") {
		return "字典管理"
	} else if strings.Contains(path, "/articles") {
		return "文章管理"
	} else if strings.Contains(path, "/comments") {
		return "评论管理"
	}
	return "其他"
}

// LogOperation 记录操作日志的辅助函数
func LogOperation(c *gin.Context, module, business string) {
	if operationLogService == nil {
		return
	}

	startTime := time.Now()
	c.Next()
	costTime := time.Since(startTime).Milliseconds()

	_, exists := GetCurrentUserID(c)
	if !exists {
		return
	}
	username := GetCurrentUsername(c)

	statusCode := c.Writer.Status()
	status := 0
	errorMsg := ""
	if statusCode >= http.StatusBadRequest {
		status = 1
	}

	operationLogService.LogOperation(
		module,
		business,
		c.Request.Method,
		c.Request.RequestURI,
		c.Request.URL.RawQuery,
		c.ClientIP(),
		"",
		c.Request.UserAgent(),
		username,
		status,
		errorMsg,
		costTime,
	)
}
