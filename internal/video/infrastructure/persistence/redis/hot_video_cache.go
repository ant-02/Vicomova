package redis

import (
	"context"
	"strconv"
	"time"

	"vicomova/internal/video/domain/repository"
	pkgredis "vicomova/pkg/infrastructure/redis"

	"github.com/redis/go-redis/v9"
)

const (
	HotVideosKey     = "hot:videos"     // ZSET: 热门视频ID + 分数
	HotMetaKeyPrefix = "hot:meta:"      // HASH: 视频元数据
	HotMetaTTL       = 30 * time.Minute // 元数据缓存 30min
)

type HotVideoCache struct {
	client *pkgredis.Client
}

func NewHotVideoCache(client *pkgredis.Client) *HotVideoCache {
	return &HotVideoCache{client: client}
}

// GetHotVideoIDs 获取热门视频 ID 列表（按分数降序）
func (c *HotVideoCache) GetHotVideoIDs(ctx context.Context, offset, limit int64) ([]repository.HotVideoScore, error) {
	if c.client == nil {
		return nil, nil
	}

	results, err := c.client.ZRevRangeWithScores(ctx, HotVideosKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	scores := make([]repository.HotVideoScore, 0, len(results))
	for _, z := range results {
		videoID, err := strconv.ParseInt(z.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		scores = append(scores, repository.HotVideoScore{
			VideoID: videoID,
			Score:   z.Score,
		})
	}

	return scores, nil
}

// GetHotVideoCount 获取热门视频总数
func (c *HotVideoCache) GetHotVideoCount(ctx context.Context) (int64, error) {
	if c.client == nil {
		return 0, nil
	}

	cmd := c.client.ZCard(ctx, HotVideosKey)
	return cmd.Result()
}

// GetVideoMeta 获取视频元数据
func (c *HotVideoCache) GetVideoMeta(ctx context.Context, videoID int64) (*repository.HotVideoMeta, error) {
	if c.client == nil {
		return nil, nil
	}

	key := HotMetaKeyPrefix + strconv.FormatInt(videoID, 10)
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, err
	}

	meta := &repository.HotVideoMeta{}
	if v, ok := data["video_id"]; ok {
		meta.VideoID, _ = strconv.ParseInt(v, 10, 64)
	}
	meta.Title = data["title"]
	meta.CoverURL = data["cover_url"]
	if v, ok := data["user_id"]; ok {
		meta.UserID, _ = strconv.ParseInt(v, 10, 64)
	}
	if v, ok := data["duration"]; ok {
		meta.Duration, _ = strconv.Atoi(v)
	}
	if v, ok := data["view_count"]; ok {
		meta.ViewCount, _ = strconv.ParseInt(v, 10, 64)
	}
	if v, ok := data["comment_count"]; ok {
		meta.CommentCount, _ = strconv.ParseInt(v, 10, 64)
	}

	return meta, nil
}

// SetVideoMeta 设置视频元数据
func (c *HotVideoCache) SetVideoMeta(ctx context.Context, videoID int64, meta *repository.HotVideoMeta) error {
	if c.client == nil || meta == nil {
		return nil
	}

	key := HotMetaKeyPrefix + strconv.FormatInt(videoID, 10)
	fields := map[string]interface{}{
		"video_id":      strconv.FormatInt(meta.VideoID, 10),
		"title":         meta.Title,
		"cover_url":     meta.CoverURL,
		"user_id":       strconv.FormatInt(meta.UserID, 10),
		"duration":      strconv.Itoa(meta.Duration),
		"view_count":    strconv.FormatInt(meta.ViewCount, 10),
		"comment_count": strconv.FormatInt(meta.CommentCount, 10),
	}

	if err := c.client.HMSet(ctx, key, fields).Err(); err != nil {
		return err
	}

	return c.client.Expire(ctx, key, HotMetaTTL).Err()
}

// SetHotVideo 设置热门视频分数
func (c *HotVideoCache) SetHotVideo(ctx context.Context, videoID int64, score float64) error {
	if c.client == nil {
		return nil
	}

	return c.client.ZAdd(ctx, HotVideosKey, redis.Z{Score: score, Member: strconv.FormatInt(videoID, 10)}).Err()
}

// RemoveHotVideo 移除热门视频
func (c *HotVideoCache) RemoveHotVideo(ctx context.Context, videoID int64) error {
	if c.client == nil {
		return nil
	}

	return c.client.ZRem(ctx, HotVideosKey, strconv.FormatInt(videoID, 10)).Err()
}

// ClearHotVideos 清空热门视频
func (c *HotVideoCache) ClearHotVideos(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	// 先删除所有元数据
	results, err := c.client.ZRange(ctx, HotVideosKey, 0, -1).Result()
	if err == nil {
		for _, z := range results {
			metaKey := HotMetaKeyPrefix + z
			c.client.Del(ctx, metaKey)
		}
	}

	return c.client.Del(ctx, HotVideosKey).Err()
}
