package service

import "context"

// ViewCountService 播放量领域服务接口
type ViewCountService interface {
	// Record 记录播放量
	Record(ctx context.Context, videoID int64) error
}