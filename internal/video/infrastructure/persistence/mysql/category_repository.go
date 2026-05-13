package mysql

import (
	"context"

	"vicomova/internal/video/domain/entity"
	"vicomova/internal/video/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
)

type CategoryRepository struct {
	mysql *sharedMysql.Client
}

func NewCategoryRepository(mysqlClient *sharedMysql.Client) repository.CategoryRepository {
	return &CategoryRepository{mysql: mysqlClient}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]*entity.Category, error) {
	var pos []CategoryPO
	if err := r.mysql.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&pos).Error; err != nil {
		return nil, err
	}
	categories := make([]*entity.Category, len(pos))
	for i := range pos {
		categories[i] = POToCategory(&pos[i])
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*entity.Category, error) {
	var po CategoryPO
	if err := r.mysql.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return POToCategory(&po), nil
}
