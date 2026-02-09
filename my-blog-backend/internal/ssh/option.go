package ssh

import (
	"fmt"
	"time"
)

type Option func(*Config) error

func WithHost(host string) func(*Config) error {
	return func(c *Config) error {
		if host == "" {
			return fmt.Errorf("主机地址不能为空")
		}
		c.Host = host
		return nil
	}
}

func WithPort(port uint) func(*Config) error {
	return func(c *Config) error {
		if port > 65535 {
			return fmt.Errorf("地址范围应该在1 - 65535之间")
		}
		c.Port = port
		return nil
	}
}

func WithUsername(username string) func(*Config) error {
	return func(c *Config) error {
		if username == "" {
			return fmt.Errorf("用户名不能为空")
		}
		c.Username = username
		return nil
	}
}

func WithPassword(password string) func(*Config) error {
	return func(c *Config) error {
		if password == "" {
			return fmt.Errorf("密码不能为空")
		}
		c.Password = password
		c.AuthType = AuthTypePassword
		return nil
	}
}

func WithKey(key []byte) func(*Config) error {
	return func(c *Config) error {
		if len(key) == 0 {
			return fmt.Errorf("私钥不能为空")
		}
		c.Key = key
		c.AuthType = AuthTypeKey
		return nil
	}
}

// 设置双重认证
func WithBoth(password string, key []byte) Option {
	return func(c *Config) error {
		if password == "" {
			return fmt.Errorf("密码不能为空")
		}
		if len(key) == 0 {
			return fmt.Errorf("私钥不能为空")
		}
		c.Password = password
		c.Key = key
		c.AuthType = AuthTypeBoth
		return nil
	}
}

// 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) error {
		if timeout <= 0 {
			return fmt.Errorf("超时时间不能小于等于0")
		}
		c.Timeout = timeout
		return nil
	}
}

// 设置认证类型（直接指定）
func WithAuthType(authType AuthType) Option {
	return func(c *Config) error {
		c.AuthType = authType
		return nil
	}
}
