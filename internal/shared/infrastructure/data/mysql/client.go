package mysql

import (
	"context"
	"fmt"
	"sync"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	db *gorm.DB
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func (c *Client) WithContext(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

func (c *Client) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (c *Client) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

var (
	client  *Client
	once    sync.Once
	initErr error
)

func Init(cfg *config.DatabaseConfig) error {
	once.Do(func() {
		dsn := cfg.DSN()
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			initErr = fmt.Errorf("failed to connect to mysql: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("failed to get sql.DB: %w", err)
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

		client = &Client{db: db}
		log.Info.Printf("MySQL connected: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	})
	return initErr
}

func GetClient() *Client {
	return client
}

func GetDB() *gorm.DB {
	if client != nil {
		return client.db
	}
	return nil
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
