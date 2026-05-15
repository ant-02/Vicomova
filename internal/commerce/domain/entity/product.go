package entity

import (
	"time"

	"gorm.io/gorm"
)

// ProductType 商品类型
type ProductType int8

const (
	ProductTypePointsItem  ProductType = 1 // 积分商品
	ProductTypeVIP         ProductType = 2 // VIP会员
	ProductTypeVideoAccess ProductType = 3 // 视频权限
	ProductTypeFlashSale   ProductType = 4 // 秒杀商品
)

// Product 商品
type Product struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ProductType     ProductType    `json:"product_type"`
	Price           int64          `json:"price"`            // 价格（积分）
	WalletPrice     float64        `json:"wallet_price"`     // 价格（元）
	Stock           int            `json:"stock"`            // 库存
	VideoID         int64          `json:"video_id"`         // 视频ID（视频权限类型）
	MembershipLevel int            `json:"membership_level"` // 会员等级（VIP类型）
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-"`
}

func NewProduct(name, description string, productType ProductType, price int64, walletPrice float64, stock int) *Product {
	now := time.Now()
	return &Product{
		Name:        name,
		Description: description,
		ProductType: productType,
		Price:       price,
		WalletPrice: walletPrice,
		Stock:       stock,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) IsValid() bool {
	return p.IsActive && p.Stock > 0
}
