package query

import (
	"context"
	"fmt"
	"time"

	"vicomova/internal/video/domain/entity"
	"vicomova/pkg/constants"

	"github.com/google/uuid"
)

// GetVideoByID 根据 ID 获取已发布视频
func (s *VideoQueryService) GetVideoByID(ctx context.Context, videoID int64) (*GetVideoResult, error) {
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if video == nil || !video.IsPublished() {
		return nil, nil
	}
	return &GetVideoResult{Video: video}, nil
}

// GetVideoStream 获取视频流信息
func (s *VideoQueryService) GetVideoStream(ctx context.Context, videoID int64) (*GetVideoStreamResult, error) {
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if video == nil || !video.IsPublished() {
		return nil, nil
	}
	return &GetVideoStreamResult{Video: video}, nil
}

// ListByCategory 按分类列出视频
func (s *VideoQueryService) ListByCategory(ctx context.Context, categoryID, page, size int) (*ListVideoResult, error) {
	videos, total, err := s.repo.ListByCategory(ctx, categoryID, page, size)
	if err != nil {
		return nil, err
	}
	return &ListVideoResult{Videos: videos, Total: total}, nil
}

// ListByUser 列出用户发布的视频
func (s *VideoQueryService) ListByUser(ctx context.Context, userID int64, page, size int) (*ListVideoResult, error) {
	videos, total, err := s.repo.ListByUser(ctx, userID, page, size)
	if err != nil {
		return nil, err
	}
	return &ListVideoResult{Videos: videos, Total: total}, nil
}

// ListHot 列出热门视频
func (s *VideoQueryService) ListHot(ctx context.Context, limit int) ([]*entity.Video, error) {
	return s.hotAlgo.GetHotVideos(ctx, limit)
}

// IncrementView 增加视频浏览量
func (s *VideoQueryService) IncrementView(ctx context.Context, videoID int64) error {
	return s.repo.IncrementView(ctx, videoID)
}

// GetVideoCover 获取视频封面
func (s *VideoQueryService) GetVideoCover(ctx context.Context, videoID int64) (*GetVideoCoverResult, error) {
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if video == nil {
		return nil, nil
	}
	return &GetVideoCoverResult{CoverURL: video.CoverURL}, nil
}

// GetUploadToken 获取上传凭证，供前端直传到 OSS（key 和过期时间由后端生成）
// key 格式：{user_id}/{year}/{month}/{day}/{uuid}
func (s *VideoQueryService) GetUploadToken(ctx context.Context, userID int64) (*GetUploadTokenResult, error) {
	// 生成 key：{user_id}/{year}/{month}/{uuid}
	now := time.Now()
	key := fmt.Sprintf("%d/%d/%02d/%02d/%s", userID, now.Year(), now.Month(), now.Day(), uuid.New().String())

	expire := constants.UploadTokenExpire * time.Second

	token, domain, err := s.oss.GetUploadToken(ctx, key, expire)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload token: %w", err)
	}

	return &GetUploadTokenResult{
		Token:  token,
		Key:    key,
		Domain: domain,
	}, nil
}
