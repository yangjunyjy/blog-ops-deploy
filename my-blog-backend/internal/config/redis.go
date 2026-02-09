package config

import "time"

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled      bool          `yaml:"enabled" env:"ENABLED" env-default:"false"`
	Address      string        `yaml:"address" env:"ADDRESS" env-default:"localhost:6379"`
	Addresses    []string      `yaml:"addresses" env:"ADDRESSES"`
	ClusterMode  bool          `yaml:"cluster_mode" env:"CLUSTER_MODE" env-default:"false"`
	Password     string        `yaml:"password" env:"PASSWORD"`
	DB           int           `yaml:"db" env:"DB" env-default:"0"`
	PoolSize     int           `yaml:"pool_size" env:"POOL_SIZE" env-default:"10"`
	MinIdleConns int           `yaml:"min_idle_conns" env:"MIN_IDLE_CONNS" env-default:"5"`
	MaxRetries   int           `yaml:"max_retries" env:"MAX_RETRIES" env-default:"3"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env:"DIAL_TIMEOUT" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"3s"`
	PoolTimeout  time.Duration `yaml:"pool_timeout" env:"POOL_TIMEOUT" env-default:"4s"`
}

