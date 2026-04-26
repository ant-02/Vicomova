package query

import (
    "context"

    "vicomova/internal/video/domain/entity"
)

func (s *VideoQueryService) GetVideo(ctx context.Context, videoID int64) (*GetVideoResult, error) {
    video, err := s.repo.GetByID(ctx, videoID)
    if err != nil {
        return nil, err
    }
    if video == nil || !video.IsPublished() {
        return nil, nil
    }
    return &GetVideoResult{Video: video}, nil
}

func (s *VideoQueryService) GetVideoStream(ctx context.Context, videoID int64) (*entity.Video, error) {
    video, err := s.repo.GetByID(ctx, videoID)
    if err != nil {
        return nil, err
    }
    if video == nil || !video.IsPublished() {
        return nil, nil
    }
    return video, nil
}

func (s *VideoQueryService) ListByCategory(ctx context.Context, categoryID int, page, size int) (*ListVideosResult, error) {
    videos, total, err := s.repo.ListByCategory(ctx, categoryID, page, size)
    if err != nil {
        return nil, err
    }
    return &ListVideosResult{Videos: videos, Total: total}, nil
}

func (s *VideoQueryService) ListPublished(ctx context.Context, page, size int) (*ListVideosResult, error) {
    videos, total, err := s.repo.ListPublished(ctx, page, size)
    if err != nil {
        return nil, err
    }
    return &ListVideosResult{Videos: videos, Total: total}, nil
}

func (s *VideoQueryService) ListByUser(ctx context.Context, userID int64, page, size int) (*ListVideosResult, error) {
    videos, total, err := s.repo.ListByUser(ctx, userID, page, size)
    if err != nil {
        return nil, err
    }
    return &ListVideosResult{Videos: videos, Total: total}, nil
}

func (s *VideoQueryService) ListHot(ctx context.Context, limit int) ([]*entity.Video, error) {
    return s.hotAlgo.GetHotVideos(ctx, limit)
}

func (s *VideoQueryService) IncrementView(ctx context.Context, videoID int64) error {
    return s.repo.IncrementView(ctx, videoID)
}
