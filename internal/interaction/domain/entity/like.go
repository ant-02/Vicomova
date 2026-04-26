package entity

import "time"

type Like struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    TargetType string   `json:"target_type"` // "video" or "comment"
    TargetID  int64     `json:"target_id"`
    CreatedAt time.Time `json:"created_at"`
}

const (
    TargetTypeVideo   = "video"
    TargetTypeComment = "comment"
)
