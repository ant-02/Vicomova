package mysql

import (
	"context"
	"errors"

	"vicomova/internal/chat/domain/entity"
	repo "vicomova/internal/chat/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type FollowRepository struct {
	mysql *sharedMysql.Client
}

func NewFollowRepository(mysqlClient *sharedMysql.Client) repo.FollowRepository {
	return &FollowRepository{mysql: mysqlClient}
}

func (r *FollowRepository) Create(ctx context.Context, follow *entity.Follow) error {
	var existing FollowPO
	result := r.mysql.WithContext(ctx).Unscoped().
		Where("follower_id = ? AND following_id = ?", follow.FollowerID, follow.FollowingID).
		First(&existing)

	if result.Error == nil {
		if existing.DeletedAt.Valid {
			if err := r.mysql.WithContext(ctx).Unscoped().Exec("UPDATE user_follows SET deleted_at = NULL WHERE id = ?", existing.ID).Error; err != nil {
				log.Error.Printf("FollowRepository.Create: restore failed: %v", err)
				return err
			}
		}
		follow.ID = existing.ID
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error.Printf("FollowRepository.Create: query failed: %v", result.Error)
		return result.Error
	}

	po := &FollowPO{
		FollowerID:  follow.FollowerID,
		FollowingID: follow.FollowingID,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("FollowRepository.Create: insert failed: %v", err)
		return err
	}
	follow.ID = po.ID
	return nil
}

func (r *FollowRepository) Delete(ctx context.Context, followerID, followingID int64) error {
	if err := r.mysql.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&FollowPO{}).Error; err != nil {
		log.Error.Printf("FollowRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *FollowRepository) Get(ctx context.Context, followerID, followingID int64) (*entity.Follow, error) {
	var po FollowPO
	err := r.mysql.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("FollowRepository.Get: failed: %v", err)
		return nil, err
	}
	return POToFollow(&po), nil
}

func (r *FollowRepository) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&FollowPO{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FollowRepository) ListFollowers(ctx context.Context, userID int64) ([]int64, error) {
	var pos []FollowPO
	if err := r.mysql.WithContext(ctx).Model(&FollowPO{}).
		Where("following_id = ?", userID).
		Find(&pos).Error; err != nil {
		log.Error.Printf("FollowRepository.ListFollowers: failed: %v", err)
		return nil, err
	}

	ids := make([]int64, len(pos))
	for i := range pos {
		ids[i] = pos[i].FollowerID
	}
	return ids, nil
}

func (r *FollowRepository) ListFollowing(ctx context.Context, userID int64) ([]int64, error) {
	var pos []FollowPO
	if err := r.mysql.WithContext(ctx).Model(&FollowPO{}).
		Where("follower_id = ?", userID).
		Find(&pos).Error; err != nil {
		log.Error.Printf("FollowRepository.ListFollowing: failed: %v", err)
		return nil, err
	}

	ids := make([]int64, len(pos))
	for i := range pos {
		ids[i] = pos[i].FollowingID
	}
	return ids, nil
}

func (r *FollowRepository) ListFriends(ctx context.Context, userID int64) ([]int64, error) {
	var pos []FollowPO
	if err := r.mysql.WithContext(ctx).Model(&FollowPO{}).
		Where("follower_id = ? AND following_id IN (SELECT follower_id FROM user_follows WHERE following_id = ? AND deleted_at IS NULL)", userID, userID).
		Find(&pos).Error; err != nil {
		log.Error.Printf("FollowRepository.ListFriends: failed: %v", err)
		return nil, err
	}

	ids := make([]int64, len(pos))
	for i := range pos {
		ids[i] = pos[i].FollowingID
	}
	return ids, nil
}
