package kafka

import (
	constants "vicomova/pkg/constants"
)

const (
	TopicUserSignIn          = constants.KafkaTopicUserSignIn
	TopicCommerceSignIn      = constants.KafkaTopicCommerceSignIn
	TopicCommerceTransaction = constants.KafkaTopicCommerceTransaction
	TopicCommerceFlashSale   = constants.KafkaTopicCommerceFlashSale
	TopicCommerceLottery     = constants.KafkaTopicCommerceLottery
)

// SignInEvent 用户签到事件
type SignInEvent struct {
	UserID          int64  `json:"user_id"`
	SignDate        string `json:"sign_date"` // 格式: 2006-01-02
	ConsecutiveDays int    `json:"consecutive_days"`
}

// TransactionEvent 钱包变动事件
type TransactionEvent struct {
	TransactionID int64  `json:"transaction_id"`
	UserID        int64  `json:"user_id"`
	Type          string `json:"type"` // earn/spend/recharge/refund
	Amount        int64  `json:"amount"`
	BalanceAfter  int64  `json:"balance_after"`
}

// FlashSaleEvent 秒杀事件
type FlashSaleEvent struct {
	SaleID         int64 `json:"sale_id"`
	UserID         int64 `json:"user_id"`
	ProductID      int64 `json:"product_id"`
	StockRemaining int   `json:"stock_remaining"`
	Success        bool  `json:"success"`
}

// LotteryEvent 抽奖事件
type LotteryEvent struct {
	DrawID int64  `json:"draw_id"`
	UserID int64  `json:"user_id"`
	Prize  string `json:"prize"`
	Won    bool   `json:"won"`
}
