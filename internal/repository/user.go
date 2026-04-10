package repository

import (
	"context"
	"vicomova/internal/data/mysql/dao"
	"vicomova/internal/domain/user"
)

type UserRepositoryImpl struct {
	dao *dao.UserDAO
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		dao: dao.NewUserDAO(),
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, u *user.User) error {
	return r.dao.Create(ctx, u)
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.User, error) {
	return r.dao.GetByID(ctx, id)
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.dao.GetByUsername(ctx, username)
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *user.User) error {
	return r.dao.Update(ctx, u)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.dao.Delete(ctx, id)
}
