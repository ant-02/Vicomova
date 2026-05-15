package mysql

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type FlashSaleStockRepository struct {
	mysql *sharedMysql.Client
}

func NewFlashSaleStockRepository(mysqlClient *sharedMysql.Client) commerceRepo.FlashSaleStockRepository {
	return &FlashSaleStockRepository{mysql: mysqlClient}
}

func (r *FlashSaleStockRepository) Create(ctx context.Context, stock *commerceEntity.FlashSaleStock) error {
	po := FlashSaleStockToPO(stock)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("FlashSaleStockRepository.Create: failed: %v", err)
		return err
	}
	stock.ID = po.ID
	return nil
}

func (r *FlashSaleStockRepository) GetByFlashSaleID(ctx context.Context, flashSaleID int64) ([]*commerceEntity.FlashSaleStock, error) {
	var pos []FlashSaleStockPO
	if err := r.mysql.WithContext(ctx).Where("flash_sale_id = ?", flashSaleID).Find(&pos).Error; err != nil {
		log.Error.Printf("FlashSaleStockRepository.GetByFlashSaleID: failed: %v", err)
		return nil, err
	}
	stocks := make([]*commerceEntity.FlashSaleStock, len(pos))
	for i := range pos {
		stocks[i] = POToFlashSaleStock(&pos[i])
	}
	return stocks, nil
}

func (r *FlashSaleStockRepository) GetStock(ctx context.Context, flashSaleID, productID int64) (*commerceEntity.FlashSaleStock, error) {
	var po FlashSaleStockPO
	err := r.mysql.WithContext(ctx).Where("flash_sale_id = ? AND product_id = ?", flashSaleID, productID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("FlashSaleStockRepository.GetStock: failed: %v", err)
		return nil, err
	}
	return POToFlashSaleStock(&po), nil
}

func (r *FlashSaleStockRepository) DecrementStock(ctx context.Context, flashSaleID, productID int64) (bool, error) {
	result := r.mysql.WithContext(ctx).Model(&FlashSaleStockPO{}).
		Where("flash_sale_id = ? AND product_id = ? AND remain_stock > 0", flashSaleID, productID).
		Updates(map[string]interface{}{
			"remain_stock": gorm.Expr("remain_stock - 1"),
		})
	if result.Error != nil {
		log.Error.Printf("FlashSaleStockRepository.DecrementStock: failed: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
