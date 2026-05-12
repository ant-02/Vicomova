package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"vicomova/internal/video/domain/entity"
	pkgredis "vicomova/pkg/infrastructure/redis"
)

type VideoCache struct {
	rdb *pkgredis.Client
}

func NewVideoCache(rdb *pkgredis.Client) *VideoCache {
	return &VideoCache{rdb: rdb}
}

func (c *VideoCache) key(id int64) string {
	return fmt.Sprintf("%s%d", VideoCacheKeyPrefix, id)
}

// Get 获取视频缓存（只读，不刷新TTL，用于旁路缓存场景）
func (c *VideoCache) Get(ctx context.Context, id int64) (*entity.Video, error) {
	data, err := c.rdb.Get(ctx, c.key(id)).Result()
	if err != nil {
		return nil, err
	}

	var video entity.Video
	if err := json.Unmarshal([]byte(data), &video); err != nil {
		return nil, err
	}

	return &video, nil
}

// GetAndRefresh 获取视频缓存并刷新TTL（读写场景）
func (c *VideoCache) GetAndRefresh(ctx context.Context, id int64) (*entity.Video, error) {
	data, err := c.rdb.Get(ctx, c.key(id)).Result()
	if err != nil {
		return nil, err
	}

	var video entity.Video
	if err := json.Unmarshal([]byte(data), &video); err != nil {
		return nil, err
	}

	// 滑动 TTL：每次访问刷新过期时间
	c.rdb.Set(ctx, c.key(id), data, VideoMetaTTL)

	return &video, nil
}

// Set 设置视频缓存
func (c *VideoCache) Set(ctx context.Context, video *entity.Video) error {
	if video == nil {
		return nil
	}
	data, err := json.Marshal(video)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, c.key(video.ID), data, VideoMetaTTL).Err()
}

// Del 删除视频缓存
func (c *VideoCache) Del(ctx context.Context, id int64) error {
	return c.rdb.Del(ctx, c.key(id)).Err()
}

// IncrViewCount 增加播放量（Redis INCR）
func IncrViewCount(ctx context.Context, videoID int64) error {
	key := fmt.Sprintf(ViewCountKeyPrefix, videoID)
	return pkgredis.GetClient().Incr(ctx, key).Err()
}
