package command

import (
    "vicomova/internal/video/domain/repository"
    "vicomova/internal/video/infrastructure/storage"
)

type VideoCommandService struct {
    repo    repository.VideoRepository
    storage storage.VideoStorage
}

func NewVideoCommandService(repo repository.VideoRepository, storage storage.VideoStorage) *VideoCommandService {
    return &VideoCommandService{
        repo:    repo,
        storage: storage,
    }
}
