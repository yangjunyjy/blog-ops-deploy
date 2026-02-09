package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host     string
	Port     uint
	Username string
	Password string
	Key      []byte
	AuthType AuthType
	Timeout  time.Duration
}

// 默认配置
func defaultConfig() *Config {
	return &Config{
		Port:     22,
		Timeout:  30 * time.Second,
		AuthType: AuthTypePassword,
	}
}

func NewConfig(opts ...Option) (*Config, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("应用配置失败: %v", err)
		}
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("当验证配置合法性时出错: %v", err)
	}
	return cfg, nil
}

// 验证配置
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("主机地址是必要的")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名是必要的")
	}

	switch c.AuthType {
	case AuthTypePassword:
		if c.Password == "" {
			return fmt.Errorf("密码认证时，密码是必要的")
		}
	case AuthTypeKey:
		if len(c.Key) == 0 {
			return fmt.Errorf("密钥认证时，密钥是必要的")
		}
	case AuthTypeBoth:
		if c.Password == "" || len(c.Key) == 0 {
			return fmt.Errorf("b双重认证时，密码和密钥都是必要的")
		}
	}

	return nil
}

// 构建认证方法
func (c *Config) BuildAuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod
	switch c.AuthType {
	case AuthTypePassword:
		authMethods = append(authMethods, ssh.Password(c.Password))

	case AuthTypeKey:
		signer, err := ssh.ParsePrivateKey(c.Key)
		if err != nil {
			return nil, fmt.Errorf("解析私钥出错: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))

	case AuthTypeBoth:
		authMethods = append(authMethods, ssh.Password(c.Password))

		signer, err := ssh.ParsePrivateKey(c.Key)
		if err != nil {
			return nil, fmt.Errorf("解析私钥出错: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	return authMethods, nil
}
