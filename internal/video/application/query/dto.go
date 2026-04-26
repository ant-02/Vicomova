package query

import (
    "vicomova/internal/video/domain/entity"
)

type GetVideoResult struct {
    Video *entity.Video
}

type ListVideosResult struct {
    Videos []*entity.Video
    Total  int64
}
