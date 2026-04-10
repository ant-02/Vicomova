package mysql

import (
	"context"
	"vicomova/internal/data/mysql"
	"vicomova/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	return mysql.GetDB().WithContext(ctx).Create(u).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := mysql.GetDB().WithContext(ctx).First(&u, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := mysql.GetDB().WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return mysql.GetDB().WithContext(ctx).Save(u).Error
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return mysql.GetDB().WithContext(ctx).Delete(&user.User{}, id).Error
}
