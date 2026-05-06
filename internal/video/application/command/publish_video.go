package command

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

func (s *VideoCommandService) Publish(ctx context.Context, cmd *PublishVideoCommand) (*PublishVideoResult, error) {
	video := &entity.Video{
		UserID:      cmd.UserID,
		Title:       cmd.Title,
		Description: cmd.Description,
		CategoryID:  cmd.CategoryID,
		CoverURL:    cmd.CoverURL,
		VideoURL:    cmd.VideoURL,
		Duration:    cmd.Duration,
		Status:      entity.VideoStatusPublished,
	}

	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}

	return &PublishVideoResult{VideoID: video.ID}, nil
}
