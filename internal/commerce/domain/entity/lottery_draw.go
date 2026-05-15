package entity

import (
	"time"

	"gorm.io/gorm"
)

// LotteryDraw 抽奖活动
type LotteryDraw struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	EntryPoints   int64          `json:"entry_points"` // 参与所需积分
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time"`
	TotalTickets  int            `json:"total_tickets"` // 总票数
	RemainTickets int            `json:"remain_tickets"`
	Prizes        string         `json:"prizes"` // JSON格式的奖品列表
	Status        string         `json:"status"` // pending/active/ended
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

const (
	LotteryStatusPending = "pending"
	LotteryStatusActive  = "active"
	LotteryStatusEnded   = "ended"
)

func NewLotteryDraw(name, description string, entryPoints int64, startTime, endTime time.Time, totalTickets int, prizes string) *LotteryDraw {
	now := time.Now()
	status := LotteryStatusPending
	if now.After(startTime) && now.Before(endTime) {
		status = LotteryStatusActive
	} else if now.After(endTime) {
		status = LotteryStatusEnded
	}
	return &LotteryDraw{
		Name:          name,
		Description:   description,
		EntryPoints:   entryPoints,
		StartTime:     startTime,
		EndTime:       endTime,
		TotalTickets:  totalTickets,
		RemainTickets: totalTickets,
		Prizes:        prizes,
		Status:        status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (l *LotteryDraw) IsActive() bool {
	now := time.Now()
	return l.Status == LotteryStatusActive && now.After(l.StartTime) && now.Before(l.EndTime) && l.RemainTickets > 0
}

func (l *LotteryDraw) HasTickets() bool {
	return l.RemainTickets > 0
}
