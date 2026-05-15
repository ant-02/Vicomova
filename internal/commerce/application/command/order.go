package command

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	commerceVO "vicomova/internal/commerce/domain/valueobject"
	"vicomova/pkg/log"
)

type OrderCommandService struct {
	orderRepo      commerceRepo.OrderRepository
	productRepo    commerceRepo.ProductRepository
	walletRepo     commerceRepo.PointsWalletRepository
	txnRepo        commerceRepo.TransactionRepository
	membershipRepo commerceRepo.MembershipRepository
	videoPermRepo  commerceRepo.VideoPermissionRepository
}

func NewOrderCommandService(
	orderRepo commerceRepo.OrderRepository,
	productRepo commerceRepo.ProductRepository,
	walletRepo commerceRepo.PointsWalletRepository,
	txnRepo commerceRepo.TransactionRepository,
	membershipRepo commerceRepo.MembershipRepository,
	videoPermRepo commerceRepo.VideoPermissionRepository,
) *OrderCommandService {
	return &OrderCommandService{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		walletRepo:     walletRepo,
		txnRepo:        txnRepo,
		membershipRepo: membershipRepo,
		videoPermRepo:  videoPermRepo,
	}
}

func (s *OrderCommandService) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) (*CreateOrderResult, error) {
	// 检查商品是否存在
	product, err := s.productRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil || !product.IsValid() {
		return nil, ErrProductNotAvailable
	}

	// 检查是否已有未完成订单
	existing, err := s.orderRepo.GetByUserAndProduct(ctx, cmd.UserID, cmd.ProductID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrOrderAlreadyExists
	}

	// 创建订单
	order := commerceEntity.NewOrder(cmd.UserID, cmd.ProductID, commerceEntity.PayType(cmd.PayType), product.Price)
	if err := s.orderRepo.Create(ctx, order); err != nil {
		log.Error.Printf("OrderCommandService.CreateOrder: failed: %v", err)
		return nil, err
	}

	return &CreateOrderResult{
		Order: &OrderDTO{
			ID:        order.ID,
			UserID:    order.UserID,
			ProductID: order.ProductID,
			PayType:   int(order.PayType),
			PricePaid: order.PricePaid,
			Status:    commerceVO.OrderStatus(order.Status).String(),
		},
	}, nil
}

func (s *OrderCommandService) PayOrder(ctx context.Context, cmd *PayOrderCommand) error {
	order, err := s.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrOrderNotFound
	}
	if order.UserID != cmd.UserID {
		return ErrOrderNotFound
	}
	if !order.Status.CanPay() {
		return ErrOrderCannotPay
	}

	product, err := s.productRepo.GetByID(ctx, order.ProductID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotAvailable
	}

	// 扣减积分或处理钱包支付
	if order.PayType == commerceEntity.PayTypePoints {
		wallet, err := s.walletRepo.GetByUserID(ctx, cmd.UserID)
		if err != nil {
			return err
		}
		if wallet == nil || wallet.Points < order.PricePaid {
			return ErrInsufficientPoints
		}

		oldVersion := wallet.Version
		success, err := s.walletRepo.UpdateVersion(ctx, cmd.UserID, oldVersion, oldVersion+1, -order.PricePaid)
		if err != nil {
			return err
		}
		if !success {
			return ErrInsufficientPoints
		}
	}

	// 标记订单已支付
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, commerceEntity.OrderStatusPaid); err != nil {
		return err
	}

	// 发放商品
	if err := s.deliverProduct(ctx, cmd.UserID, product); err != nil {
		log.Error.Printf("OrderCommandService.PayOrder: deliver product failed: %v", err)
		// 回滚订单状态
		s.orderRepo.UpdateStatus(ctx, order.ID, commerceEntity.OrderStatusPending)
		return err
	}

	// 标记订单完成
	if err := s.orderRepo.UpdateStatus(ctx, order.ID, commerceEntity.OrderStatusCompleted); err != nil {
		log.Warn.Printf("OrderCommandService.PayOrder: update status failed: %v", err)
	}

	// 扣减库存
	if err := s.productRepo.UpdateStock(ctx, product.ID, -1); err != nil {
		log.Warn.Printf("OrderCommandService.PayOrder: update stock failed: %v", err)
	}

	return nil
}

func (s *OrderCommandService) deliverProduct(ctx context.Context, userID int64, product *commerceEntity.Product) error {
	switch product.ProductType {
	case commerceEntity.ProductTypeVIP:
		// 发放会员
		membership := commerceEntity.NewMembership(userID, commerceEntity.MembershipLevel(product.MembershipLevel), 30) // 默认30天
		return s.membershipRepo.Create(ctx, membership)
	case commerceEntity.ProductTypeVideoAccess:
		// 发放视频权限
		perm := commerceEntity.NewVideoPermission(userID, product.VideoID, "purchased", nil)
		return s.videoPermRepo.Create(ctx, perm)
	default:
		// 其他商品类型暂不处理
		return nil
	}
}

var ErrProductNotAvailable = &CommerceError{Code: "PRODUCT_NOT_AVAILABLE", Message: "商品不可用"}
var ErrOrderAlreadyExists = &CommerceError{Code: "ORDER_ALREADY_EXISTS", Message: "订单已存在"}
var ErrOrderNotFound = &CommerceError{Code: "ORDER_NOT_FOUND", Message: "订单不存在"}
var ErrOrderCannotPay = &CommerceError{Code: "ORDER_CANNOT_PAY", Message: "订单状态不可支付"}
