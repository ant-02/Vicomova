package command

import (
	"vicomova/internal/video/domain/repository"
	"vicomova/pkg/infrastructure/oss"
)

type VideoCommandService struct {
	repo repository.VideoRepository
	oss  oss.OSS
}

func NewVideoCommandService(repo repository.VideoRepository, oss oss.OSS) *VideoCommandService {
	return &VideoCommandService{
		repo: repo,
		oss:  oss,
	}
}
