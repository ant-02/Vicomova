package repository

import (
	"context"
	"vicomova/internal/search/domain/entity"
)

type VideoSearchRepository interface {
	SearchByIDs(ctx context.Context, ids []int64) ([]*entity.VideoSearchResult, error)
}

type VideoMetaRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]*entity.VideoSearchResult, error)
}
