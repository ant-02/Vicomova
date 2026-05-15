package mysql

import (
	"context"
	"strconv"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type OrderRepository struct {
	mysql *sharedMysql.Client
}

func NewOrderRepository(mysqlClient *sharedMysql.Client) commerceRepo.OrderRepository {
	return &OrderRepository{mysql: mysqlClient}
}

func (r *OrderRepository) Create(ctx context.Context, order *commerceEntity.Order) error {
	po := OrderToPO(order)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("OrderRepository.Create: failed: %v", err)
		return err
	}
	order.ID = po.ID
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Order, error) {
	var po OrderPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("OrderRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToOrder(&po), nil
}

func (r *OrderRepository) GetByUserAndProduct(ctx context.Context, userID, productID int64) (*commerceEntity.Order, error) {
	var po OrderPO
	err := r.mysql.WithContext(ctx).Where("user_id = ? AND product_id = ? AND status NOT IN (?, ?)",
		userID, productID, int8(commerceEntity.OrderStatusCancelled), int8(commerceEntity.OrderStatusRefunded)).
		Order("id DESC").First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("OrderRepository.GetByUserAndProduct: failed: %v", err)
		return nil, err
	}
	return POToOrder(&po), nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status commerceEntity.OrderStatus) error {
	if err := r.mysql.WithContext(ctx).Model(&OrderPO{}).Where("id = ?", id).
		Update("status", int8(status)).Error; err != nil {
		log.Error.Printf("OrderRepository.UpdateStatus: failed: %v", err)
		return err
	}
	return nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64, cursor string, limit int) ([]*commerceEntity.Order, bool, error) {
	var pos []OrderPO
	query := r.mysql.WithContext(ctx).Where("user_id = ?", userID)
	if cursor != "" {
		if id, err := strconv.ParseInt(cursor, 10, 64); err == nil {
			query = query.Where("id < ?", id)
		}
	}
	if err := query.Order("id DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		log.Error.Printf("OrderRepository.ListByUser: failed: %v", err)
		return nil, false, err
	}
	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}
	orders := make([]*commerceEntity.Order, len(pos))
	for i := range pos {
		orders[i] = POToOrder(&pos[i])
	}
	return orders, hasMore, nil
}
