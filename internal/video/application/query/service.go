package query

import (
	usergrpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/internal/video/domain/repository"
	"vicomova/internal/video/domain/service"
	"vicomova/pkg/infrastructure/oss"
)

type VideoQueryService struct {
	repo              repository.VideoRepository
	cache             repository.VideoCache
	hotAlgo           service.HotAlgorithm
	oss               oss.OSS
	viewCountProducer repository.ViewCountProducer
	hotCache          repository.HotVideoCache
	userClient        *usergrpc.UserClient
}

func NewVideoQueryService(
	repo repository.VideoRepository,
	cache repository.VideoCache,
	hotAlgo service.HotAlgorithm,
	ossClient oss.OSS,
	viewCountProducer repository.ViewCountProducer,
	hotCache repository.HotVideoCache,
	userClient *usergrpc.UserClient,
) *VideoQueryService {
	return &VideoQueryService{
		repo:              repo,
		cache:             cache,
		hotAlgo:           hotAlgo,
		oss:               ossClient,
		viewCountProducer: viewCountProducer,
		hotCache:          hotCache,
		userClient:        userClient,
	}
}
