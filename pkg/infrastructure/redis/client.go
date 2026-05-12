package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vicomova/pkg/config"
	"vicomova/pkg/log"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func (c *Client) RDB() *redis.Client {
	return c.rdb
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	return c.rdb.Set(ctx, key, value, ttl)
}

func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.rdb.Get(ctx, key)
}

func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.rdb.Del(ctx, keys...)
}

func (c *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.rdb.Scan(ctx, cursor, match, count)
}

func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.rdb.Incr(ctx, key)
}

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
	onceErr error
	initMu  sync.Mutex
)

func Init(cfg *config.Redis) error {
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			onceErr = fmt.Errorf("failed to connect to redis: %w", err)
			return
		}

		client = &Client{rdb: rdb}
		log.Info.Printf("Redis connected: %s", addr)
	})
	return onceErr
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

func Reload(cfg *config.Redis) error {
	initMu.Lock()
	defer initMu.Unlock()

	if client != nil {
		if err := client.Close(); err != nil {
			log.Error.Printf("failed to close redis: %v", err)
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	client = &Client{rdb: rdb}
	log.Info.Printf("Redis reloaded: %s", addr)
	return nil
}
