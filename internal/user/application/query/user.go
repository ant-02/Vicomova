package query

import (
	"context"

	"vicomova/internal/user/application/command"
	errorsPkg "vicomova/internal/shared/pkg/errors"
	"vicomova/internal/shared/pkg/log"
)

func (s *UserQueryService) GetUser(ctx context.Context, query *GetUserQuery) (*command.UserResult, error) {
	u, err := s.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		log.Error.Printf("GetUser: GetByID failed for userID=%d: %v", query.UserID, err)
		return nil, errorsPkg.ErrInternalServer
	}
	if u == nil {
		log.Warn.Printf("GetUser: user not found, userID=%d", query.UserID)
		return nil, errorsPkg.ErrUserNotFound
	}

	return &command.UserResult{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
