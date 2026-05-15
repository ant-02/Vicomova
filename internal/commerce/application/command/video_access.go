package command

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	"vicomova/pkg/log"
)

type VideoAccessCommandService struct {
	videoPermRepo commerceRepo.VideoPermissionRepository
}

func NewVideoAccessCommandService(videoPermRepo commerceRepo.VideoPermissionRepository) *VideoAccessCommandService {
	return &VideoAccessCommandService{
		videoPermRepo: videoPermRepo,
	}
}

func (s *VideoAccessCommandService) BuyVideoAccess(ctx context.Context, cmd *BuyVideoAccessCommand) (*BuyVideoAccessResult, error) {
	perm := commerceEntity.NewVideoPermission(cmd.UserID, cmd.VideoID, "purchased", nil)
	if err := s.videoPermRepo.Create(ctx, perm); err != nil {
		log.Error.Printf("VideoAccessCommandService.BuyVideoAccess: failed: %v", err)
		return nil, err
	}
	return &BuyVideoAccessResult{
		Success: true,
	}, nil
}
