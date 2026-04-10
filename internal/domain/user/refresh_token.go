package user

import "time"

type RefreshToken struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
