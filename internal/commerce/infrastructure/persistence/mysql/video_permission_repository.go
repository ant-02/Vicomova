package mysql

import (
	"context"
	"time"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type VideoPermissionRepository struct {
	mysql *sharedMysql.Client
}

func NewVideoPermissionRepository(mysqlClient *sharedMysql.Client) commerceRepo.VideoPermissionRepository {
	return &VideoPermissionRepository{mysql: mysqlClient}
}

func (r *VideoPermissionRepository) Create(ctx context.Context, p *commerceEntity.VideoPermission) error {
	po := VideoPermissionToPO(p)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("VideoPermissionRepository.Create: failed: %v", err)
		return err
	}
	p.ID = po.ID
	return nil
}

func (r *VideoPermissionRepository) Check(ctx context.Context, userID, videoID int64) (bool, error) {
	var count int64
	err := r.mysql.WithContext(ctx).Model(&VideoPermissionPO{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Where("expire_at IS NULL OR expire_at > ?", time.Now()).
		Count(&count).Error
	if err != nil {
		log.Error.Printf("VideoPermissionRepository.Check: failed: %v", err)
		return false, err
	}
	return count > 0, nil
}

func (r *VideoPermissionRepository) Get(ctx context.Context, userID, videoID int64) (*commerceEntity.VideoPermission, error) {
	var po VideoPermissionPO
	err := r.mysql.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("VideoPermissionRepository.Get: failed: %v", err)
		return nil, err
	}
	return POToVideoPermission(&po), nil
}

func (r *VideoPermissionRepository) ListByUser(ctx context.Context, userID int64) ([]*commerceEntity.VideoPermission, error) {
	var pos []VideoPermissionPO
	if err := r.mysql.WithContext(ctx).Where("user_id = ?", userID).
		Where("expire_at IS NULL OR expire_at > ?", time.Now()).
		Find(&pos).Error; err != nil {
		log.Error.Printf("VideoPermissionRepository.ListByUser: failed: %v", err)
		return nil, err
	}
	perms := make([]*commerceEntity.VideoPermission, len(pos))
	for i := range pos {
		perms[i] = POToVideoPermission(&pos[i])
	}
	return perms, nil
}

func (r *VideoPermissionRepository) ListByVideo(ctx context.Context, videoID int64) ([]*commerceEntity.VideoPermission, error) {
	var pos []VideoPermissionPO
	if err := r.mysql.WithContext(ctx).Where("video_id = ?", videoID).
		Where("expire_at IS NULL OR expire_at > ?", time.Now()).
		Find(&pos).Error; err != nil {
		log.Error.Printf("VideoPermissionRepository.ListByVideo: failed: %v", err)
		return nil, err
	}
	perms := make([]*commerceEntity.VideoPermission, len(pos))
	for i := range pos {
		perms[i] = POToVideoPermission(&pos[i])
	}
	return perms, nil
}
