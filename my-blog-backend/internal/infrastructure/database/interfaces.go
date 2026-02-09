package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my-blog-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


// Transaction 事务接口
type Transaction interface {
	Commit() error
	Rollback() error
	Begin() *gorm.DB
	SavePoint(name string) error
	RollbackTo(name string) error
}

// databaseTransaction 事务实现
type databaseTransaction struct {
	db          *gorm.DB
	redisClient interface{}
	startTime   time.Time
}

// Commit 提交事务
func (t *databaseTransaction) Commit() error {
	return t.db.Commit().Error
}

// Rollback 回滚事务
func (t *databaseTransaction) Rollback() error {
	return t.db.Rollback().Error
}

// Begin 开始嵌套事务
func (t *databaseTransaction) Begin() *gorm.DB {
	return t.db.Begin()
}

// SavePoint 创建保存点
func (t *databaseTransaction) SavePoint(name string) error {
	return t.db.SavePoint(name).Error
}

// RollbackTo 回滚到保存点
func (t *databaseTransaction) RollbackTo(name string) error {
	return t.db.RollbackTo(name).Error
}

// TransactionContext 事务上下文
type TransactionContext struct {
	context.Context
	tx Transaction
}

// NewTransactionContext 创建事务上下文
func NewTransactionContext(ctx context.Context, tx Transaction) *TransactionContext {
	return &TransactionContext{
		Context: ctx,
		tx:      tx,
	}
}

// GetTx 获取事务
func (c *TransactionContext) GetTx() Transaction {
	return c.tx
}

// RepositoryFactory 仓库工厂接口
type RepositoryFactory interface {
	GetUserRepository() interface{}
	GetRoleRepository() interface{}
	GetMenuRepository() interface{}
	GetPermRepository() interface{}
	GetSystemRepository() interface{}
	GetCacheRepository() interface{}
}

// HealthChecker 健康检查器
type HealthChecker struct {
	checks []func() error
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make([]func() error, 0),
	}
}

// Register 注册健康检查
func (h *HealthChecker) Register(name string, checker interface{}) {
	// TODO: 实现健康检查注册
}

// Start 启动健康检查
func (h *HealthChecker) Start() {
	// TODO: 实现健康检查定时任务
}

// Stop 停止健康检查
func (h *HealthChecker) Stop() {
	// TODO: 停止健康检查
}

// CreateFactory 创建仓库工厂
func CreateFactory(dbType string, cfg interface{}, redisClient interface{}) (RepositoryFactory, error) {
	// 简化版工厂，直接返回空工厂
	if cfg == nil {
		return &simpleFactory{}, nil
	}
	return &simpleFactory{}, nil
}


// simpleFactory 简单工厂实现
type simpleFactory struct{}

func (f *simpleFactory) GetUserRepository() interface{}          { return nil }
func (f *simpleFactory) GetRoleRepository() interface{}          { return nil }
func (f *simpleFactory) GetMenuRepository() interface{}          { return nil }
func (f *simpleFactory) GetPermRepository() interface{}          { return nil }
func (f *simpleFactory) GetSystemRepository() interface{}        { return nil }
func (f *simpleFactory) GetCacheRepository() interface{}         { return nil }

// databaseDriver 数据库驱动包装
type databaseDriver struct {
	driverType string
}

func (d *databaseDriver) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		switch d.driverType {
		case "mysql":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
		case "postgres", "postgresql":
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
		case "sqlite":
			dsn = fmt.Sprintf("%s.db", cfg.Database)
		}
	}

	switch d.driverType {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres", "postgresql":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported driver: %s", d.driverType)
	}
}

func (d *databaseDriver) Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if sqlDB == nil {
		return errors.New("sqlDB is nil")
	}
	return sqlDB.Ping()
}

