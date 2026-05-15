package mysql

import (
	"time"

	"gorm.io/gorm"

	constants "vicomova/pkg/constants"
)

type PointsWalletPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"uniqueIndex;not null"`
	Points    int64          `gorm:"not null;default:0"`
	Version   int64          `gorm:"not null;default:0"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (PointsWalletPO) TableName() string {
	return constants.TablePointsWallets
}

type SignInRecordPO struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	UserID          int64          `gorm:"index:idx_user_date,unique"`
	SignDate        time.Time      `gorm:"type:date;not null"`
	ConsecutiveDays int            `gorm:"not null;default:1"`
	BonusPoints     int64          `gorm:"not null;default:0"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (SignInRecordPO) TableName() string {
	return constants.TableSignInRecords
}

type TransactionPO struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	UserID       int64          `gorm:"index;not null"`
	Type         string         `gorm:"size:32;not null"`
	Amount       int64          `gorm:"not null"`
	BalanceAfter int64          `gorm:"not null"`
	Memo         string         `gorm:"size:255"`
	Status       string         `gorm:"size:32;not null"`
	ReferenceID  string         `gorm:"size:64"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (TransactionPO) TableName() string {
	return constants.TableTransactions
}

type ProductPO struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	Name            string         `gorm:"size:128;not null"`
	Description     string         `gorm:"size:512"`
	ProductType     int8           `gorm:"not null;index"`
	Price           int64          `gorm:"not null;default:0"`
	WalletPrice     float64        `gorm:"type:decimal(10,2);default:0"`
	Stock           int            `gorm:"not null;default:0"`
	VideoID         int64          `gorm:"index"`
	MembershipLevel int            `gorm:"default:0"`
	IsActive        bool           `gorm:"not null;default:true"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (ProductPO) TableName() string {
	return constants.TableProducts
}

type OrderPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"index;not null"`
	ProductID int64          `gorm:"index;not null"`
	PayType   int8           `gorm:"not null"`
	PricePaid int64          `gorm:"not null;default:0"`
	Status    int8           `gorm:"not null;default:1"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	PaidAt    *time.Time     `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (OrderPO) TableName() string {
	return constants.TableOrders
}

type MembershipPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"uniqueIndex;not null"`
	Level     int8           `gorm:"not null;default:1"`
	StartAt   time.Time      `gorm:"not null"`
	ExpireAt  time.Time      `gorm:"not null;index"`
	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (MembershipPO) TableName() string {
	return constants.TableMemberships
}

type VideoPermissionPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"index:idx_user_video,unique;not null"`
	VideoID   int64          `gorm:"index:idx_user_video,unique;not null"`
	Type      string         `gorm:"size:32;not null"` // vip/purchased
	ExpireAt  *time.Time     `gorm:"index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (VideoPermissionPO) TableName() string {
	return constants.TableVideoPermissions
}

type FlashSalePO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"size:128;not null"`
	StartTime time.Time      `gorm:"not null;index"`
	EndTime   time.Time      `gorm:"not null;index"`
	Status    string         `gorm:"size:32;not null;default:pending"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (FlashSalePO) TableName() string {
	return constants.TableFlashSales
}

type FlashSaleStockPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	FlashSaleID int64          `gorm:"index:idx_flash_product,unique;not null"`
	ProductID   int64          `gorm:"index:idx_flash_product,unique;not null"`
	Stock       int            `gorm:"not null;default:0"`
	RemainStock int            `gorm:"not null;default:0"`
	FlashPrice  int64          `gorm:"not null;default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (FlashSaleStockPO) TableName() string {
	return constants.TableFlashSaleStocks
}

type LotteryDrawPO struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"size:128;not null"`
	Description   string         `gorm:"size:512"`
	EntryPoints   int64          `gorm:"not null;default:0"`
	StartTime     time.Time      `gorm:"not null;index"`
	EndTime       time.Time      `gorm:"not null;index"`
	TotalTickets  int            `gorm:"not null;default:0"`
	RemainTickets int            `gorm:"not null;default:0"`
	Prizes        string         `gorm:"type:text"` // JSON
	Status        string         `gorm:"size:32;not null;default:pending"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (LotteryDrawPO) TableName() string {
	return constants.TableLotteryDraws
}

type LotteryRecordPO struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	UserID     int64          `gorm:"index:idx_user_lottery,unique;not null"`
	LotteryID  int64          `gorm:"index:idx_user_lottery,unique;not null"`
	Prize      string         `gorm:"size:128"`
	IsWin      bool           `gorm:"not null;default:false"`
	DrawNumber int            `gorm:"not null;default:1"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (LotteryRecordPO) TableName() string {
	return constants.TableLotteryRecords
}
