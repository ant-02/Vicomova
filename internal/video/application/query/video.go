package query

import (
	"context"
	"fmt"
	"time"

	"vicomova/internal/video/domain/entity"
	"vicomova/pkg/constants"
	videoErr "vicomova/pkg/errors"
)

// GetVideoStream 获取视频流信息（同时触发播放量统计）
func (s *VideoQueryService) GetVideoStream(ctx context.Context, videoID int64) (*GetVideoStreamResult, error) {
	// 先查缓存，命中则刷新TTL
	var video *entity.Video
	if s.cache != nil {
		cached, err := s.cache.GetAndRefresh(ctx, videoID)
		if err == nil && cached != nil {
			if !cached.IsPublished() {
				return nil, videoErr.ErrVideoNotPublished
			}
			video = cached
		}
	}

	// 缓存未命中，查数据库
	if video == nil {
		var err error
		video, err = s.repo.GetByID(ctx, videoID)
		if err != nil {
			return nil, err
		}
		if video == nil {
			return nil, videoErr.ErrVideoNotFound
		}
		if !video.IsPublished() {
			return nil, videoErr.ErrVideoNotPublished
		}

		// 写入缓存
		if s.cache != nil {
			s.cache.Set(ctx, video)
		}
	}

	// 触发播放量统计（只生产一次）
	s.IncrementView(ctx, videoID)

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

// GetUploadToken 获取视频封面，供前端直传到 OSS（video_id 生成唯一 key）
// key 格式：{video_id}/{upload_type}/{year}/{month}/{day}/{timestamp}
// uploadType: 1=video, 2=cover
func (s *VideoQueryService) GetUploadToken(ctx context.Context, videoID int64, uploadType int32) (*GetUploadTokenResult, error) {
	// 检查视频是否存在
	video, err := s.repo.GetByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	if video == nil {
		return nil, fmt.Errorf("video not found")
	}

	// 生成 key：{video_id}/{upload_type}/{year}/{month}/{day}/{timestamp}
	now := time.Now()
	uploadTypeStr := map[int32]string{constants.UploadTokenTypeVideo: "video", constants.UploadTokenTypeCover: "cover"}[uploadType]
	key := fmt.Sprintf("%d/%s/%d/%02d/%02d/%d", videoID, uploadTypeStr, now.Year(), now.Month(), now.Day(), now.Unix())

	expire := constants.UploadTokenExpire * time.Second

	token, host, domain, err := s.oss.GetUploadToken(ctx, key, expire)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload token: %w", err)
	}

	return &GetUploadTokenResult{
		Token:  token,
		Key:    key,
		Domain: domain,
		Host:   host,
	}, nil
}