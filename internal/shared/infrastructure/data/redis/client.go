package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func (c *Client) RDB() *redis.Client {
	return c.rdb
}

// Set delegates to underlying redis client.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	return c.rdb.Set(ctx, key, value, ttl)
}

// Get delegates to underlying redis client.
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.rdb.Get(ctx, key)
}

// Del delegates to underlying redis client.
func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.rdb.Del(ctx, keys...)
}

// Scan delegates to underlying redis client.
func (c *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.rdb.Scan(ctx, cursor, match, count)
}

// Ping delegates to underlying redis client.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	if c.rdb != nil {
		return c.rdb.Close()
	}
	return nil
}

var (
	client  *Client
	once    sync.Once
	initErr error
)

func Init(cfg *config.RedisConfig) error {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.Addr(),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("failed to connect to redis: %w", err)
			return
		}

		client = &Client{rdb: rdb}
		log.Info.Printf("Redis connected: %s", cfg.Addr())
	})
	return initErr
}

func GetClient() *Client {
	return client
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
