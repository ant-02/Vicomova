package entity

import (
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	ID          int64          `json:"id"`
	FollowerID  int64          `json:"follower_id"`  // 关注者
	FollowingID int64          `json:"following_id"` // 被关注者
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
