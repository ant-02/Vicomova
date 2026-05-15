package query

// PointsWalletQuery 积分钱包查询
type GetWalletQuery struct {
	UserID int64
}

type GetWalletResult struct {
	Wallet *WalletDTO
}

type WalletDTO struct {
	UserID  int64 `json:"user_id"`
	Points  int64 `json:"points"`
	Version int64 `json:"version"`
}

// SignInQuery 签到查询
type GetSignInInfoQuery struct {
	UserID int64
}

type GetSignInInfoResult struct {
	Info *SignInInfoDTO
}

type SignInInfoDTO struct {
	UserID          int64 `json:"user_id"`
	ConsecutiveDays int   `json:"consecutive_days"`
	LastSignInAt    int64 `json:"last_sign_in_at"`
	TodayBonus      int64 `json:"today_bonus"`
	SignedInToday   bool  `json:"signed_in_today"`
}

// ProductQuery 商品查询
type ListProductsQuery struct {
	ProductType int
	Cursor      string
	Limit       int
}

type ListProductsResult struct {
	Products   []*ProductDTO
	NextCursor string
	HasMore    bool
}

type ProductDTO struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	ProductType     int     `json:"product_type"`
	Price           int64   `json:"price"`
	WalletPrice     float64 `json:"wallet_price"`
	Stock           int     `json:"stock"`
	VideoID         int64   `json:"video_id"`
	MembershipLevel int     `json:"membership_level"`
	IsActive        bool    `json:"is_active"`
}

// OrderQuery 订单查询
type ListOrdersQuery struct {
	UserID int64
	Cursor string
	Limit  int
}

type ListOrdersResult struct {
	Orders     []*OrderDTO
	HasMore    bool
	NextCursor string
}

type OrderDTO struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	PayType   int   `json:"pay_type"`
	PricePaid int64 `json:"price_paid"`
	Status    int   `json:"status"` // OrderStatus as int for conversion
}

// MembershipQuery 会员查询
type GetMembershipQuery struct {
	UserID int64
}

type GetMembershipResult struct {
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

// FlashSaleQuery 秒杀查询
type ListFlashSalesQuery struct{}

type ListFlashSalesResult struct {
	FlashSales []*FlashSaleDTO
}

type FlashSaleDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Status    string `json:"status"`
}

// LotteryQuery 抽奖查询
type ListLotteriesQuery struct{}

type ListLotteriesResult struct {
	Lotteries []*LotteryDTO
}

type LotteryDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	EntryPoints   int64  `json:"entry_points"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	RemainTickets int    `json:"remain_tickets"`
	Status        string `json:"status"`
}
