package entity

import (
	"time"

	"gorm.io/gorm"
)

// MembershipLevel 会员等级
type MembershipLevel int8

const (
	MembershipOrdinary MembershipLevel = 1 // 普通会员
	MembershipSilver   MembershipLevel = 2 // 银牌会员
	MembershipGold     MembershipLevel = 3 // 金牌会员
)

// Membership 会员记录
type Membership struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Level     MembershipLevel `json:"level"`
	StartAt   time.Time       `json:"start_at"`
	ExpireAt  time.Time       `json:"expire_at"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"-"`
}

func NewMembership(userID int64, level MembershipLevel, durationDays int) *Membership {
	now := time.Now()
	return &Membership{
		UserID:    userID,
		Level:     level,
		StartAt:   now,
		ExpireAt:  now.AddDate(0, 0, durationDays),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *Membership) IsValid() bool {
	return m.IsActive && time.Now().Before(m.ExpireAt)
}

func (m *Membership) CanAccess(requiredLevel MembershipLevel) bool {
	return m.IsValid() && m.Level >= requiredLevel
}

func (m *Membership) Extend(days int) {
	m.ExpireAt = m.ExpireAt.AddDate(0, 0, days)
	m.UpdatedAt = time.Now()
}
