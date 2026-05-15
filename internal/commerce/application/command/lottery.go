package command

import (
	"context"
	"math/rand"
	"time"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	redisCache "vicomova/internal/commerce/infrastructure/persistence/redis"
	"vicomova/pkg/log"
)

type LotteryCommandService struct {
	lotteryRepo  commerceRepo.LotteryRepository
	recordRepo   commerceRepo.LotteryRecordRepository
	lotteryCache *redisCache.LotteryCache
}

func NewLotteryCommandService(
	lotteryRepo commerceRepo.LotteryRepository,
	recordRepo commerceRepo.LotteryRecordRepository,
	lotteryCache *redisCache.LotteryCache,
) *LotteryCommandService {
	return &LotteryCommandService{
		lotteryRepo:  lotteryRepo,
		recordRepo:   recordRepo,
		lotteryCache: lotteryCache,
	}
}

func (s *LotteryCommandService) Participate(ctx context.Context, cmd *ParticipateLotteryCommand) (*ParticipateLotteryResult, error) {
	// 检查是否已参与
	participated, err := s.lotteryCache.HasUserEntry(ctx, cmd.LotteryID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if participated {
		return nil, ErrLotteryAlreadyParticipated
	}

	// 检查抽奖活动
	lottery, err := s.lotteryRepo.GetByID(ctx, cmd.LotteryID)
	if err != nil {
		return nil, err
	}
	if lottery == nil {
		return nil, ErrLotteryNotFound
	}
	if !lottery.IsActive() {
		return nil, ErrLotteryNotActive
	}

	// 标记用户已参与
	if err := s.lotteryCache.SetUserEntry(ctx, cmd.LotteryID, cmd.UserID); err != nil {
		return nil, ErrLotteryAlreadyParticipated
	}

	// 模拟抽奖概率（简化处理，默认30%中奖率）
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	win := r.Float32() < 0.3
	prize := ""
	if win {
		prize = "恭喜中奖" // 简化，实际应解析 lottery.Prizes
	}

	// 记录抽奖
	record := commerceEntity.NewLotteryRecord(cmd.UserID, cmd.LotteryID, prize, win, 1)
	if err := s.recordRepo.RecordParticipation(ctx, record); err != nil {
		log.Error.Printf("LotteryCommandService.Participate: record participation failed: %v", err)
	}

	return &ParticipateLotteryResult{
		Success: true,
		Won:     win,
		Prize:   prize,
		Message: "抽奖完成",
	}, nil
}

var ErrLotteryNotFound = &CommerceError{Code: "LOTTERY_NOT_FOUND", Message: "抽奖不存在"}
var ErrLotteryNotActive = &CommerceError{Code: "LOTTERY_NOT_ACTIVE", Message: "抽奖未开始"}
var ErrLotteryAlreadyParticipated = &CommerceError{Code: "LOTTERY_ALREADY_PARTICIPATED", Message: "已参与过抽奖"}
