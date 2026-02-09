// internal/pkg/database/manager.go
package database

import (
	"context"
	"errors"
	"fmt"
	"my-blog-backend/internal/config"
	"my-blog-backend/internal/infrastructure/cache"
	"my-blog-backend/internal/pkg/logger"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrDriverNotFound         = errors.New("driver not found")
)

// Manager 数据库管理器
type Manager interface {
	// 连接管理
	Connect() error
	Ping(context.Context) error

	// 工厂管理
	GetRepositoryFactory() RepositoryFactory
	GetDriverName() string

	// 事务支持
	BeginTx(context.Context) (Transaction, error)
	WithTransaction(context.Context, func(context.Context) error) error

	// 健康检查
	HealthCheck(context.Context) (HealthStatus, error)
	GetStats() Stats

	// 连接池管理
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
	SetConnMaxLifetime(time.Duration)

	// 读写分离
	UseReplica(bool)
	GetPrimaryDB() *gorm.DB
	GetReplicaDB() *gorm.DB
	GetRedisClient() *redis.Client
	GetRedisCluster() *redis.ClusterClient

	// 关闭
	Close() error
}

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`
	Details   map[string]interface{} `json:"details,omitempty"`
	LastCheck time.Time              `json:"last_check"`
	Latency   time.Duration          `json:"latency,omitempty"`
}

// Stats 统计信息
type Stats struct {
	Uptime            time.Duration `json:"uptime"`
	TotalConnections  int64         `json:"total_connections"`
	ActiveConnections int           `json:"active_connections"`
	IdleConnections   int           `json:"idle_connections"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
	QueriesExecuted   int64         `json:"queries_executed"`
	Transactions      int64         `json:"transactions"`
	CacheHits         int64         `json:"cache_hits"`
	CacheMisses       int64         `json:"cache_misses"`
}

// DatabaseManager 数据库管理器实现
type DatabaseManager struct {
	config        *config.Config
	driver        *databaseDriver
	primaryDB     *gorm.DB
	replicaDB     *gorm.DB
	redisClient   *redis.Client
	redisCluster  *redis.ClusterClient
	repoFactory   RepositoryFactory
	healthChecker *HealthChecker
	stats         Stats
	mu            sync.RWMutex
	startTime     time.Time
	connected     bool
	useReplica    bool
}

// NewManager 创建数据库管理器
func NewManager(cfg *config.Config) (*DatabaseManager, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	manager := &DatabaseManager{
		config:        cfg,
		startTime:     time.Now(),
		healthChecker: NewHealthChecker(),
	}

	// 初始化驱动
	if err := manager.initDriver(); err != nil {
		return nil, fmt.Errorf("failed to init driver: %w", err)
	}

	// 初始化主数据库
	if err := manager.initPrimaryDB(); err != nil {
		return nil, fmt.Errorf("failed to init primary DB: %w", err)
	}

	if err := manager.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to init Redis: %w", err)
	}

	// 初始化仓库工厂
	if err := manager.initRepositoryFactory(); err != nil {
		return nil, fmt.Errorf("failed to init repository factory: %w", err)
	}

	// 启动健康检查
	manager.healthChecker.Start()

	return manager, nil
}

func (m *DatabaseManager) initDriver() error {
	if m.config == nil {
		return errors.New("config is nil in manager")
	}
	driverName := string(m.config.Database.Type)
	m.driver = &databaseDriver{driverType: driverName}
	return nil
}

func (m *DatabaseManager) initPrimaryDB() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.driver == nil {
		return errors.New("driver not initialized")
	}

	// 检查 config.Database 是否为 nil
	if m.config == nil {
		return errors.New("config is nil")
	}

	db, err := m.driver.Connect(&m.config.Database)
	if err != nil {
		return err
	}

	// 配置连接池（简化版，添加nil检查）
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	poolCfg := m.config.Database.Pool
	if poolCfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(poolCfg.MaxIdleConns)
	}
	if poolCfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(poolCfg.MaxOpenConns)
	}
	if poolCfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(poolCfg.ConnMaxLifetime)
	}
	if poolCfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)
	}

	// 启用查询日志（开发环境）
	if m.config.App.Env == "development" {
		db = db.Debug()
	}

	m.primaryDB = db
	return nil
}

func (m *DatabaseManager) initReplicaDB() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Database.ReadReplica.Enabled {
		return nil
	}

	replicaCfg := config.DatabaseConfig{
		Type:     m.config.Database.Type,
		Host:     m.config.Database.ReadReplica.Host,
		Port:     m.config.Database.ReadReplica.Port,
		Database: m.config.Database.ReadReplica.Database,
		Username: m.config.Database.ReadReplica.Username,
		Password: m.config.Database.ReadReplica.Password,
		DSN:      m.config.Database.ReadReplica.DSN,
		Charset:  m.config.Database.Charset,
	}

	db, err := m.driver.Connect(&replicaCfg)
	if err != nil {
		return err
	}

	m.replicaDB = db
	return nil
}

func (m *DatabaseManager) initRedis() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	redisCfg := m.config.Redis

	if !redisCfg.Enabled {
		logger.Info("Redis is disabled")
		return nil
	}

	m.redisClient = cache.InitRedisCache(&redisCfg)

	return nil
}

func (m *DatabaseManager) initRepositoryFactory() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取Redis客户端
	var redisClient interface{}
	if m.redisCluster != nil {
		redisClient = m.redisCluster
	} else {
		redisClient = m.redisClient
	}

	// 创建工厂
	factory, err := CreateFactory(string(m.config.Database.Type), &m.config.Database, redisClient)
	if err != nil {
		return err
	}

	m.repoFactory = factory
	return nil
}

// Connect 连接数据库
func (m *DatabaseManager) Connect() error {
	m.mu.Lock()

	if m.connected {
		m.mu.Unlock()
		return nil
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 在持有锁的情况下直接检查，不调用Ping方法避免死锁
	if m.primaryDB == nil {
		m.mu.Unlock()
		return ErrDatabaseNotInitialized
	}

	if m.driver == nil {
		m.mu.Unlock()
		return ErrDriverNotFound
	}

	// 临时释放锁来调用driver.Ping
	m.mu.Unlock()
	err := m.driver.Ping(ctx, m.primaryDB)
	m.mu.Lock()

	if err != nil {
		m.mu.Unlock()
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	m.connected = true
	logger.Info("Database connected successfully")
	m.mu.Unlock()
	return nil
}

// Ping 测试连接
func (m *DatabaseManager) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 检查主数据库
	if m.primaryDB == nil {
		return ErrDatabaseNotInitialized
	}

	if err := m.driver.Ping(ctx, m.primaryDB); err != nil {
		return err
	}

	// 检查从数据库
	if m.replicaDB != nil {
		if err := m.driver.Ping(ctx, m.replicaDB); err != nil {
			return err
		}
	}

	// 检查Redis
	if m.redisClient != nil {
		if err := m.redisClient.Ping(ctx).Err(); err != nil {
			return err
		}
	}

	return nil
}

// GetRepositoryFactory 获取仓库工厂
func (m *DatabaseManager) GetRepositoryFactory() RepositoryFactory {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.repoFactory
}

// BeginTx 开始事务
func (m *DatabaseManager) BeginTx(ctx context.Context) (Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.primaryDB == nil {
		return nil, ErrDatabaseNotInitialized
	}

	tx := m.primaryDB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &databaseTransaction{
		db:          tx,
		redisClient: m.redisClient,
		startTime:   time.Now(),
	}, nil
}

// WithTransaction 执行事务
func (m *DatabaseManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.BeginTx(ctx)
	if err != nil {
		return err
	}

	// 确保事务被提交或回滚
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// 创建带事务的上下文
	txCtx := NewTransactionContext(ctx, tx)

	// 执行事务函数
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// HealthCheck 健康检查
func (m *DatabaseManager) HealthCheck(ctx context.Context) (HealthStatus, error) {
	start := time.Now()

	// 检查数据库连接
	if err := m.Ping(ctx); err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Details:   map[string]interface{}{"error": err.Error()},
			LastCheck: time.Now(),
			Latency:   time.Since(start),
		}, err
	}

	// 检查连接池
	stats := m.GetStats()

	return HealthStatus{
		Status:    "healthy",
		Details:   map[string]interface{}{"stats": stats},
		LastCheck: time.Now(),
		Latency:   time.Since(start),
	}, nil
}

// GetStats 获取统计信息
func (m *DatabaseManager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := m.stats
	stats.Uptime = time.Since(m.startTime)

	// 获取数据库连接池统计
	if m.primaryDB != nil {
		if sqlDB, err := m.primaryDB.DB(); err == nil {
			stats.ActiveConnections = sqlDB.Stats().OpenConnections
			stats.IdleConnections = sqlDB.Stats().Idle
			stats.MaxIdleClosed = int64(sqlDB.Stats().MaxIdleClosed)
			stats.MaxLifetimeClosed = int64(sqlDB.Stats().MaxLifetimeClosed)
		}
	}

	return stats
}

// GetDriverName 获取驱动名称
func (m *DatabaseManager) GetDriverName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return string(m.config.Database.Type)
}

// SetMaxIdleConns 设置最大空闲连接数
func (m *DatabaseManager) SetMaxIdleConns(n int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.primaryDB != nil {
		if sqlDB, err := m.primaryDB.DB(); err == nil {
			sqlDB.SetMaxIdleConns(n)
		}
	}
}

// SetMaxOpenConns 设置最大打开连接数
func (m *DatabaseManager) SetMaxOpenConns(n int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.primaryDB != nil {
		if sqlDB, err := m.primaryDB.DB(); err == nil {
			sqlDB.SetMaxOpenConns(n)
		}
	}
}

// SetConnMaxLifetime 设置连接最大生命周期
func (m *DatabaseManager) SetConnMaxLifetime(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.primaryDB != nil {
		if sqlDB, err := m.primaryDB.DB(); err == nil {
			sqlDB.SetConnMaxLifetime(d)
		}
	}
}

// UseReplica 切换使用从库
func (m *DatabaseManager) UseReplica(use bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.useReplica = use
}

// GetPrimaryDB 获取主数据库连接
func (m *DatabaseManager) GetPrimaryDB() *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primaryDB
}

// GetReplicaDB 获取从数据库连接
func (m *DatabaseManager) GetReplicaDB() *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.useReplica && m.replicaDB != nil {
		return m.replicaDB
	}
	return m.primaryDB
}

// GetRedisClient 获取Redis客户端
func (m *DatabaseManager) GetRedisClient() *redis.Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.redisClient
}

// GetRedisCluster 获取Redis集群客户端
func (m *DatabaseManager) GetRedisCluster() *redis.ClusterClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.redisCluster
}

// Close 关闭连接
func (m *DatabaseManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error

	// 停止健康检查
	m.healthChecker.Stop()

	// 关闭主数据库
	if m.primaryDB != nil {
		if sqlDB, err := m.primaryDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("primary DB close error: %w", err))
			}
		}
	}

	// 关闭从数据库
	if m.replicaDB != nil {
		if sqlDB, err := m.replicaDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("replica DB close error: %w", err))
			}
		}
	}

	// 关闭Redis
	if m.redisClient != nil {
		if err := m.redisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("Redis close error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	m.connected = false
	logger.Info("Database connections closed")
	return nil
}
