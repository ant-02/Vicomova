package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vicomova/internal/video/domain/entity"
	pkgredis "vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/utils"

	"github.com/redis/go-redis/v9"
)

const (
	CategoryVideosKeyPrefix  = "category:%d:videos"    // ZSET: 分类视频ID + 分数
	CategoryVideoMetaPrefix  = "category:meta:%d"       // HASH: 分类视频元数据
	CategoryVideoMetaTTL     = 30 * time.Minute         // 元数据缓存 30min
)

type CategoryVideoCache struct {
	client *pkgredis.Client
}

func NewCategoryVideoCache(client *pkgredis.Client) *CategoryVideoCache {
	return &CategoryVideoCache{client: client}
}

func (c *CategoryVideoCache) categoryKey(categoryID int) string {
	return fmt.Sprintf(CategoryVideosKeyPrefix, categoryID)
}

func (c *CategoryVideoCache) metaKey(videoID int64) string {
	return fmt.Sprintf(CategoryVideoMetaPrefix, videoID)
}

// GetVideoIDs 获取分类视频 ID 列表（按分数降序）
func (c *CategoryVideoCache) GetVideoIDs(ctx context.Context, categoryID int, offset, limit int64) ([]entity.HotVideoScore, error) {
	if c.client == nil {
		return nil, nil
	}

	key := c.categoryKey(categoryID)
	results, err := c.client.ZRevRangeWithScores(ctx, key, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	scores := make([]entity.HotVideoScore, 0, len(results))
	for _, z := range results {
		videoID, err := strconv.ParseInt(z.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		scores = append(scores, entity.HotVideoScore{
			VideoID: videoID,
			Score:   z.Score,
		})
	}

	return scores, nil
}

// GetVideoMeta 获取视频元数据
func (c *CategoryVideoCache) GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error) {
	if c.client == nil {
		return nil, nil
	}

	key := c.metaKey(videoID)
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, err
	}

	meta := &entity.HotVideoMeta{}
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
func (c *CategoryVideoCache) SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error {
	if c.client == nil || meta == nil {
		return nil
	}

	key := c.metaKey(videoID)
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

	return c.client.Expire(ctx, key, CategoryVideoMetaTTL).Err()
}

// GetVideoCount 获取分类视频总数
func (c *CategoryVideoCache) GetVideoCount(ctx context.Context, categoryID int) (int64, error) {
	if c.client == nil {
		return 0, nil
	}

	key := c.categoryKey(categoryID)
	cmd := c.client.ZCard(ctx, key)
	return cmd.Result()
}

// SetVideoIDs 批量设置分类视频（用于缓存回填）
func (c *CategoryVideoCache) SetVideoIDs(ctx context.Context, categoryID int, videos []*entity.Video) error {
	if c.client == nil || len(videos) == 0 {
		return nil
	}

	key := c.categoryKey(categoryID)
	pipe := c.client.RDB().Pipeline()

	for _, v := range videos {
		score := utils.WilsonScore(v.LikeCount, v.ViewCount, 1.96)
		pipe.ZAdd(ctx, key, redis.Z{Score: score, Member: strconv.FormatInt(v.ID, 10)})

		// 同时缓存元数据
		meta := &entity.HotVideoMeta{}
		meta.FromVideo(v)
		metaKey := c.metaKey(v.ID)
		fields := map[string]interface{}{
			"video_id":      strconv.FormatInt(meta.VideoID, 10),
			"title":         meta.Title,
			"cover_url":     meta.CoverURL,
			"user_id":       strconv.FormatInt(meta.UserID, 10),
			"duration":      strconv.Itoa(meta.Duration),
			"view_count":    strconv.FormatInt(meta.ViewCount, 10),
			"comment_count": strconv.FormatInt(meta.CommentCount, 10),
		}
		pipe.HMSet(ctx, metaKey, fields)
		pipe.Expire(ctx, metaKey, CategoryVideoMetaTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// ClearCategoryVideos 清空分类视频缓存
func (c *CategoryVideoCache) ClearCategoryVideos(ctx context.Context, categoryID int) error {
	if c.client == nil {
		return nil
	}

	key := c.categoryKey(categoryID)

	// 先删除所有元数据
	results, err := c.client.ZRange(ctx, key, 0, -1).Result()
	if err == nil {
		for _, member := range results {
			videoID, _ := strconv.ParseInt(member, 10, 64)
			metaKey := c.metaKey(videoID)
			c.client.Del(ctx, metaKey)
		}
	}

	return c.client.Del(ctx, key).Err()
}
