package valueobject

import (
	"time"
)

type EmailCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

func NewEmailCode(email, code string, expiresAt time.Time) *EmailCode {
	return &EmailCode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
	}
}

func (e *EmailCode) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

func (e *EmailCode) IsValid() bool {
	return !e.IsExpired() && len(e.Code) > 0
}
