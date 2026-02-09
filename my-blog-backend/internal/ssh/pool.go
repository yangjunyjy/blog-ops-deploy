package ssh

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	clients map[uint]*SSHClient // hostID -> SSHClient
	mu      sync.RWMutex
	maxIdle time.Duration
	cleanup chan struct{}
	done    chan struct{} // 用于协程退出
}

func NewPool(maxIdle time.Duration) *Pool {
	pool := &Pool{
		clients: make(map[uint]*SSHClient),
		maxIdle: maxIdle,
		cleanup: make(chan struct{}),
		done:    make(chan struct{}),
	}

	// 启动清理协程
	go pool.startCleanup()

	return pool
}

// 获取或创建 SSH 客户端
func (p *Pool) Get(ctx context.Context, cfg *Config, hostID uint) (*SSHClient, error) {
	p.mu.Lock()

	// 检查是否有可用的连接
	if client, exists := p.clients[hostID]; exists {
		if client.IsAlive() {
			p.mu.Unlock()
			client.UpdateLastUsed()
			return client, nil
		}
		// 连接已失效，删除
		delete(p.clients, hostID)
		client.Close()
	}
	p.mu.Unlock()

	// 根据认证类型构建选项
	opts := p.buildOptions(cfg)

	// 创建新连接
	client, err := NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create ssh client failed: %w", err)
	}

	p.mu.Lock()
	// 再次检查，防止并发创建
	if existing, exists := p.clients[hostID]; exists {
		p.mu.Unlock()
		client.Close() // 关闭新创建的，使用已存在的
		existing.UpdateLastUsed()
		return existing, nil
	}
	client.SetHostID(hostID)
	p.clients[hostID] = client
	p.mu.Unlock()

	return client, nil
}

// 根据配置构建选项
func (p *Pool) buildOptions(cfg *Config) []Option {
	opts := []Option{
		WithHost(cfg.Host),
		WithPort(cfg.Port),
		WithUsername(cfg.Username),
		WithTimeout(cfg.Timeout),
	}

	// 根据认证类型添加认证选项
	switch cfg.AuthType {
	case AuthTypePassword:
		if cfg.Password != "" {
			opts = append(opts, WithPassword(cfg.Password))
		}
	case AuthTypeKey:
		if len(cfg.Key) > 0 {
			opts = append(opts, WithKey(cfg.Key))
		}
	case AuthTypeBoth:
		if cfg.Password != "" && len(cfg.Key) > 0 {
			opts = append(opts, WithBoth(cfg.Password, cfg.Key))
		}
	}

	return opts
}

// 清理空闲连接
func (p *Pool) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.cleanupIdle()
		case <-p.cleanup:
			return
		case <-p.done:
			return
		}
	}
}

func (p *Pool) cleanupIdle() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for hostID, client := range p.clients {
		// 需要在 SSHClient 中添加 GetLastUsed() 方法
		lastUsed := client.GetLastUsed()
		if now.Sub(lastUsed) > p.maxIdle {
			client.Close()
			delete(p.clients, hostID)
		}
	}
}

// 关闭连接池
func (p *Pool) Close() error {
	// 关闭清理协程
	select {
	case <-p.cleanup:
		// 已经关闭
	default:
		close(p.cleanup)
	}

	// 发送退出信号
	select {
	case <-p.done:
		// 已经关闭
	default:
		close(p.done)
	}

	// 等待清理协程退出
	time.Sleep(100 * time.Millisecond)

	p.mu.Lock()
	defer p.mu.Unlock()

	// 关闭所有连接
	for _, client := range p.clients {
		_ = client.Close()
	}
	p.clients = make(map[uint]*SSHClient)

	return nil
}

// 获取连接池状态
func (p *Pool) Stats() map[uint]bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := make(map[uint]bool)
	for hostID, client := range p.clients {
		stats[hostID] = client.IsAlive()
	}
	return stats
}

// 主动释放指定主机的连接
func (p *Pool) Release(hostID uint) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if client, exists := p.clients[hostID]; exists {
		_ = client.Close()
		delete(p.clients, hostID)
	}
	return nil
}
