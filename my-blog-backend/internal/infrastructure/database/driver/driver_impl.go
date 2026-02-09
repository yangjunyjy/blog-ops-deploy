package driver

import (
	"fmt"
	"my-blog-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MySQLDriver MySQL驱动
type MySQLDriver struct {
	BaseDriver
}

func (d *MySQLDriver) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// PostgresDriver PostgreSQL驱动
type PostgresDriver struct {
	BaseDriver
}

func (d *PostgresDriver) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// SQLiteDriver SQLite驱动
type SQLiteDriver struct {
	BaseDriver
}

func (d *SQLiteDriver) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s.db", cfg.Database)
	}
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
