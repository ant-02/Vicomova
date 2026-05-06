package entity

import "time"

type Favorite struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	VideoID   int64     `json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
}
