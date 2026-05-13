package repository

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

// HotVideoCache 热门视频缓存接口
type HotVideoCache interface {
	// GetHotVideoIDs 获取热门视频 ID 列表（按分数降序）
	GetHotVideoIDs(ctx context.Context, offset, limit int64) ([]entity.HotVideoScore, error)
	// GetHotVideoCount 获取热门视频总数
	GetHotVideoCount(ctx context.Context) (int64, error)
	// SetHotVideoIDs 批量设置热门视频（用于缓存回填）
	SetHotVideoIDs(ctx context.Context, videos []*entity.Video) error
	// SetHotVideoIDsStaging 批量设置热门视频到 staging key（原子性刷新）
	SetHotVideoIDsStaging(ctx context.Context, videos []*entity.Video) error
	// SwapHotVideoIDs 原子性交换 staging 到主 key
	SwapHotVideoIDs(ctx context.Context) error
	// GetVideoMeta 获取视频元数据
	GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error)
	// SetVideoMeta 设置视频元数据
	SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error
	// SetHotVideo 设置热门视频分数
	SetHotVideo(ctx context.Context, videoID int64, score float64) error
	// RemoveHotVideo 移除热门视频
	RemoveHotVideo(ctx context.Context, videoID int64) error
	// ClearHotVideos 清空热门视频
	ClearHotVideos(ctx context.Context) error
}
