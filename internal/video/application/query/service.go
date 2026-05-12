package query

import (
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