package entity

import (
	"time"
)

type OfflineMessage struct {
	ID         int64     `json:"id"`
	ToUserID   int64     `json:"to_user_id"`
	FromUserID int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `json:"is_read"`
}
