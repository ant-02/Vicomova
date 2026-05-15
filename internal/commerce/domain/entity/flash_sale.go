package entity

import (
	"time"

	"gorm.io/gorm"
)

// FlashSale 秒杀场次
type FlashSale struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Status    string         `json:"status"` // pending/active/ended
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

const (
	FlashSaleStatusPending = "pending"
	FlashSaleStatusActive  = "active"
	FlashSaleStatusEnded   = "ended"
)

func NewFlashSale(name string, startTime, endTime time.Time) *FlashSale {
	now := time.Now()
	status := FlashSaleStatusPending
	if now.After(startTime) && now.Before(endTime) {
		status = FlashSaleStatusActive
	} else if now.After(endTime) {
		status = FlashSaleStatusEnded
	}
	return &FlashSale{
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (f *FlashSale) IsActive() bool {
	now := time.Now()
	return f.Status == FlashSaleStatusActive && now.After(f.StartTime) && now.Before(f.EndTime)
}
