package query

import (
	"context"

	"vicomova/internal/video/domain/repository"
	"vicomova/internal/video/domain/service"
	"vicomova/pkg/infrastructure/oss"
)

type VideoQueryService struct {
	repo              repository.VideoRepository
	cache             repository.VideoCache
	hotAlgo           service.HotAlgorithm
	oss               oss.OSS
	viewCountService  service.ViewCountService
}

func NewVideoQueryService(
	repo repository.VideoRepository,
	cache repository.VideoCache,
	hotAlgo service.HotAlgorithm,
	ossClient oss.OSS,
	viewCountService service.ViewCountService,
) *VideoQueryService {
	return &VideoQueryService{
		repo:              repo,
		cache:             cache,
		hotAlgo:           hotAlgo,
		oss:               ossClient,
		viewCountService:  viewCountService,
	}
}

func (s *VideoQueryService) IncrementView(ctx context.Context, videoID int64) error {
	if s.viewCountService == nil {
		return nil
	}
	return s.viewCountService.Record(ctx, videoID)
}