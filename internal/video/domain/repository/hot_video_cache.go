package repository

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

// HotVideoScore 热门视频分数
type HotVideoScore struct {
	VideoID int64
	Score   float64
}

// HotVideoMeta 热门视频元数据
type HotVideoMeta struct {
	VideoID      int64
	Title        string
	CoverURL     string
	UserID       int64
	Duration     int
	ViewCount    int64
	CommentCount int64
}

// FromVideo 从 Video 实体转换
func (m *HotVideoMeta) FromVideo(v *entity.Video) {
	m.VideoID = v.ID
	m.Title = v.Title
	m.CoverURL = v.CoverURL
	m.UserID = v.UserID
	m.Duration = v.Duration
	m.ViewCount = v.ViewCount
	m.CommentCount = v.CommentCount
}

// HotVideoCache 热门视频缓存接口
type HotVideoCache interface {
	// GetHotVideoIDs 获取热门视频 ID 列表（按分数降序）
	GetHotVideoIDs(ctx context.Context, offset, limit int64) ([]HotVideoScore, error)
	// GetHotVideoCount 获取热门视频总数
	GetHotVideoCount(ctx context.Context) (int64, error)
	// GetVideoMeta 获取视频元数据
	GetVideoMeta(ctx context.Context, videoID int64) (*HotVideoMeta, error)
	// SetVideoMeta 设置视频元数据
	SetVideoMeta(ctx context.Context, videoID int64, meta *HotVideoMeta) error
	// SetHotVideo 设置热门视频分数
	SetHotVideo(ctx context.Context, videoID int64, score float64) error
	// RemoveHotVideo 移除热门视频
	RemoveHotVideo(ctx context.Context, videoID int64) error
	// ClearHotVideos 清空热门视频
	ClearHotVideos(ctx context.Context) error
}
