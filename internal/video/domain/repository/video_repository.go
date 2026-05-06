package repository

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

type VideoRepository interface {
	Create(ctx context.Context, v *entity.Video) error
	GetByID(ctx context.Context, id int64) (*entity.Video, error)
	Update(ctx context.Context, v *entity.Video) error
	Delete(ctx context.Context, id int64) error
	ListByCategory(ctx context.Context, categoryID int, page, size int) ([]*entity.Video, int64, error)
	ListByUser(ctx context.Context, userID int64, page, size int) ([]*entity.Video, int64, error)
	ListPublished(ctx context.Context, page, size int) ([]*entity.Video, int64, error)
	IncrementView(ctx context.Context, id int64) error
	UpdateCounts(ctx context.Context, id int64, likeDelta, commentDelta int64) error
}
