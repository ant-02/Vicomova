package entity

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态
type OrderStatus int8

const (
	OrderStatusPending   OrderStatus = 1 // 待支付
	OrderStatusPaid      OrderStatus = 2 // 已支付
	OrderStatusCompleted OrderStatus = 3 // 已完成
	OrderStatusCancelled OrderStatus = 4 // 已取消
	OrderStatusRefunded  OrderStatus = 5 // 已退款
)

// PayType 支付类型
type PayType int8

const (
	PayTypePoints PayType = 1 // 积分支付
	PayTypeWallet PayType = 2 // 钱包支付
)

// Order 订单
type Order struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ProductID int64          `json:"product_id"`
	PayType   PayType        `json:"pay_type"`
	PricePaid int64          `json:"price_paid"` // 实际支付积分
	Status    OrderStatus    `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	PaidAt    *time.Time     `json:"paid_at"` // 支付时间
	DeletedAt gorm.DeletedAt `json:"-"`
}

func NewOrder(userID, productID int64, payType PayType, pricePaid int64) *Order {
	now := time.Now()
	return &Order{
		UserID:    userID,
		ProductID: productID,
		PayType:   payType,
		PricePaid: pricePaid,
		Status:    OrderStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *Order) MarkPaid() {
	now := time.Now()
	o.Status = OrderStatusPaid
	o.PaidAt = &now
	o.UpdatedAt = now
}

func (o *Order) MarkCompleted() {
	o.Status = OrderStatusCompleted
	o.UpdatedAt = time.Now()
}

func (o *Order) MarkCancelled() {
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
}

// CanPay 判断订单是否可以支付
func (s OrderStatus) CanPay() bool {
	return s == OrderStatusPending
}

// CanCancel 判断订单是否可以取消
func (s OrderStatus) CanCancel() bool {
	return s == OrderStatusPending
}
