package repository

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

// CategoryVideoCache 分类视频缓存接口
type CategoryVideoCache interface {
	// GetVideoIDs 获取分类视频 ID 列表（按分数降序）
	GetVideoIDs(ctx context.Context, categoryID int, offset, limit int64) ([]entity.HotVideoScore, error)
	// GetVideoCount 获取分类视频总数
	GetVideoCount(ctx context.Context, categoryID int) (int64, error)
	// GetVideoMeta 获取视频元数据
	GetVideoMeta(ctx context.Context, videoID int64) (*entity.HotVideoMeta, error)
	// SetVideoMeta 设置视频元数据
	SetVideoMeta(ctx context.Context, videoID int64, meta *entity.HotVideoMeta) error
	// SetVideoIDs 批量设置分类视频（用于缓存回填）
	SetVideoIDs(ctx context.Context, categoryID int, videos []*entity.Video) error
	// ClearCategoryVideos 清空分类视频缓存
	ClearCategoryVideos(ctx context.Context, categoryID int) error
}
