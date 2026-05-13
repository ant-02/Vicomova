package redis

import (
	"context"
	"strconv"
	"time"

	"vicomova/internal/video/domain/entity"
	"vicomova/pkg/constants"
	pkgredis "vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
	"vicomova/pkg/utils"

	"github.com/redis/go-redis/v9"
)

const (
	HotVideosKey        = "hot:videos"         // ZSET: 热门视频ID + 分数 (主key)
	HotVideosStagingKey = "hot:videos:staging" // ZSET: 刷新时的临时key
	HotMetaKeyPrefix    = "hot:meta:"          // HASH: 视频元数据
	HotMetaTTL          = 30 * time.Minute     // 元数据缓存 30min
)

type HotVideoCache struct {
	client *pkgredis.Client
}

func NewHotVideoCache(client *pkgredis.Client) *HotVideoCache {
	return &HotVideoCache{client: client}
}

// GetHotVideoIDs 获取热门视频 ID 列表（按分数降序）
func (c *HotVideoCache) GetHotVideoIDs(ctx context.Context, offset, limit int64) ([]entity.HotVideoScore, error) {
	if c.client == nil {
		return nil, nil
	}

	results, err := c.client.ZRevRangeWithScores(ctx, HotVideosKey, offset, offset+limit-1).Result()
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

// GetHotVideoCount 获取热门视频总数
func (c *HotVideoCache) GetHotVideoCount(ctx context.Context) (int64, error) {
	if c.client == nil {
		return 0, nil
	}

	cmd := c.client.ZCard(ctx, HotVideosKey)
	return cmd.Result()
}

// GetVideoMeta 获取视频元数据
func (c *HotVideoCache) GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error) {
	if c.client == nil {
		return nil, nil
	}

	key := HotMetaKeyPrefix + strconv.FormatInt(videoID, 10)
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
func (c *HotVideoCache) SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error {
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

// SetHotVideoIDs 批量设置热门视频（用于缓存回填）
func (c *HotVideoCache) SetHotVideoIDs(ctx context.Context, videos []*entity.Video) error {
	if c.client == nil || len(videos) == 0 {
		return nil
	}

	pipe := c.client.RDB().Pipeline()
	for _, v := range videos {
		score := utils.WilsonScore(v.LikeCount, v.ViewCount, constants.WilsonZ)
		pipe.ZAdd(ctx, HotVideosKey, redis.Z{Score: score, Member: strconv.FormatInt(v.ID, 10)})

		// 同时缓存元数据
		meta := &entity.HotVideoMeta{}
		meta.FromVideo(v)
		key := HotMetaKeyPrefix + strconv.FormatInt(v.ID, 10)
		fields := map[string]interface{}{
			"video_id":      strconv.FormatInt(meta.VideoID, 10),
			"title":         meta.Title,
			"cover_url":     meta.CoverURL,
			"user_id":       strconv.FormatInt(meta.UserID, 10),
			"duration":      strconv.Itoa(meta.Duration),
			"view_count":    strconv.FormatInt(meta.ViewCount, 10),
			"comment_count": strconv.FormatInt(meta.CommentCount, 10),
		}
		pipe.HMSet(ctx, key, fields)
		pipe.Expire(ctx, key, HotMetaTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// SetHotVideoIDsStaging 批量设置热门视频到 staging key（原子性刷新）
func (c *HotVideoCache) SetHotVideoIDsStaging(ctx context.Context, videos []*entity.Video) error {
	if c.client == nil || len(videos) == 0 {
		log.Debug.Printf("SetHotVideoIDsStaging: client nil or empty, returning")
		return nil
	}

	log.Debug.Printf("SetHotVideoIDsStaging: deleting staging key...")
	// 先清空 staging key
	if err := c.client.Del(ctx, HotVideosStagingKey).Err(); err != nil {
		log.Error.Printf("SetHotVideoIDsStaging: del staging failed: %v", err)
	}

	log.Debug.Printf("SetHotVideoIDsStaging: writing %d videos to staging...", len(videos))
	pipe := c.client.RDB().Pipeline()
	for _, v := range videos {
		score := utils.WilsonScore(v.LikeCount, v.ViewCount, constants.WilsonZ)
		pipe.ZAdd(ctx, HotVideosStagingKey, redis.Z{Score: score, Member: strconv.FormatInt(v.ID, 10)})

		// 同时缓存元数据
		meta := &entity.HotVideoMeta{}
		meta.FromVideo(v)
		key := HotMetaKeyPrefix + strconv.FormatInt(v.ID, 10)
		fields := map[string]interface{}{
			"video_id":      strconv.FormatInt(meta.VideoID, 10),
			"title":         meta.Title,
			"cover_url":     meta.CoverURL,
			"user_id":       strconv.FormatInt(meta.UserID, 10),
			"duration":      strconv.Itoa(meta.Duration),
			"view_count":    strconv.FormatInt(meta.ViewCount, 10),
			"comment_count": strconv.FormatInt(meta.CommentCount, 10),
		}
		pipe.HMSet(ctx, key, fields)
		pipe.Expire(ctx, key, HotMetaTTL)
	}

	_, err := pipe.Exec(ctx)
	log.Debug.Printf("SetHotVideoIDsStaging: pipeline exec done, err=%v", err)
	return err
}

// SwapHotVideoIDs 原子性交换 staging 到主 key
func (c *HotVideoCache) SwapHotVideoIDs(ctx context.Context) error {
	if c.client == nil {
		log.Debug.Printf("SwapHotVideoIDs: client nil, returning")
		return nil
	}

	log.Debug.Printf("SwapHotVideoIDs: renaming staging to main...")
	// Redis RENAME 是原子操作
	// 如果 staging 不存在会返回错误，忽略即可
	err := c.client.RDB().Rename(ctx, HotVideosStagingKey, HotVideosKey).Err()
	if err != nil {
		log.Error.Printf("SwapHotVideoIDs: rename failed: %v", err)
	} else {
		log.Debug.Printf("SwapHotVideoIDs: rename success")
	}
	return err
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
