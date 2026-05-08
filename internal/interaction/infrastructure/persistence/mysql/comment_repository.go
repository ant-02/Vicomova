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

type CommentRepository struct {
	mysql *sharedMysql.Client
}

func NewCommentRepository(mysqlClient *sharedMysql.Client) repo.CommentRepository {
	return &CommentRepository{mysql: mysqlClient}
}

func (r *CommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	po := &CommentPO{
		UserID:   comment.UserID,
		VideoID:  comment.VideoID,
		ParentID: comment.ParentID,
		Content:  comment.Content,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("CommentRepository.Create: failed: %v", err)
		return err
	}
	comment.ID = po.ID
	return nil
}

func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	if err := r.mysql.WithContext(ctx).Delete(&CommentPO{}, id).Error; err != nil {
		log.Error.Printf("CommentRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *CommentRepository) Get(ctx context.Context, id int64) (*entity.Comment, error) {
	var po CommentPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("CommentRepository.Get: failed: %v", err)
		return nil, err
	}
	return POToComment(&po), nil
}

func (r *CommentRepository) ListByVideo(ctx context.Context, videoID int64, page, size int) ([]*entity.Comment, int64, error) {
	var pos []CommentPO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&CommentPO{}).
		Where("video_id = ? AND parent_id = 0", videoID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	comments := make([]*entity.Comment, len(pos))
	for i := range pos {
		comments[i] = POToComment(&pos[i])
	}
	return comments, total, nil
}

func (r *CommentRepository) ListByParent(ctx context.Context, parentID int64, page, size int) ([]*entity.Comment, int64, error) {
	var pos []CommentPO
	var total int64

	db := r.mysql.WithContext(ctx).Model(&CommentPO{}).
		Where("parent_id = ?", parentID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at ASC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	comments := make([]*entity.Comment, len(pos))
	for i := range pos {
		comments[i] = POToComment(&pos[i])
	}
	return comments, total, nil
}

func (r *CommentRepository) IncrementLike(ctx context.Context, id int64) error {
	if err := r.mysql.WithContext(ctx).Model(&CommentPO{}).Where("id = ?", id).
		Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
		log.Error.Printf("CommentRepository.IncrementLike: failed: %v", err)
		return err
	}
	return nil
}
