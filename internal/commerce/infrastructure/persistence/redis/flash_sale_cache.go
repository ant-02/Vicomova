package redis

import (
	"context"
	"fmt"
	"strconv"

	sharedRedis "vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type FlashSaleCache struct {
	client *sharedRedis.Client
}

func NewFlashSaleCache(client *sharedRedis.Client) *FlashSaleCache {
	return &FlashSaleCache{client: client}
}

func (c *FlashSaleCache) stockKey(flashSaleID, productID int64) string {
	return fmt.Sprintf("%s%d:%d", FlashSaleStockPrefix, flashSaleID, productID)
}

func (c *FlashSaleCache) InitStock(ctx context.Context, flashSaleID, productID int64, stock int) error {
	key := c.stockKey(flashSaleID, productID)
	return c.client.Set(ctx, key, stock, FlashSaleTTL).Err()
}

func (c *FlashSaleCache) GetStock(ctx context.Context, flashSaleID, productID int64) (int, error) {
	key := c.stockKey(flashSaleID, productID)
	result, err := c.client.Get(ctx, key).Int()
	if err != nil {
		return 0, nil
	}
	return result, nil
}

func (c *FlashSaleCache) DecrementStock(ctx context.Context, flashSaleID, productID int64) (bool, error) {
	key := c.stockKey(flashSaleID, productID)
	// 使用 DECR 并检查结果
	result, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		log.Error.Printf("FlashSaleCache.DecrementStock: failed: %v", err)
		return false, err
	}
	// 如果值小于0，说明库存已空，恢复并返回false
	if result < 0 {
		if rstErr := c.RestoreStock(ctx, flashSaleID, productID); rstErr != nil {
			log.Error.Printf("FlashSaleCache.DecrementStock: restore stock failed: %v", rstErr)
		}
		return false, nil
	}
	return true, nil
}

func (c *FlashSaleCache) RestoreStock(ctx context.Context, flashSaleID, productID int64) error {
	key := c.stockKey(flashSaleID, productID)
	return c.client.Incr(ctx, key).Err()
}

var _ = strconv.FormatInt // import check
