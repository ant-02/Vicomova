package mysql

import (
	"context"

	userEntity "vicomova/internal/user/domain/entity"
	userRepo "vicomova/internal/user/domain/repository"
	userVO "vicomova/internal/user/domain/valueobject"
	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	"vicomova/internal/shared/pkg/log"

	"gorm.io/gorm"
)

type UserRepository struct {
	mysql *sharedMysql.Client
}

func NewUserRepository(mysqlClient *sharedMysql.Client) userRepo.UserRepository {
	return &UserRepository{mysql: mysqlClient}
}

func (r *UserRepository) Create(ctx context.Context, u *userEntity.User) error {
	po := UserToPO(u)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("UserRepository.Create: failed for username=%s: %v", u.Username.String(), err)
		return err
	}
	log.Info.Printf("UserRepository.Create: created user username=%s", u.Username.String())
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*userEntity.User, error) {
	var po UserPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("UserRepository.GetByID: failed for id=%d: %v", id, err)
		return nil, err
	}
	return POToUser(&po), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username *userVO.Username) (*userEntity.User, error) {
	var po UserPO
	err := r.mysql.WithContext(ctx).Where("username = ?", username.String()).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("UserRepository.GetByUsername: failed for username=%s: %v", username.String(), err)
		return nil, err
	}
	return POToUser(&po), nil
}

func (r *UserRepository) Update(ctx context.Context, u *userEntity.User) error {
	po := UserToPO(u)
	if err := r.mysql.WithContext(ctx).Save(po).Error; err != nil {
		log.Error.Printf("UserRepository.Update: failed for userID=%d: %v", u.ID, err)
		return err
	}
	log.Info.Printf("UserRepository.Update: updated userID=%d", u.ID)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if err := r.mysql.WithContext(ctx).Delete(&UserPO{}, id).Error; err != nil {
		log.Error.Printf("UserRepository.Delete: failed for id=%d: %v", id, err)
		return err
	}
	log.Info.Printf("UserRepository.Delete: deleted userID=%d", id)
	return nil
}