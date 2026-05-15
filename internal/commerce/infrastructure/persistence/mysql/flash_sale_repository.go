package mysql

import (
	"context"
	"time"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type FlashSaleRepository struct {
	mysql *sharedMysql.Client
}

func NewFlashSaleRepository(mysqlClient *sharedMysql.Client) commerceRepo.FlashSaleRepository {
	return &FlashSaleRepository{mysql: mysqlClient}
}

func (r *FlashSaleRepository) Create(ctx context.Context, fs *commerceEntity.FlashSale) error {
	po := FlashSaleToPO(fs)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("FlashSaleRepository.Create: failed: %v", err)
		return err
	}
	fs.ID = po.ID
	return nil
}

func (r *FlashSaleRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.FlashSale, error) {
	var po FlashSalePO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("FlashSaleRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToFlashSale(&po), nil
}

func (r *FlashSaleRepository) GetActiveFlashSale(ctx context.Context) (*commerceEntity.FlashSale, error) {
	now := time.Now()
	var po FlashSalePO
	err := r.mysql.WithContext(ctx).Where("status = ? AND start_time <= ? AND end_time > ?",
		commerceEntity.FlashSaleStatusActive, now, now).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("FlashSaleRepository.GetActiveFlashSale: failed: %v", err)
		return nil, err
	}
	return POToFlashSale(&po), nil
}

func (r *FlashSaleRepository) ListActive(ctx context.Context) ([]*commerceEntity.FlashSale, error) {
	now := time.Now()
	var pos []FlashSalePO
	if err := r.mysql.WithContext(ctx).Where("status = ? AND end_time > ?",
		commerceEntity.FlashSaleStatusActive, now).Order("start_time ASC").Find(&pos).Error; err != nil {
		log.Error.Printf("FlashSaleRepository.ListActive: failed: %v", err)
		return nil, err
	}
	sales := make([]*commerceEntity.FlashSale, len(pos))
	for i := range pos {
		sales[i] = POToFlashSale(&pos[i])
	}
	return sales, nil
}
