package cache

import (
	"my-blog-backend/internal/config"
	"my-blog-backend/internal/pkg/logger"

	"github.com/go-redis/redis/v8"
)

func InitRedisCache(redisCfg *config.RedisConfig) *redis.Client {
	var redis2 redis.Client
	redis2 = *redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addresses[0],
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
		MaxRetries:   redisCfg.MaxRetries,
		DialTimeout:  redisCfg.DialTimeout,
		ReadTimeout:  redisCfg.ReadTimeout,
		WriteTimeout: redisCfg.WriteTimeout,
		PoolTimeout:  redisCfg.PoolTimeout,
	})
	logger.Info("redis 已经单节点初始化")
	return &redis2
}
