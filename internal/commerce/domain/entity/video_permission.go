package entity

import (
	"time"

	"gorm.io/gorm"
)

// VideoPermission 用户视频购买记录
type VideoPermission struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	VideoID   int64          `json:"video_id"`
	Type      string         `json:"type"`      // vip/purchased
	ExpireAt  *time.Time     `json:"expire_at"` // nil表示永久
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func NewVideoPermission(userID, videoID int64, permType string, expireAt *time.Time) *VideoPermission {
	now := time.Now()
	return &VideoPermission{
		UserID:    userID,
		VideoID:   videoID,
		Type:      permType,
		ExpireAt:  expireAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *VideoPermission) IsValid() bool {
	if p.ExpireAt == nil {
		return true
	}
	return time.Now().Before(*p.ExpireAt)
}
