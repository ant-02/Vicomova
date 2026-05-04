package mysql

import (
	"context"
	"errors"

	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	"vicomova/internal/shared/pkg/log"
	"vicomova/internal/interaction/domain/entity"
	repo "vicomova/internal/interaction/domain/repository"

	"gorm.io/gorm"
)

type LikeRepository struct {
	mysql *sharedMysql.Client
}

func NewLikeRepository(mysqlClient *sharedMysql.Client) repo.LikeRepository {
	return &LikeRepository{mysql: mysqlClient}
}

func (r *LikeRepository) Create(ctx context.Context, like *entity.Like) error {
	po := &LikePO{
		UserID:     like.UserID,
		TargetType: like.TargetType,
		TargetID:  like.TargetID,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("LikeRepository.Create: failed: %v", err)
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

func (r *LikeRepository) ListByUser(ctx context.Context, userID int64, targetType string, page, size int) ([]*entity.Like, int64, error) {
	var pos []LikePO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&LikePO{}).
		Where("user_id = ? AND target_type = ?", userID, targetType)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	likes := make([]*entity.Like, len(pos))
	for i := range pos {
		likes[i] = POToLike(&pos[i])
	}
	return likes, total, nil
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