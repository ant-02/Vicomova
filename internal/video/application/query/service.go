package query

import (
	"vicomova/internal/video/domain/repository"
	"vicomova/internal/video/domain/service"
	"vicomova/pkg/infrastructure/oss"
)

type VideoQueryService struct {
	repo    repository.VideoRepository
	hotAlgo service.HotAlgorithm
	oss     oss.OSS
}

func NewVideoQueryService(repo repository.VideoRepository, hotAlgo service.HotAlgorithm, oss oss.OSS) *VideoQueryService {
	return &VideoQueryService{
		repo:    repo,
		hotAlgo: hotAlgo,
		oss:     oss,
	}
}
