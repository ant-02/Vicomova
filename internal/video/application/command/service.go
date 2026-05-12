package command

import (
	"vicomova/internal/video/domain/repository"
	"vicomova/pkg/infrastructure/oss"
)

type VideoCommandService struct {
	repo  repository.VideoRepository
	cache repository.VideoCache
	oss   oss.OSS
}

func NewVideoCommandService(repo repository.VideoRepository, cache repository.VideoCache, ossClient oss.OSS) *VideoCommandService {
	return &VideoCommandService{
		repo:  repo,
		cache: cache,
		oss:   ossClient,
	}
}
