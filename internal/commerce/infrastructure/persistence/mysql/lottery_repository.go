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

type LotteryRepository struct {
	mysql *sharedMysql.Client
}

func NewLotteryRepository(mysqlClient *sharedMysql.Client) commerceRepo.LotteryRepository {
	return &LotteryRepository{mysql: mysqlClient}
}

func (r *LotteryRepository) Create(ctx context.Context, l *commerceEntity.LotteryDraw) error {
	po := LotteryDrawToPO(l)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("LotteryRepository.Create: failed: %v", err)
		return err
	}
	l.ID = po.ID
	return nil
}

func (r *LotteryRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.LotteryDraw, error) {
	var po LotteryDrawPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("LotteryRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToLotteryDraw(&po), nil
}

func (r *LotteryRepository) GetActiveLottery(ctx context.Context) (*commerceEntity.LotteryDraw, error) {
	now := time.Now()
	var po LotteryDrawPO
	err := r.mysql.WithContext(ctx).Where("status = ? AND start_time <= ? AND end_time > ? AND remain_tickets > 0",
		commerceEntity.LotteryStatusActive, now, now).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("LotteryRepository.GetActiveLottery: failed: %v", err)
		return nil, err
	}
	return POToLotteryDraw(&po), nil
}

func (r *LotteryRepository) ListActive(ctx context.Context) ([]*commerceEntity.LotteryDraw, error) {
	now := time.Now()
	var pos []LotteryDrawPO
	if err := r.mysql.WithContext(ctx).Where("status = ? AND end_time > ?",
		commerceEntity.LotteryStatusActive, now).Order("start_time ASC").Find(&pos).Error; err != nil {
		log.Error.Printf("LotteryRepository.ListActive: failed: %v", err)
		return nil, err
	}
	lotteries := make([]*commerceEntity.LotteryDraw, len(pos))
	for i := range pos {
		lotteries[i] = POToLotteryDraw(&pos[i])
	}
	return lotteries, nil
}

type LotteryRecordRepository struct {
	mysql *sharedMysql.Client
}

func NewLotteryRecordRepository(mysqlClient *sharedMysql.Client) commerceRepo.LotteryRecordRepository {
	return &LotteryRecordRepository{mysql: mysqlClient}
}

func (r *LotteryRecordRepository) RecordParticipation(ctx context.Context, record *commerceEntity.LotteryRecord) error {
	po := LotteryRecordToPO(record)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("LotteryRecordRepository.RecordParticipation: failed: %v", err)
		return err
	}
	record.ID = po.ID
	return nil
}

func (r *LotteryRecordRepository) GetUserParticipation(ctx context.Context, userID, lotteryID int64) (*commerceEntity.LotteryRecord, error) {
	var po LotteryRecordPO
	err := r.mysql.WithContext(ctx).Where("user_id = ? AND lottery_id = ?", userID, lotteryID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("LotteryRecordRepository.GetUserParticipation: failed: %v", err)
		return nil, err
	}
	return POToLotteryRecord(&po), nil
}

func (r *LotteryRecordRepository) CountByUser(ctx context.Context, userID, lotteryID int64) (int, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&LotteryRecordPO{}).
		Where("user_id = ? AND lottery_id = ?", userID, lotteryID).Count(&count).Error; err != nil {
		log.Error.Printf("LotteryRecordRepository.CountByUser: failed: %v", err)
		return 0, err
	}
	return int(count), nil
}
