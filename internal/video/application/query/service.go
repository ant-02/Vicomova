package query

import (
    "vicomova/internal/video/domain/repository"
    "vicomova/internal/video/domain/service"
)

type VideoQueryService struct {
    repo    repository.VideoRepository
    hotAlgo service.HotAlgorithm
}

func NewVideoQueryService(repo repository.VideoRepository, hotAlgo service.HotAlgorithm) *VideoQueryService {
    return &VideoQueryService{
        repo:    repo,
        hotAlgo: hotAlgo,
    }
}
