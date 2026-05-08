package mysql

import (
	"context"
	"errors"

	"vicomova/internal/interaction/domain/entity"
	repo "vicomova/internal/interaction/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	mysql *sharedMysql.Client
}

func NewFavoriteRepository(mysqlClient *sharedMysql.Client) repo.FavoriteRepository {
	return &FavoriteRepository{mysql: mysqlClient}
}

func (r *FavoriteRepository) Create(ctx context.Context, fav *entity.Favorite) error {
	po := &FavoritePO{
		UserID:  fav.UserID,
		VideoID: fav.VideoID,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("FavoriteRepository.Create: failed: %v", err)
		return err
	}
	fav.ID = po.ID
	return nil
}

func (r *FavoriteRepository) Delete(ctx context.Context, userID, videoID int64) error {
	if err := r.mysql.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&FavoritePO{}).Error; err != nil {
		log.Error.Printf("FavoriteRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *FavoriteRepository) Get(ctx context.Context, userID, videoID int64) (*entity.Favorite, error) {
	var po FavoritePO
	err := r.mysql.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("FavoriteRepository.Get: failed: %v", err)
		return nil, err
	}
	return POToFavorite(&po), nil
}

func (r *FavoriteRepository) ListByUser(ctx context.Context, userID int64, page, size int) ([]*entity.Favorite, int64, error) {
	var pos []FavoritePO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&FavoritePO{}).
		Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	favorites := make([]*entity.Favorite, len(pos))
	for i := range pos {
		favorites[i] = POToFavorite(&pos[i])
	}
	return favorites, total, nil
}
