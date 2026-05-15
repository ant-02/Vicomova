package entity

import (
	"time"

	"gorm.io/gorm"
)

// Transaction 积分/钱包变动流水
type Transaction struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	Type         string         `json:"type"`          // earn/spend/recharge/refund
	Amount       int64          `json:"amount"`        // 变动积分数量
	BalanceAfter int64          `json:"balance_after"` // 变动后余额
	Memo         string         `json:"memo"`          // 备注
	Status       string         `json:"status"`        // pending/completed/failed
	ReferenceID  string         `json:"reference_id"`  // 关联业务ID，如order_id
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

func NewTransaction(userID int64, txType string, amount, balanceAfter int64, memo, status, referenceID string) *Transaction {
	now := time.Now()
	return &Transaction{
		UserID:       userID,
		Type:         txType,
		Amount:       amount,
		BalanceAfter: balanceAfter,
		Memo:         memo,
		Status:       status,
		ReferenceID:  referenceID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
