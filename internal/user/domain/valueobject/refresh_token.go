package valueobject

import (
	"time"
)

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

func NewRefreshToken(userID int64, token string, expiry time.Duration) *RefreshToken {
	return &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(expiry),
		CreatedAt: time.Now(),
		Revoked:   false,
	}
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.Revoked
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked() && !rt.IsExpired()
}