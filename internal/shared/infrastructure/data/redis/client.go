package redis

import (
	"context"
	"fmt"
	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info.Printf("Redis connected: %s", cfg.Addr())
	return nil
}

func GetClient() *redis.Client {
	return Client
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}