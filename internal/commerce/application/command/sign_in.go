package command

import (
	"context"
	"time"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	commerceSvc "vicomova/internal/commerce/domain/service"
	"vicomova/pkg/log"
)

type SignInCommandService struct {
	signInRepo commerceRepo.SignInRecordRepository
	walletRepo commerceRepo.PointsWalletRepository
	txnRepo    commerceRepo.TransactionRepository
	bonusCalc  *commerceSvc.SignInBonusCalculator
}

func NewSignInCommandService(
	signInRepo commerceRepo.SignInRecordRepository,
	walletRepo commerceRepo.PointsWalletRepository,
	txnRepo commerceRepo.TransactionRepository,
	bonusCalc *commerceSvc.SignInBonusCalculator,
) *SignInCommandService {
	return &SignInCommandService{
		signInRepo: signInRepo,
		walletRepo: walletRepo,
		txnRepo:    txnRepo,
		bonusCalc:  bonusCalc,
	}
}

// ProcessSignIn 处理签到事件（由 Kafka 消费者调用）
func (s *SignInCommandService) ProcessSignIn(ctx context.Context, cmd *ProcessSignInCommand) (*ProcessSignInResult, error) {
	// 解析签到日期
	signDate, err := time.Parse("2006-01-02", cmd.SignDate)
	if err != nil {
		log.Error.Printf("SignInCommandService.ProcessSignIn: parse date failed: %v", err)
		return nil, err
	}

	// 检查今日是否已签到
	todayRecord, err := s.signInRepo.GetTodaySignIn(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if todayRecord != nil {
		return &ProcessSignInResult{
			Success: false,
			Message: "今日已签到",
		}, nil
	}

	// 计算积分
	bonus := s.bonusCalc.Calculate(cmd.ConsecutiveDays)

	// 记录签到
	record := commerceEntity.NewSignInRecord(cmd.UserID, signDate, cmd.ConsecutiveDays, bonus)
	if err := s.signInRepo.Create(ctx, record); err != nil {
		log.Error.Printf("SignInCommandService.ProcessSignIn: create record failed: %v", err)
		return nil, err
	}

	// 确保钱包存在
	wallet, err := s.walletRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		wallet = commerceEntity.NewPointsWallet(cmd.UserID)
		if err := s.walletRepo.Create(ctx, wallet); err != nil {
			return nil, err
		}
	}

	// 发放积分
	oldVersion := wallet.Version
	newVersion := oldVersion + 1
	success, err := s.walletRepo.UpdateVersion(ctx, cmd.UserID, oldVersion, newVersion, bonus)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, err
	}

	// 刷新钱包获取最新余额
	wallet, _ = s.walletRepo.GetByUserID(ctx, cmd.UserID)

	// 记录积分变动
	txn := commerceEntity.NewTransaction(cmd.UserID, "earn", bonus, wallet.Points, "签到奖励", "completed", "")
	if err := s.txnRepo.Create(ctx, txn); err != nil {
		log.Warn.Printf("SignInCommandService.ProcessSignIn: create transaction failed: %v", err)
	}

	return &ProcessSignInResult{
		Success:         true,
		PointsEarned:    bonus,
		ConsecutiveDays: cmd.ConsecutiveDays,
		Message:         "签到成功",
	}, nil
}
