package redis

import (
	"context"
	"strconv"

	sharedRedis "vicomova/pkg/infrastructure/redis"
)

type LotteryCache struct {
	client *sharedRedis.Client
}

func NewLotteryCache(client *sharedRedis.Client) *LotteryCache {
	return &LotteryCache{client: client}
}

func (c *LotteryCache) entryKey(lotteryID int64) string {
	return LotteryOddsPrefix + strconv.FormatInt(lotteryID, 10)
}

// HasUserEntry 检查用户是否已参与
func (c *LotteryCache) HasUserEntry(ctx context.Context, lotteryID, userID int64) (bool, error) {
	key := c.entryKey(lotteryID)
	_, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return false, nil // redis.Nil or any error means not exists
	}
	return true, nil
}

// SetUserEntry 设置用户已参与（简单实现）
func (c *LotteryCache) SetUserEntry(ctx context.Context, lotteryID, userID int64) error {
	key := c.entryKey(lotteryID)
	return c.client.Set(ctx, key, userID, LotteryTTL).Err()
}

func (c *LotteryCache) RemoveUserEntry(ctx context.Context, lotteryID, userID int64) error {
	key := c.entryKey(lotteryID)
	return c.client.Del(ctx, key).Err()
}
