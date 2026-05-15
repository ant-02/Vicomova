package mysql

import (
	"context"

	"vicomova/internal/search/domain/entity"
	"vicomova/pkg/infrastructure/mysql"
)

type VideoMetaRepository struct {
	db *mysql.Client
}

func NewVideoMetaRepository(db *mysql.Client) *VideoMetaRepository {
	return &VideoMetaRepository{db: db}
}

func (r *VideoMetaRepository) GetByIDs(ctx context.Context, ids []int64) ([]*entity.VideoSearchResult, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var pos []VideoPO
	if err := r.db.WithContext(ctx).Where("id IN ? AND status = 2", ids).Find(&pos).Error; err != nil {
		return nil, err
	}

	results := make([]*entity.VideoSearchResult, len(pos))
	for i, po := range pos {
		results[i] = &entity.VideoSearchResult{
			VideoID:      po.ID,
			Title:        po.Title,
			Description:  po.Description,
			CoverURL:     po.CoverURL,
			ViewCount:    po.ViewCount,
			LikeCount:    po.LikeCount,
			CommentCount: po.CommentCount,
		}
	}
	return results, nil
}

type VideoPO struct {
	ID           int64  `gorm:"column:id"`
	Title        string `gorm:"column:title"`
	Description  string `gorm:"column:description"`
	CoverURL     string `gorm:"column:cover_url"`
	ViewCount    int64  `gorm:"column:view_count"`
	LikeCount    int64  `gorm:"column:like_count"`
	CommentCount int64  `gorm:"column:comment_count"`
	Status       int8   `gorm:"column:status"`
}

func (VideoPO) TableName() string {
	return "videos"
}
