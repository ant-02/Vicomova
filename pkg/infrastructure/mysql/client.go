package mysql

import (
	"context"
	"fmt"
	"sync"

	"vicomova/pkg/config"
	"vicomova/pkg/log"

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
	onceErr error
	initMu  sync.Mutex
)

func Init(cfg *config.Database) error {
	once.Do(func() {
		initMu.Lock()
		defer initMu.Unlock()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			onceErr = fmt.Errorf("failed to connect to mysql: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			onceErr = fmt.Errorf("failed to get sql.DB: %w", err)
			return
		}

		if cfg.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		}

		client = &Client{db: db}
		log.Info.Printf("MySQL connected: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	})
	return onceErr
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

func Reload(cfg *config.Database) error {
	initMu.Lock()
	defer initMu.Unlock()

	if client != nil {
		if err := client.Close(); err != nil {
			log.Error.Printf("failed to close mysql: %v", err)
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	client = &Client{db: db}
	log.Info.Printf("MySQL reloaded: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	return nil
}
