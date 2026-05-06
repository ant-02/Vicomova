package query

import (
	"context"

	"vicomova/internal/shared/pkg/errors"
	"vicomova/internal/shared/pkg/log"
	userVO "vicomova/internal/user/domain/valueobject"
)

func (s *UserQueryService) GetUser(ctx context.Context, query *GetUserQuery) (*UserResult, error) {
	username, err := userVO.NewUsername(query.Username)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Error.Printf("GetUser: GetByUsername failed for userID=%d: %v", query.UserID, err)
		return nil, errors.ErrInternalServer
	}
	if u == nil {
		log.Warn.Printf("GetUser: user not found, username=%s", query.Username)
		return nil, errors.ErrUserNotFound
	}

	return &UserResult{
		UserID:   u.ID,
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}, nil
}
