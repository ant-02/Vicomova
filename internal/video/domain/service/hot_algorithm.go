package service

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

type HotAlgorithm interface {
	CalculateScore(ctx context.Context, videoID int64) (float64, error)
	CalculateAndUpdateScore(ctx context.Context, video *entity.Video) error
	IncrementView(ctx context.Context, videoID int64) error
	GetHotVideos(ctx context.Context, limit int) ([]*entity.Video, error)
}
