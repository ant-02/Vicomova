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

type LikeRepository struct {
	mysql *sharedMysql.Client
}

func NewLikeRepository(mysqlClient *sharedMysql.Client) repo.LikeRepository {
	return &LikeRepository{mysql: mysqlClient}
}

func (r *LikeRepository) Create(ctx context.Context, like *entity.Like) error {
	var existing LikePO
	result := r.mysql.WithContext(ctx).Unscoped().
		Where("user_id = ? AND target_type = ? AND target_id = ?", like.UserID, like.TargetType, like.TargetID).
		First(&existing)

	if result.Error == nil {
		if existing.DeletedAt.Valid {
			restore := LikePO{UpdatedAt: time.Now()}
			if err := r.mysql.WithContext(ctx).Unscoped().Model(&LikePO{}).Where("id = ?", existing.ID).Updates(&restore).Error; err != nil {
				log.Error.Printf("LikeRepository.Create: restore updated_at failed: %v", err)
				return err
			}
			if err := r.mysql.WithContext(ctx).Unscoped().Exec("UPDATE likes SET deleted_at = NULL WHERE id = ?", existing.ID).Error; err != nil {
				log.Error.Printf("LikeRepository.Create: restore deleted_at failed: %v", err)
				return err
			}
		}
		like.ID = existing.ID
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error.Printf("LikeRepository.Create: query failed: %v", result.Error)
		return result.Error
	}

	po := &LikePO{
		UserID:     like.UserID,
		TargetType: like.TargetType,
		TargetID:   like.TargetID,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("LikeRepository.Create: insert failed: %v", err)
		return err
	}
	like.ID = po.ID
	return nil
}

func (r *LikeRepository) Delete(ctx context.Context, userID int64, targetType string, targetID int64) error {
	if err := r.mysql.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&LikePO{}).Error; err != nil {
		log.Error.Printf("LikeRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *LikeRepository) Get(ctx context.Context, userID int64, targetType string, targetID int64) (*entity.Like, error) {
	var po LikePO
	err := r.mysql.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("LikeRepository.Get: failed: %v", err)
		return nil, err
	}
	return POToLike(&po), nil
}

func (r *LikeRepository) ListByUser(ctx context.Context, userID int64, targetType string, cursor int64, limit int) ([]*entity.Like, bool, error) {
	var pos []LikePO

	db := r.mysql.WithContext(ctx).Model(&LikePO{}).
		Where("user_id = ? AND target_type = ?", userID, targetType)

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

	likes := make([]*entity.Like, len(pos))
	for i := range pos {
		likes[i] = POToLike(&pos[i])
	}
	return likes, hasMore, nil
}

func (r *LikeRepository) Count(ctx context.Context, targetType string, targetID int64) (int64, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&LikePO{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
