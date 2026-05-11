package entity

import (
	"time"

	videoVO "vicomova/internal/video/domain/valueobject"

	"gorm.io/gorm"
)

type Video struct {
	ID           int64               `json:"id"`
	UserID       int64               `json:"user_id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	CoverURL     string              `json:"cover_url"`
	VideoURL     string              `json:"video_url"`
	CategoryID   int                 `json:"category_id"`
	ViewCount    int64               `json:"view_count"`
	LikeCount    int64               `json:"like_count"`
	CommentCount int64               `json:"comment_count"`
	Duration     int                 `json:"duration"`
	Status       videoVO.VideoStatus `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `json:"-"`
}

func (v *Video) IsPublished() bool {
	return v.Status == videoVO.VideoStatusPublished
}
