package entity

import (
	"time"

	"gorm.io/gorm"
)

// LotteryRecord 用户抽奖记录
type LotteryRecord struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"user_id"`
	LotteryID  int64          `json:"lottery_id"`
	Prize      string         `json:"prize"` // 获得的奖品
	IsWin      bool           `json:"is_win"`
	DrawNumber int            `json:"draw_number"` // 第几次抽奖（同一个活动）
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

func NewLotteryRecord(userID, lotteryID int64, prize string, isWin bool, drawNumber int) *LotteryRecord {
	now := time.Now()
	return &LotteryRecord{
		UserID:     userID,
		LotteryID:  lotteryID,
		Prize:      prize,
		IsWin:      isWin,
		DrawNumber: drawNumber,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
