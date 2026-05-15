package entity

import (
	"time"

	"gorm.io/gorm"
)

// SignInRecord 用户签到记录
type SignInRecord struct {
	ID              int64          `json:"id"`
	UserID          int64          `json:"user_id"`
	SignDate        time.Time      `json:"sign_date"`        // 签到日期（只到日期）
	ConsecutiveDays int            `json:"consecutive_days"` // 当日连续签到天数
	BonusPoints     int64          `json:"bonus_points"`     // 本次获得积分
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-"`
}

func NewSignInRecord(userID int64, signDate time.Time, consecutiveDays int, bonusPoints int64) *SignInRecord {
	now := time.Now()
	return &SignInRecord{
		UserID:          userID,
		SignDate:        signDate,
		ConsecutiveDays: consecutiveDays,
		BonusPoints:     bonusPoints,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
