package command

import (
	"context"

	"vicomova/internal/video/domain/entity"
	videoVO "vicomova/internal/video/domain/valueobject"
	"vicomova/pkg/errors"
	"vicomova/pkg/log"
)

// Save 保存视频草稿（创建或更新），状态为 editing
func (s *VideoCommandService) Save(ctx context.Context, cmd *SaveVideoCommand) (*SaveVideoResult, error) {
	if cmd.VideoID != 0 {
		// 更新已有视频
		v, err := s.repo.GetByID(ctx, cmd.VideoID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		if v == nil {
			return nil, errors.ErrVideoNotFound
		}
		if v.UserID != cmd.UserID {
			return nil, errors.ErrForbidden
		}
		v.Title = cmd.Title
		v.Description = cmd.Description
		v.CategoryID = cmd.CategoryID
		v.CoverURL = cmd.CoverURL
		v.VideoURL = cmd.VideoURL
		v.Duration = cmd.Duration
		v.Status = videoVO.VideoStatusEditing
		if err := s.repo.Update(ctx, v); err != nil {
			log.Error.Printf("VideoCommandService.Save: Update failed: %v", err)
			return nil, errors.ErrInternalServer
		}
		// 更新后删除缓存
		if s.cache != nil {
			s.cache.Del(ctx, v.ID)
		}
		return &SaveVideoResult{VideoID: v.ID}, nil
	}

	// 创建新视频
	v := &entity.Video{
		UserID:      cmd.UserID,
		Title:       cmd.Title,
		Description: cmd.Description,
		CategoryID:  cmd.CategoryID,
		CoverURL:    cmd.CoverURL,
		VideoURL:    cmd.VideoURL,
		Duration:    cmd.Duration,
		Status:      videoVO.VideoStatusEditing,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		log.Error.Printf("VideoCommandService.Save: Create failed: %v", err)
		return nil, errors.ErrInternalServer
	}
	return &SaveVideoResult{VideoID: v.ID}, nil
}

// Submit 提交审核，状态从 editing 变为 pending
func (s *VideoCommandService) Submit(ctx context.Context, cmd *SubmitVideoCommand) (*SubmitVideoResult, error) {
	v, err := s.repo.GetByID(ctx, cmd.VideoID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if v == nil {
		return nil, errors.ErrVideoNotFound
	}
	if v.UserID != cmd.UserID {
		return nil, errors.ErrForbidden
	}
	if v.Status != videoVO.VideoStatusEditing {
		return nil, errors.ErrVideoStatusInvalid
	}
	v.Status = videoVO.VideoStatusPending
	if err := s.repo.Update(ctx, v); err != nil {
		log.Error.Printf("VideoCommandService.Submit: Update failed: %v", err)
		return nil, errors.ErrInternalServer
	}
	return &SubmitVideoResult{VideoID: v.ID}, nil
}

// Publish 发布视频（管理员审核通过），状态从 pending 变为 published
func (s *VideoCommandService) Publish(ctx context.Context, cmd *PublishVideoCommand) (*PublishVideoResult, error) {
	v, err := s.repo.GetByID(ctx, cmd.VideoID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if v == nil {
		return nil, errors.ErrVideoNotFound
	}
	if v.Status != videoVO.VideoStatusPending {
		return nil, errors.ErrVideoStatusInvalid
	}
	v.Status = videoVO.VideoStatusPublished
	if err := s.repo.Update(ctx, v); err != nil {
		log.Error.Printf("VideoCommandService.Publish: Update failed: %v", err)
		return nil, errors.ErrInternalServer
	}
	// 删除缓存
	if s.cache != nil {
		s.cache.Del(ctx, v.ID)
	}
	return &PublishVideoResult{VideoID: v.ID}, nil
}
