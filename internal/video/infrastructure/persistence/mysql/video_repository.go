package mysql

import (
	"context"

	"vicomova/internal/video/domain/entity"
	"vicomova/internal/video/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type VideoRepository struct {
	mysql *sharedMysql.Client
}

func NewVideoRepository(mysqlClient *sharedMysql.Client) repository.VideoRepository {
	return &VideoRepository{mysql: mysqlClient}
}

func (r *VideoRepository) Create(ctx context.Context, v *entity.Video) error {
	po := VideoToPO(v)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("VideoRepository.Create: failed: %v", err)
		return err
	}
	v.ID = po.ID
	return nil
}

func (r *VideoRepository) GetByID(ctx context.Context, id int64) (*entity.Video, error) {
	var po VideoPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("VideoRepository.GetByID: failed for id=%d: %v", id, err)
		return nil, err
	}
	return POToVideo(&po), nil
}

func (r *VideoRepository) Update(ctx context.Context, v *entity.Video) error {
	po := VideoToPO(v)
	if err := r.mysql.WithContext(ctx).Save(po).Error; err != nil {
		log.Error.Printf("VideoRepository.Update: failed: %v", err)
		return err
	}
	return nil
}

func (r *VideoRepository) Delete(ctx context.Context, id int64) error {
	if err := r.mysql.WithContext(ctx).Delete(&VideoPO{}, id).Error; err != nil {
		log.Error.Printf("VideoRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *VideoRepository) ListByCategory(ctx context.Context, categoryID int, page, size int) ([]*entity.Video, int64, error) {
	var pos []VideoPO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&VideoPO{}).Where("category_id = ? AND status = ?", categoryID, entity.VideoStatusPublished)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	videos := make([]*entity.Video, len(pos))
	for i := range pos {
		videos[i] = POToVideo(&pos[i])
	}
	return videos, total, nil
}

func (r *VideoRepository) ListByUser(ctx context.Context, userID int64, page, size int) ([]*entity.Video, int64, error) {
	var pos []VideoPO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&VideoPO{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	videos := make([]*entity.Video, len(pos))
	for i := range pos {
		videos[i] = POToVideo(&pos[i])
	}
	return videos, total, nil
}

func (r *VideoRepository) ListPublished(ctx context.Context, page, size int) ([]*entity.Video, int64, error) {
	var pos []VideoPO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&VideoPO{}).Where("status = ?", entity.VideoStatusPublished)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	videos := make([]*entity.Video, len(pos))
	for i := range pos {
		videos[i] = POToVideo(&pos[i])
	}
	return videos, total, nil
}

func (r *VideoRepository) IncrementView(ctx context.Context, id int64) error {
	if err := r.mysql.WithContext(ctx).Model(&VideoPO{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		log.Error.Printf("VideoRepository.IncrementView: failed: %v", err)
		return err
	}
	return nil
}

func (r *VideoRepository) UpdateCounts(ctx context.Context, id int64, likeDelta, commentDelta int64) error {
	updates := map[string]interface{}{}
	if likeDelta != 0 {
		updates["like_count"] = gorm.Expr("like_count + ?", likeDelta)
	}
	if commentDelta != 0 {
		updates["comment_count"] = gorm.Expr("comment_count + ?", commentDelta)
	}
	if len(updates) == 0 {
		return nil
	}
	if err := r.mysql.WithContext(ctx).Model(&VideoPO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error.Printf("VideoRepository.UpdateCounts: failed: %v", err)
		return err
	}
	return nil
}

func VideoToPO(v *entity.Video) *VideoPO {
	return &VideoPO{
		ID:           v.ID,
		UserID:       v.UserID,
		Title:        v.Title,
		Description:  v.Description,
		CoverURL:     v.CoverURL,
		VideoURL:     v.VideoURL,
		CategoryID:   v.CategoryID,
		ViewCount:    v.ViewCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		Duration:     v.Duration,
		Status:       v.Status,
	}
}

func POToVideo(po *VideoPO) *entity.Video {
	return &entity.Video{
		ID:           po.ID,
		UserID:       po.UserID,
		Title:        po.Title,
		Description:  po.Description,
		CoverURL:     po.CoverURL,
		VideoURL:     po.VideoURL,
		CategoryID:   po.CategoryID,
		ViewCount:    po.ViewCount,
		LikeCount:    po.LikeCount,
		CommentCount: po.CommentCount,
		Duration:     po.Duration,
		Status:       po.Status,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}
