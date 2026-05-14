package mysql

import (
	"context"
	"errors"
	"time"

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
	var existing FavoritePO
	result := r.mysql.WithContext(ctx).Unscoped().
		Where("user_id = ? AND video_id = ?", fav.UserID, fav.VideoID).
		First(&existing)

	if result.Error == nil {
		if existing.DeletedAt.Valid {
			restore := FavoritePO{UpdatedAt: time.Now()}
			if err := r.mysql.WithContext(ctx).Unscoped().Model(&FavoritePO{}).Where("id = ?", existing.ID).Updates(&restore).Error; err != nil {
				log.Error.Printf("FavoriteRepository.Create: restore updated_at failed: %v", err)
				return err
			}
			if err := r.mysql.WithContext(ctx).Unscoped().Exec("UPDATE favorites SET deleted_at = NULL WHERE id = ?", existing.ID).Error; err != nil {
				log.Error.Printf("FavoriteRepository.Create: restore deleted_at failed: %v", err)
				return err
			}
		}
		fav.ID = existing.ID
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error.Printf("FavoriteRepository.Create: query failed: %v", result.Error)
		return result.Error
	}

	po := &FavoritePO{
		UserID:  fav.UserID,
		VideoID: fav.VideoID,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("FavoriteRepository.Create: insert failed: %v", err)
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

func (r *FavoriteRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*entity.Favorite, bool, error) {
	var pos []FavoritePO

	db := r.mysql.WithContext(ctx).Model(&FavoritePO{}).
		Where("user_id = ?", userID)

	if cursor > 0 {
		db = db.Where("created_at < ?", cursor)
	}

	if err := db.Order("created_at DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		return nil, false, err
	}

	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}

	favorites := make([]*entity.Favorite, len(pos))
	for i := range pos {
		favorites[i] = POToFavorite(&pos[i])
	}
	return favorites, hasMore, nil
}
