package redis

import (
	"context"
	"fmt"
	"strconv"

	pkgredis "vicomova/pkg/infrastructure/redis"

	"github.com/redis/go-redis/v9"
)

const viewCountKeyFormat = "video:%d:views"

// GetViewCount 获取播放量增量
func GetViewCount(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf(viewCountKeyFormat, videoID)
	val, err := pkgredis.GetClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// GetPendingViewCounts 扫描所有待同步的播放量
func GetPendingViewCounts(ctx context.Context) (map[int64]int64, error) {
	result := make(map[int64]int64)
	var cursor uint64
	rdb := pkgredis.GetClient()

	for {
		keys, scanCursor, err := rdb.Scan(ctx, cursor, "video:*:views", 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				continue
			}
			count, err := strconv.ParseInt(val, 10, 64)
			if err != nil || count == 0 {
				continue
			}
			var videoID int64
			fmt.Sscanf(key, "video:%d:views", &videoID)
			if videoID > 0 {
				result[videoID] = count
			}
		}

		cursor = scanCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// DelViewCountKey 删除播放量 key
func DelViewCountKey(ctx context.Context, videoID int64) error {
	key := fmt.Sprintf(viewCountKeyFormat, videoID)
	return pkgredis.GetClient().Del(ctx, key).Err()
}
