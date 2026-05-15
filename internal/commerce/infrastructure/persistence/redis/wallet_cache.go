package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	sharedRedis "vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type WalletCache struct {
	client *sharedRedis.Client
}

func NewWalletCache(client *sharedRedis.Client) *WalletCache {
	return &WalletCache{client: client}
}

func (c *WalletCache) key(userID int64) string {
	return PointsWalletPrefix + strconv.FormatInt(userID, 10)
}

func (c *WalletCache) Get(ctx context.Context, userID int64) (*PointsWallet, error) {
	data, err := c.client.Get(ctx, c.key(userID)).Bytes()
	if err != nil {
		return nil, nil
	}
	var wallet PointsWallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (c *WalletCache) Set(ctx context.Context, wallet *PointsWallet) error {
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(wallet.UserID), data, WalletCacheTTL).Err()
}

func (c *WalletCache) Del(ctx context.Context, userID int64) error {
	return c.client.Del(ctx, c.key(userID)).Err()
}

type PointsWallet struct {
	UserID  int64 `json:"user_id"`
	Points  int64 `json:"points"`
	Version int64 `json:"version"`
}

func (c *WalletCache) IncrementPoints(ctx context.Context, userID int64, delta int64) (int64, error) {
	// 直接使用底层的 redis.Client
	result, err := c.client.Incr(ctx, c.key(userID)).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *WalletCache) GetOrSet(ctx context.Context, userID int64, factory func() (*PointsWallet, error)) (*PointsWallet, error) {
	wallet, err := c.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		return wallet, nil
	}
	wallet, err = factory()
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		if cacheErr := c.Set(ctx, wallet); cacheErr != nil {
			log.Warn.Printf("WalletCache.GetOrSet: cache set failed: %v", cacheErr)
		}
	}
	return wallet, nil
}

var _ = fmt.Sprintf // import check
