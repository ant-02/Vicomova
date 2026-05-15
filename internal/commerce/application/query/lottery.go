package query

import (
	"context"

	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type LotteryQueryService struct {
	lotteryRepo commerceRepo.LotteryRepository
}

func NewLotteryQueryService(lotteryRepo commerceRepo.LotteryRepository) *LotteryQueryService {
	return &LotteryQueryService{lotteryRepo: lotteryRepo}
}

func (s *LotteryQueryService) ListLotteries(ctx context.Context, query *ListLotteriesQuery) (*ListLotteriesResult, error) {
	lotteries, err := s.lotteryRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*LotteryDTO, len(lotteries))
	for i := range lotteries {
		dtos[i] = &LotteryDTO{
			ID:            lotteries[i].ID,
			Name:          lotteries[i].Name,
			Description:   lotteries[i].Description,
			EntryPoints:   lotteries[i].EntryPoints,
			StartTime:     lotteries[i].StartTime.Unix(),
			EndTime:       lotteries[i].EndTime.Unix(),
			RemainTickets: lotteries[i].RemainTickets,
			Status:        lotteries[i].Status,
		}
	}

	return &ListLotteriesResult{Lotteries: dtos}, nil
}
