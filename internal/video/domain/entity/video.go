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
	HotScore     float64             `json:"hot_score"`
	Status       videoVO.VideoStatus `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `json:"-"`
}

func (v *Video) IsPublished() bool {
	return v.Status == videoVO.VideoStatusPublished
}

// HotVideoScore 视频热度分数（用于缓存返回）
type HotVideoScore struct {
	VideoID int64
	Score   float64
}

// HotVideoMeta 热门视频元数据（精简版，用于缓存和列表展示）
type HotVideoMeta struct {
	VideoID      int64
	Title        string
	CoverURL     string
	UserID       int64
	Duration     int
	ViewCount    int64
	CommentCount int64
}

// FromVideo 从 Video 转换
func (m *HotVideoMeta) FromVideo(v *Video) {
	m.VideoID = v.ID
	m.Title = v.Title
	m.CoverURL = v.CoverURL
	m.UserID = v.UserID
	m.Duration = v.Duration
	m.ViewCount = v.ViewCount
	m.CommentCount = v.CommentCount
}
