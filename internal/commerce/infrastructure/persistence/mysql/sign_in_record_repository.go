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

type SignInRecordRepository struct {
	mysql *sharedMysql.Client
}

func NewSignInRecordRepository(mysqlClient *sharedMysql.Client) commerceRepo.SignInRecordRepository {
	return &SignInRecordRepository{mysql: mysqlClient}
}

func (r *SignInRecordRepository) Create(ctx context.Context, record *commerceEntity.SignInRecord) error {
	po := SignInRecordToPO(record)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("SignInRecordRepository.Create: failed: %v", err)
		return err
	}
	record.ID = po.ID
	return nil
}

func (r *SignInRecordRepository) GetTodaySignIn(ctx context.Context, userID int64) (*commerceEntity.SignInRecord, error) {
	today := time.Now().Format("2006-01-02")
	var po SignInRecordPO
	err := r.mysql.WithContext(ctx).Where("user_id = ? AND sign_date = ?", userID, today).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("SignInRecordRepository.GetTodaySignIn: failed: %v", err)
		return nil, err
	}
	return POToSignInRecord(&po), nil
}

func (r *SignInRecordRepository) GetLastSignIn(ctx context.Context, userID int64) (*commerceEntity.SignInRecord, error) {
	var po SignInRecordPO
	err := r.mysql.WithContext(ctx).Where("user_id = ?", userID).Order("sign_date DESC").First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("SignInRecordRepository.GetLastSignIn: failed: %v", err)
		return nil, err
	}
	return POToSignInRecord(&po), nil
}

func (r *SignInRecordRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*commerceEntity.SignInRecord, bool, error) {
	var pos []SignInRecordPO
	query := r.mysql.WithContext(ctx).Where("user_id = ?", userID)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}
	if err := query.Order("id DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		log.Error.Printf("SignInRecordRepository.ListByUser: failed: %v", err)
		return nil, false, err
	}
	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}
	records := make([]*commerceEntity.SignInRecord, len(pos))
	for i := range pos {
		records[i] = POToSignInRecord(&pos[i])
	}
	return records, hasMore, nil
}
