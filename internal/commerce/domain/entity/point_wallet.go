package entity

import (
	"time"

	"gorm.io/gorm"
)

// PointsWallet 用户积分钱包
type PointsWallet struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Points    int64          `json:"points"`  // 当前积分
	Version   int64          `json:"version"` // 乐观锁版本号
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func NewPointsWallet(userID int64) *PointsWallet {
	now := time.Now()
	return &PointsWallet{
		UserID:    userID,
		Points:    0,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *PointsWallet) AddPoints(amount int64) {
	w.Points += amount
	w.UpdatedAt = time.Now()
}

func (w *PointsWallet) DeductPoints(amount int64) bool {
	if w.Points < amount {
		return false
	}
	w.Points -= amount
	w.UpdatedAt = time.Now()
	return true
}
