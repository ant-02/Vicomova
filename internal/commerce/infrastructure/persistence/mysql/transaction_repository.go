package mysql

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	mysql *sharedMysql.Client
}

func NewTransactionRepository(mysqlClient *sharedMysql.Client) commerceRepo.TransactionRepository {
	return &TransactionRepository{mysql: mysqlClient}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *commerceEntity.Transaction) error {
	po := TransactionToPO(tx)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("TransactionRepository.Create: failed: %v", err)
		return err
	}
	tx.ID = po.ID
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Transaction, error) {
	var po TransactionPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("TransactionRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToTransaction(&po), nil
}

func (r *TransactionRepository) ListByUser(ctx context.Context, userID int64, cursor int64, limit int) ([]*commerceEntity.Transaction, bool, error) {
	var pos []TransactionPO
	query := r.mysql.WithContext(ctx).Where("user_id = ?", userID)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}
	if err := query.Order("id DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		log.Error.Printf("TransactionRepository.ListByUser: failed: %v", err)
		return nil, false, err
	}
	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}
	txs := make([]*commerceEntity.Transaction, len(pos))
	for i := range pos {
		txs[i] = POToTransaction(&pos[i])
	}
	return txs, hasMore, nil
}
