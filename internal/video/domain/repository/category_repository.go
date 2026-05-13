package repository

import (
	"context"

	"vicomova/internal/video/domain/entity"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]*entity.Category, error)
	GetByID(ctx context.Context, id int) (*entity.Category, error)
}
