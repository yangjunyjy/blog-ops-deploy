package ssh

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type AuthType uint

const TCpNetwork string = "tcp"

const (
	AuthTypePassword AuthType = iota
	AuthTypeKey
	AuthTypeBoth
)

type SSHClient struct {
	client   *ssh.Client
	sftp     *sftp.Client
	hostID   uint
	config   *Config
	lastUsed time.Time
	mu       sync.RWMutex
}

func NewClient(ctx context.Context, opts ...Option) (*SSHClient, error) {
	cfg, err := NewConfig(opts...)
	if err != nil {
		return nil, fmt.Errorf("创建配置时出错: %v", err)
	}
	authMethods, err := cfg.BuildAuthMethods()
	if err != nil {
		return nil, fmt.Errorf("使用BuildAuthMethods()函数创建认证失败: %v", err)
	}
	sshConfig := &ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            authMethods,
		Timeout:         cfg.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	client, err := DailWithContext(ctx, TCpNetwork, fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH 连接失败: %w", err)
	}
	return &SSHClient{
		client:   client,
		config:   cfg,
		lastUsed: time.Now(),
	}, nil
}

// 手动封装带超时控制的ssh连接
func DailWithContext(ctx context.Context, network, address string, config *ssh.ClientConfig) (*ssh.Client, error) {
	type result struct {
		client *ssh.Client
		err    error
	}
	res := make(chan result, 1)

	// 在携程中启动连接
	go func() {
		client, err := ssh.Dial(network, address, config)
		res <- result{client: client, err: err}
	}()

	// 创建超时定时器
	timeoutTimer := time.NewTimer(config.Timeout)
	defer timeoutTimer.Stop()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resdata := <-res:
		return resdata.client, resdata.err
	case <-timeoutTimer.C:
		return nil, fmt.Errorf("连接超时，超时时长%v", config.Timeout)
	}
}

func (sshClient *SSHClient) Close() error {
	sshClient.mu.Lock()
	defer sshClient.mu.Unlock()
	return sshClient.client.Close()
}

func (c *SSHClient) IsAlive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, _, err := c.client.SendRequest("keepalive", true, nil)
	return err == nil
}

func (c *SSHClient) UpdateLastUsed() {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
}

func (c *SSHClient) GetClient() *ssh.Client {
	return c.client
}

func (c *SSHClient) GetLastUsed() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUsed
}

func (c *SSHClient) SetHostID(hostID uint) {
	c.mu.Lock()
	c.hostID = hostID
	c.mu.Unlock()
}

func (c *SSHClient) GetHostID() uint {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hostID
}

// CreateSFTP 创建 SFTP 客户端（如果尚未创建）
func (c *SSHClient) CreateSFTP() (*sftp.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果 SFTP 客户端已存在且可用，直接返回
	if c.sftp != nil {
		// 测试连接是否还活着
		if _, err := c.sftp.Stat("."); err == nil {
			return c.sftp, nil
		}
		// 连接已断开，不调用 Close() 避免关闭底层 SSH 连接
		// c.sftp.Close() // 注释掉，防止关闭底层 SSH 连接
		c.sftp = nil
	}

	// 创建新的 SFTP 客户端
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return nil, fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}

	c.sftp = sftpClient
	return c.sftp, nil
}

// GetSFTP 获取 SFTP 客户端（自动创建如果不存在）
func (c *SSHClient) GetSFTP() (*sftp.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果 SFTP 客户端已存在，直接返回
	if c.sftp != nil {
		// 测试连接是否还活着
		if _, err := c.sftp.Stat("."); err == nil {
			return c.sftp, nil
		}
		// 连接已断开，不调用 Close() 避免关闭底层 SSH 连接
		// c.sftp.Close() // 注释掉，防止关闭底层 SSH 连接
		c.sftp = nil
	}

	// 创建新的 SFTP 客户端
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return nil, fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}

	c.sftp = sftpClient
	return c.sftp, nil
}
