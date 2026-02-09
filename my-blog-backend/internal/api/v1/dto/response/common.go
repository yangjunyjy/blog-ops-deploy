package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// NewSuccess 创建成功响应
func NewSuccess(data any) *Response {
	return &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// NewError 创建错误响应
func NewError(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// Success 成功响应
func Success(c *gin.Context, data any, message string) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, message string, err error) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权访问"
	}
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	Error(c, http.StatusForbidden, message, nil)
}

// BadRequest 请求错误响应
func BadRequest(c *gin.Context, message string, err error) {
	if message == "" {
		message = "请求参数错误"
	}
	Error(c, http.StatusBadRequest, message, err)
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Error(c, http.StatusNotFound, message, nil)
}

// InternalError 内部错误响应
func InternalError(c *gin.Context, message string, err error) {
	if message == "" {
		message = "服务器内部错误"
	}
	Error(c, http.StatusInternalServerError, message, err)
}

// SafeDerefString 安全解引用字符串指针
func SafeDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SafeDerefUint64 安全解引用uint64指针
func SafeDerefUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

// StringPtr 返回字符串的指针
func StringPtr(s string) *string {
	return &s
}
