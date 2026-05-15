package repository

import (
	"context"

	"vicomova/internal/commerce/domain/entity"
)

type PointsWalletRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*entity.PointsWallet, error)
	Create(ctx context.Context, wallet *entity.PointsWallet) error
	Update(ctx context.Context, wallet *entity.PointsWallet) error
	// UpdateVersion 乐观锁更新，返回是否成功
	UpdateVersion(ctx context.Context, userID int64, oldVersion int64, newVersion int64, delta int64) (bool, error)
}

type SignInRecordRepository interface {
	GetTodaySignIn(ctx context.Context, userID int64) (*entity.SignInRecord, error)
	GetLastSignIn(ctx context.Context, userID int64) (*entity.SignInRecord, error)
	Create(ctx context.Context, record *entity.SignInRecord) error
	ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*entity.SignInRecord, bool, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *entity.Transaction) error
	GetByID(ctx context.Context, id int64) (*entity.Transaction, error)
	ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*entity.Transaction, bool, error)
}

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
	List(ctx context.Context, productType entity.ProductType, cursor string, limit int) ([]*entity.Product, string, bool, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error)
	ListFlashSaleProducts(ctx context.Context) ([]*entity.Product, error)
	ListActive(ctx context.Context) ([]*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	UpdateStock(ctx context.Context, id int64, delta int) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id int64) (*entity.Order, error)
	GetByUserAndProduct(ctx context.Context, userID, productID int64) (*entity.Order, error)
	UpdateStatus(ctx context.Context, id int64, status entity.OrderStatus) error
	ListByUser(ctx context.Context, userID int64, cursor string, limit int) ([]*entity.Order, bool, error)
}

type MembershipRepository interface {
	GetActiveMembership(ctx context.Context, userID int64) (*entity.Membership, error)
	GetByID(ctx context.Context, id int64) (*entity.Membership, error)
	Create(ctx context.Context, m *entity.Membership) error
	Update(ctx context.Context, m *entity.Membership) error
	ListByUser(ctx context.Context, userID int64) ([]*entity.Membership, error)
}

type VideoPermissionRepository interface {
	Check(ctx context.Context, userID, videoID int64) (bool, error)
	Get(ctx context.Context, userID, videoID int64) (*entity.VideoPermission, error)
	Create(ctx context.Context, p *entity.VideoPermission) error
	ListByUser(ctx context.Context, userID int64) ([]*entity.VideoPermission, error)
	ListByVideo(ctx context.Context, videoID int64) ([]*entity.VideoPermission, error)
}

type FlashSaleRepository interface {
	GetActiveFlashSale(ctx context.Context) (*entity.FlashSale, error)
	GetByID(ctx context.Context, id int64) (*entity.FlashSale, error)
	Create(ctx context.Context, fs *entity.FlashSale) error
	ListActive(ctx context.Context) ([]*entity.FlashSale, error)
}

type FlashSaleStockRepository interface {
	GetByFlashSaleID(ctx context.Context, flashSaleID int64) ([]*entity.FlashSaleStock, error)
	GetStock(ctx context.Context, flashSaleID, productID int64) (*entity.FlashSaleStock, error)
	Create(ctx context.Context, stock *entity.FlashSaleStock) error
	// DecrementStock 原子扣减库存，返回是否成功
	DecrementStock(ctx context.Context, flashSaleID, productID int64) (bool, error)
}

type LotteryRepository interface {
	GetActiveLottery(ctx context.Context) (*entity.LotteryDraw, error)
	GetByID(ctx context.Context, id int64) (*entity.LotteryDraw, error)
	Create(ctx context.Context, l *entity.LotteryDraw) error
	ListActive(ctx context.Context) ([]*entity.LotteryDraw, error)
}

type LotteryRecordRepository interface {
	RecordParticipation(ctx context.Context, record *entity.LotteryRecord) error
	GetUserParticipation(ctx context.Context, userID, lotteryID int64) (*entity.LotteryRecord, error)
	CountByUser(ctx context.Context, userID, lotteryID int64) (int, error)
}
