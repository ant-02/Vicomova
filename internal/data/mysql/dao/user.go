package dao

import (
	"context"
	"vicomova/internal/data/mysql"
	"vicomova/internal/domain/user"

	"gorm.io/gorm"
)

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (d *UserDAO) Create(ctx context.Context, u *user.User) error {
	return mysql.GetDB().WithContext(ctx).Create(u).Error
}

func (d *UserDAO) GetByID(ctx context.Context, id int64) (*user.User, error) {
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

func (d *UserDAO) GetByUsername(ctx context.Context, username string) (*user.User, error) {
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

func (d *UserDAO) Update(ctx context.Context, u *user.User) error {
	return mysql.GetDB().WithContext(ctx).Save(u).Error
}

func (d *UserDAO) Delete(ctx context.Context, id int64) error {
	return mysql.GetDB().WithContext(ctx).Delete(&user.User{}, id).Error
}
