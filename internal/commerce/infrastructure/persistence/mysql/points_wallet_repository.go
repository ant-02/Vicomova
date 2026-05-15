package mysql

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type PointsWalletRepository struct {
	mysql *sharedMysql.Client
}

func NewPointsWalletRepository(mysqlClient *sharedMysql.Client) commerceRepo.PointsWalletRepository {
	return &PointsWalletRepository{mysql: mysqlClient}
}

func (r *PointsWalletRepository) Create(ctx context.Context, wallet *commerceEntity.PointsWallet) error {
	po := PointsWalletToPO(wallet)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("PointsWalletRepository.Create: failed: %v", err)
		return err
	}
	wallet.ID = po.ID
	return nil
}

func (r *PointsWalletRepository) GetByUserID(ctx context.Context, userID int64) (*commerceEntity.PointsWallet, error) {
	var po PointsWalletPO
	err := r.mysql.WithContext(ctx).Where("user_id = ?", userID).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("PointsWalletRepository.GetByUserID: failed: %v", err)
		return nil, err
	}
	return POToPointsWallet(&po), nil
}

func (r *PointsWalletRepository) Update(ctx context.Context, wallet *commerceEntity.PointsWallet) error {
	po := PointsWalletToPO(wallet)
	if err := r.mysql.WithContext(ctx).Save(po).Error; err != nil {
		log.Error.Printf("PointsWalletRepository.Update: failed: %v", err)
		return err
	}
	return nil
}

func (r *PointsWalletRepository) UpdateVersion(ctx context.Context, userID int64, oldVersion int64, newVersion int64, delta int64) (bool, error) {
	result := r.mysql.WithContext(ctx).Model(&PointsWalletPO{}).
		Where("user_id = ? AND version = ?", userID, oldVersion).
		Updates(map[string]interface{}{
			"points":  gorm.Expr("points + ?", delta),
			"version": newVersion,
		})
	if result.Error != nil {
		log.Error.Printf("PointsWalletRepository.UpdateVersion: failed: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
