package command

// PointsWalletCommand 积分钱包命令
type RechargePointsCommand struct {
	UserID    int64
	Amount    int64
	PaymentID string
}

type DeductPointsCommand struct {
	UserID int64
	Amount int64
	Memo   string
}

type RechargePointsResult struct {
	Wallet    *PointsWalletDTO
	NewPoints int64
}

type PointsWalletDTO struct {
	UserID  int64 `json:"user_id"`
	Points  int64 `json:"points"`
	Version int64 `json:"version"`
}

// SignInCommand 签到命令
type ProcessSignInCommand struct {
	UserID          int64
	SignDate        string
	ConsecutiveDays int
}

type ProcessSignInResult struct {
	Success         bool   `json:"success"`
	PointsEarned    int64  `json:"points_earned"`
	ConsecutiveDays int    `json:"consecutive_days"`
	Message         string `json:"message"`
}

// OrderCommand 订单命令
type CreateOrderCommand struct {
	UserID    int64
	ProductID int64
	PayType   int // 1=points, 2=wallet
}

type PayOrderCommand struct {
	OrderID int64
	UserID  int64
}

type CreateOrderResult struct {
	Order *OrderDTO
}

type OrderDTO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	PayType   int    `json:"pay_type"`
	PricePaid int64  `json:"price_paid"`
	Status    string `json:"status"`
}

// MembershipCommand 会员命令
type BuyMembershipCommand struct {
	UserID  int64
	Level   int
	PayType int
}

type BuyMembershipResult struct {
	Membership *MembershipDTO
}

type MembershipDTO struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"user_id"`
	Level    int   `json:"level"`
	StartAt  int64 `json:"start_at"`
	ExpireAt int64 `json:"expire_at"`
	IsActive bool  `json:"is_active"`
}

// VideoAccessCommand 视频权限命令
type BuyVideoAccessCommand struct {
	UserID  int64
	VideoID int64
	PayType int
}

type BuyVideoAccessResult struct {
	Success  bool  `json:"success"`
	ExpireAt int64 `json:"expire_at,omitempty"`
}

// FlashSaleCommand 秒杀命令
type FlashSalePurchaseCommand struct {
	UserID      int64
	FlashSaleID int64
	ProductID   int64
	PayType     int
}

type FlashSalePurchaseResult struct {
	Success bool   `json:"success"`
	OrderID int64  `json:"order_id,omitempty"`
	Message string `json:"message"`
}

// LotteryCommand 抽奖命令
type ParticipateLotteryCommand struct {
	UserID    int64
	LotteryID int64
}

type ParticipateLotteryResult struct {
	Success bool   `json:"success"`
	Won     bool   `json:"won"`
	Prize   string `json:"prize"`
	Message string `json:"message"`
}
