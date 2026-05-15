package command

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	redisCache "vicomova/internal/commerce/infrastructure/persistence/redis"
	"vicomova/pkg/log"
)

type FlashSaleCommandService struct {
	flashSaleRepo  commerceRepo.FlashSaleRepository
	stockRepo      commerceRepo.FlashSaleStockRepository
	productRepo    commerceRepo.ProductRepository
	orderRepo      commerceRepo.OrderRepository
	flashSaleCache *redisCache.FlashSaleCache
}

func NewFlashSaleCommandService(
	flashSaleRepo commerceRepo.FlashSaleRepository,
	stockRepo commerceRepo.FlashSaleStockRepository,
	productRepo commerceRepo.ProductRepository,
	orderRepo commerceRepo.OrderRepository,
	flashSaleCache *redisCache.FlashSaleCache,
) *FlashSaleCommandService {
	return &FlashSaleCommandService{
		flashSaleRepo:  flashSaleRepo,
		stockRepo:      stockRepo,
		productRepo:    productRepo,
		orderRepo:      orderRepo,
		flashSaleCache: flashSaleCache,
	}
}

func (s *FlashSaleCommandService) Purchase(ctx context.Context, cmd *FlashSalePurchaseCommand) (*FlashSalePurchaseResult, error) {
	// 检查秒杀场次
	flashSale, err := s.flashSaleRepo.GetByID(ctx, cmd.FlashSaleID)
	if err != nil {
		return nil, err
	}
	if flashSale == nil {
		return nil, ErrFlashSaleNotFound
	}
	if !flashSale.IsActive() {
		return nil, ErrFlashSaleNotActive
	}

	// 获取秒杀库存
	stock, err := s.stockRepo.GetStock(ctx, cmd.FlashSaleID, cmd.ProductID)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, ErrFlashSaleProductNotFound
	}
	if !stock.Available() {
		return nil, ErrFlashSaleOutOfStock
	}

	// 尝试原子扣减库存
	success, err := s.flashSaleCache.DecrementStock(ctx, cmd.FlashSaleID, cmd.ProductID)
	if err != nil {
		return nil, err
	}
	if !success {
		return &FlashSalePurchaseResult{
			Success: false,
			Message: "秒杀已售罄",
		}, nil
	}

	// 创建订单
	order := commerceEntity.NewOrder(cmd.UserID, cmd.ProductID, commerceEntity.PayType(cmd.PayType), stock.FlashPrice)
	if err := s.orderRepo.Create(ctx, order); err != nil {
		// 回滚库存
		s.flashSaleCache.RestoreStock(ctx, cmd.FlashSaleID, cmd.ProductID)
		log.Error.Printf("FlashSaleCommandService.Purchase: create order failed: %v", err)
		return nil, err
	}

	return &FlashSalePurchaseResult{
		Success: true,
		OrderID: order.ID,
		Message: "秒杀成功",
	}, nil
}

var ErrFlashSaleNotFound = &CommerceError{Code: "FLASH_SALE_NOT_FOUND", Message: "秒杀场次不存在"}
var ErrFlashSaleNotActive = &CommerceError{Code: "FLASH_SALE_NOT_ACTIVE", Message: "秒杀未开始"}
var ErrFlashSaleProductNotFound = &CommerceError{Code: "FLASH_SALE_PRODUCT_NOT_FOUND", Message: "秒杀商品不存在"}
var ErrFlashSaleOutOfStock = &CommerceError{Code: "FLASH_SALE_OUT_OF_STOCK", Message: "秒杀已售罄"}
