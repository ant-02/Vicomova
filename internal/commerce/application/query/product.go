package query

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type ProductQueryService struct {
	productRepo commerceRepo.ProductRepository
}

func NewProductQueryService(productRepo commerceRepo.ProductRepository) *ProductQueryService {
	return &ProductQueryService{productRepo: productRepo}
}

func (s *ProductQueryService) ListProducts(ctx context.Context, query *ListProductsQuery) (*ListProductsResult, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	products, nextCursor, hasMore, err := s.productRepo.List(ctx, commerceEntity.ProductType(query.ProductType), query.Cursor, query.Limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]*ProductDTO, len(products))
	for i := range products {
		dtos[i] = &ProductDTO{
			ID:              products[i].ID,
			Name:            products[i].Name,
			Description:     products[i].Description,
			ProductType:     int(products[i].ProductType),
			Price:           products[i].Price,
			WalletPrice:     products[i].WalletPrice,
			Stock:           products[i].Stock,
			VideoID:         products[i].VideoID,
			MembershipLevel: products[i].MembershipLevel,
			IsActive:        products[i].IsActive,
		}
	}

	return &ListProductsResult{
		Products:   dtos,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
