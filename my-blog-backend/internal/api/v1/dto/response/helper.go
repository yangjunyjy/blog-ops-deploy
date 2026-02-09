package response

import (
	"strings"

	"my-blog-backend/internal/config"
)

var (
	urlPrefix string
)

// InitURLPrefix 初始化URL前缀
func InitURLPrefix(cfg *config.Config) {
	urlPrefix = cfg.Upload.URLPrefix
	if urlPrefix == "" {
		urlPrefix = "http://localhost:8081"
	}
	// 移除末尾的斜杠
	urlPrefix = strings.TrimSuffix(urlPrefix, "/")
}

// BuildFullURL 构建完整的URL
func BuildFullURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}

	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	// 移除开头的斜杠
	relativePath = strings.TrimPrefix(relativePath, "/")

	return urlPrefix + "/" + relativePath
}

// BuildFullURLPtr 构建完整的URL（指针版本）
func BuildFullURLPtr(relativePath *string) *string {
	if relativePath == nil {
		return nil
	}
	fullURL := BuildFullURL(*relativePath)
	if fullURL == "" {
		return nil
	}
	return &fullURL
}
