package user

import (
	"context"

	"vicomova/internal/domain/user"
	errorsPkg "vicomova/internal/pkg/errors"
)

type GetUserQuery struct {
	UserID int64
}

type UserQueryService struct {
	userRepo user.UserRepository
}

func NewUserQueryService(userRepo user.UserRepository) *UserQueryService {
	return &UserQueryService{userRepo: userRepo}
}

func (s *UserQueryService) GetUser(ctx context.Context, query *GetUserQuery) (*UserResult, error) {
	u, err := s.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if u == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	return &UserResult{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
