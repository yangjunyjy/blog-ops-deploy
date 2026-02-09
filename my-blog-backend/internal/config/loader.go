package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/myblog")

	// 环境变量
	v.SetEnvPrefix("MYBLOG")
	v.AutomaticEnv()

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	fmt.Println(cfg)
	// 从环境变量覆盖敏感配置
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		cfg.Redis.Password = redisPassword
	}

	return &cfg, nil
}

func GetConfigPath() string {
	// 从环境变量获取配置路径
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 默认路径
	return filepath.Join(".", "config")
}
