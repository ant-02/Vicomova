package query

import (
	"context"
	"strconv"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type OrderQueryService struct {
	orderRepo commerceRepo.OrderRepository
}

func NewOrderQueryService(orderRepo commerceRepo.OrderRepository) *OrderQueryService {
	return &OrderQueryService{orderRepo: orderRepo}
}

func (s *OrderQueryService) ListOrders(ctx context.Context, query *ListOrdersQuery) (*ListOrdersResult, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	orders, hasMore, err := s.orderRepo.ListByUser(ctx, query.UserID, query.Cursor, query.Limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]*OrderDTO, len(orders))
	var nextCursor string
	for i := range orders {
		dtos[i] = &OrderDTO{
			ID:        orders[i].ID,
			UserID:    orders[i].UserID,
			ProductID: orders[i].ProductID,
			PayType:   int(orders[i].PayType),
			PricePaid: orders[i].PricePaid,
			Status:    int(orders[i].Status),
		}
		if i == len(orders)-1 && len(orders) > 0 {
			nextCursor = strconv.FormatInt(orders[i].ID, 10)
		}
	}

	return &ListOrdersResult{
		Orders:     dtos,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

var _ = commerceEntity.OrderStatusPending // suppress unused import warning
