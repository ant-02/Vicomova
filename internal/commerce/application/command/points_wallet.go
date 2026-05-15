package command

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	"vicomova/pkg/log"
)

type PointsWalletCommandService struct {
	walletRepo commerceRepo.PointsWalletRepository
	txnRepo    commerceRepo.TransactionRepository
}

func NewPointsWalletCommandService(
	walletRepo commerceRepo.PointsWalletRepository,
	txnRepo commerceRepo.TransactionRepository,
) *PointsWalletCommandService {
	return &PointsWalletCommandService{
		walletRepo: walletRepo,
		txnRepo:    txnRepo,
	}
}

func (s *PointsWalletCommandService) RechargePoints(ctx context.Context, cmd *RechargePointsCommand) (*RechargePointsResult, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		log.Error.Printf("PointsWalletCommandService.RechargePoints: get wallet failed: %v", err)
		return nil, err
	}
	if wallet == nil {
		wallet = commerceEntity.NewPointsWallet(cmd.UserID)
		if err := s.walletRepo.Create(ctx, wallet); err != nil {
			log.Error.Printf("PointsWalletCommandService.RechargePoints: create wallet failed: %v", err)
			return nil, err
		}
	}

	oldVersion := wallet.Version
	newVersion := oldVersion + 1
	success, err := s.walletRepo.UpdateVersion(ctx, cmd.UserID, oldVersion, newVersion, cmd.Amount)
	if err != nil {
		log.Error.Printf("PointsWalletCommandService.RechargePoints: update version failed: %v", err)
		return nil, err
	}
	if !success {
		log.Warn.Printf("PointsWalletCommandService.RechargePoints: optimistic lock failed for user=%d", cmd.UserID)
		return nil, err
	}

	// 刷新获取最新钱包
	wallet, err = s.walletRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 记录交易
	txn := commerceEntity.NewTransaction(cmd.UserID, "recharge", cmd.Amount, wallet.Points, "积分充值", "completed", cmd.PaymentID)
	if err := s.txnRepo.Create(ctx, txn); err != nil {
		log.Warn.Printf("PointsWalletCommandService.RechargePoints: create transaction failed: %v", err)
	}

	return &RechargePointsResult{
		Wallet: &PointsWalletDTO{
			UserID:  wallet.UserID,
			Points:  wallet.Points,
			Version: wallet.Version,
		},
		NewPoints: wallet.Points,
	}, nil
}

func (s *PointsWalletCommandService) DeductPoints(ctx context.Context, cmd *DeductPointsCommand) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if wallet == nil || wallet.Points < cmd.Amount {
		return ErrInsufficientPoints
	}

	oldVersion := wallet.Version
	newVersion := oldVersion + 1
	success, err := s.walletRepo.UpdateVersion(ctx, cmd.UserID, oldVersion, newVersion, -cmd.Amount)
	if err != nil {
		return err
	}
	if !success {
		return ErrInsufficientPoints
	}

	// 记录交易
	wallet, _ = s.walletRepo.GetByUserID(ctx, cmd.UserID)
	txn := commerceEntity.NewTransaction(cmd.UserID, "spend", cmd.Amount, wallet.Points, cmd.Memo, "completed", "")
	if err := s.txnRepo.Create(ctx, txn); err != nil {
		log.Warn.Printf("PointsWalletCommandService.DeductPoints: create transaction failed: %v", err)
	}

	return nil
}

var ErrInsufficientPoints = &CommerceError{Code: "INSUFFICIENT_POINTS", Message: "积分不足"}

type CommerceError struct {
	Code    string
	Message string
}

func (e *CommerceError) Error() string {
	return e.Message
}

var _ error = (*CommerceError)(nil)
