package query

import (
	"context"

	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type PointsWalletQueryService struct {
	walletRepo commerceRepo.PointsWalletRepository
}

func NewPointsWalletQueryService(walletRepo commerceRepo.PointsWalletRepository) *PointsWalletQueryService {
	return &PointsWalletQueryService{walletRepo: walletRepo}
}

func (s *PointsWalletQueryService) GetWallet(ctx context.Context, query *GetWalletQuery) (*GetWalletResult, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return &GetWalletResult{
			Wallet: &WalletDTO{
				UserID: query.UserID,
				Points: 0,
			},
		}, nil
	}
	return &GetWalletResult{
		Wallet: &WalletDTO{
			UserID:  wallet.UserID,
			Points:  wallet.Points,
			Version: wallet.Version,
		},
	}, nil
}
