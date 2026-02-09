package repository

import (
	"context"
	"time"
)

// CacheRepo 缓存仓库接口
type CacheRepo interface {
	// 验证码相关
	SetCaptcha(ctx context.Context, key string, value string, expiration time.Duration) error
	GetCaptcha(ctx context.Context, key string) (string, error)
	DelCaptcha(ctx context.Context, key string) error
	// 兼容旧方法名
	GetCapetcha(ctx context.Context, key string) (string, error)

	// 通用缓存
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
