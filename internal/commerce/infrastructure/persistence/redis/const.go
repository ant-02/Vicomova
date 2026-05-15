package redis

import (
	"time"
)

const (
	PointsWalletPrefix    = "commerce:wallet:"      // + userID
	SignInStreakPrefix    = "commerce:signin:"      // + userID
	FlashSaleStockPrefix  = "commerce:flash:stock:" // + flashSaleID:productID
	FlashSaleConfigPrefix = "commerce:flash:cfg:"   // + flashSaleID
	LotteryOddsPrefix     = "commerce:lottery:"     // + lotteryID

	WalletCacheTTL = 24 * time.Hour
	SignInCacheTTL = 25 * time.Hour
	FlashSaleTTL   = 5 * time.Minute
	LotteryTTL     = 1 * time.Hour
)
