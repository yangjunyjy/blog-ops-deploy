package driver

import (
	"my-blog-backend/internal/config"

	"gorm.io/gorm"
)

// Driver 数据库驱动接口
type Driver interface {
	Connect(cfg *config.DatabaseConfig) (*gorm.DB, error)
	Ping(ctx interface{}, db *gorm.DB) error
}

// drivers 注册的驱动
var drivers = make(map[string]Driver)

// RegisterDriver 注册驱动
func RegisterDriver(name string, driver Driver) {
	drivers[name] = driver
}

// GetDriver 获取驱动
func GetDriver(name string) (Driver, error) {
	driver, ok := drivers[name]
	if !ok {
		return nil, nil
	}
	return driver, nil
}

// init 初始化驱动
func init() {
	RegisterDriver("mysql", &MySQLDriver{})
	RegisterDriver("postgres", &PostgresDriver{})
	RegisterDriver("postgresql", &PostgresDriver{})
	RegisterDriver("sqlite", &SQLiteDriver{})
}

// BaseDriver 基础驱动
type BaseDriver struct{}

func (d *BaseDriver) Ping(ctx interface{}, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
