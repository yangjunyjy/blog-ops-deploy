package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCacheRepo Redis缓存仓库实现
type RedisCacheRepo struct {
	client *redis.Client
}

// NewRedisCacheRepo 创建Redis缓存仓库
func NewRedisCacheRepo(client *redis.Client) *RedisCacheRepo {
	return &RedisCacheRepo{
		client: client,
	}
}

// SetCaptcha 设置验证码
func (r *RedisCacheRepo) SetCaptcha(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("设置验证码失败: %w", err)
	}
	return nil
}

// GetCaptcha 获取验证码
func (r *RedisCacheRepo) GetCaptcha(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("获取验证码失败: %w", err)
	}
	return val, nil
}

// GetCapetcha 获取验证码（兼容旧方法名）
func (r *RedisCacheRepo) GetCapetcha(ctx context.Context, key string) (string, error) {
	return r.GetCaptcha(ctx, key)
}

// DelCaptcha 删除验证码
func (r *RedisCacheRepo) DelCaptcha(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("删除验证码失败: %w", err)
	}
	return nil
}

// Set 设置缓存
func (r *RedisCacheRepo) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}
	return nil
}

// Get 获取缓存
func (r *RedisCacheRepo) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("获取缓存失败: %w", err)
	}
	return val, nil
}

// Delete 删除缓存
func (r *RedisCacheRepo) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// Exists 检查键是否存在
func (r *RedisCacheRepo) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查键存在失败: %w", err)
	}
	return result > 0, nil
}
