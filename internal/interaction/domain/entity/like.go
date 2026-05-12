package entity

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"user_id"`
	TargetType string         `json:"target_type"` // "video" or "comment"
	TargetID   int64          `json:"target_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

const (
	TargetTypeVideo   = "video"
	TargetTypeComment = "comment"
)
