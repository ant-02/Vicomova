package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	VideoID   int64          `json:"video_id"`
	ParentID  int64          `json:"parent_id"` // 0: root comment, >0: reply
	Content   string         `json:"content"`
	LikeCount int64          `json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
