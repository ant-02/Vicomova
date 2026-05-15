package query

import (
	"context"

	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type FlashSaleQueryService struct {
	flashSaleRepo commerceRepo.FlashSaleRepository
}

func NewFlashSaleQueryService(flashSaleRepo commerceRepo.FlashSaleRepository) *FlashSaleQueryService {
	return &FlashSaleQueryService{flashSaleRepo: flashSaleRepo}
}

func (s *FlashSaleQueryService) ListFlashSales(ctx context.Context, query *ListFlashSalesQuery) (*ListFlashSalesResult, error) {
	flashSales, err := s.flashSaleRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*FlashSaleDTO, len(flashSales))
	for i := range flashSales {
		dtos[i] = &FlashSaleDTO{
			ID:        flashSales[i].ID,
			Name:      flashSales[i].Name,
			StartTime: flashSales[i].StartTime.Unix(),
			EndTime:   flashSales[i].EndTime.Unix(),
			Status:    flashSales[i].Status,
		}
	}

	return &ListFlashSalesResult{FlashSales: dtos}, nil
}
