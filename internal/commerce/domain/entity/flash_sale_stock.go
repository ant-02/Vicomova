package entity

import (
	"time"

	"gorm.io/gorm"
)

// FlashSaleStock 秒杀库存
type FlashSaleStock struct {
	ID          int64          `json:"id"`
	FlashSaleID int64          `json:"flash_sale_id"`
	ProductID   int64          `json:"product_id"`
	Stock       int            `json:"stock"`        // 总库存
	RemainStock int            `json:"remain_stock"` // 剩余库存
	FlashPrice  int64          `json:"flash_price"`  // 秒杀价格（积分）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func NewFlashSaleStock(flashSaleID, productID int64, stock int, flashPrice int64) *FlashSaleStock {
	now := time.Now()
	return &FlashSaleStock{
		FlashSaleID: flashSaleID,
		ProductID:   productID,
		Stock:       stock,
		RemainStock: stock,
		FlashPrice:  flashPrice,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *FlashSaleStock) Decrement() bool {
	if s.RemainStock <= 0 {
		return false
	}
	s.RemainStock--
	s.UpdatedAt = time.Now()
	return true
}

func (s *FlashSaleStock) Available() bool {
	return s.RemainStock > 0
}
