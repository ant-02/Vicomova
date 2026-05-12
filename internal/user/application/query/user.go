package query

import (
	"context"

	userEntity "vicomova/internal/user/domain/entity"
	userVO "vicomova/internal/user/domain/valueobject"
	"vicomova/pkg/errors"
	"vicomova/pkg/log"
)

func (s *UserQueryService) GetUser(ctx context.Context, query *GetUserQuery) (*UserResult, error) {
	var u *userEntity.User
	var err error

	// Query by UserID if provided
	if query.UserID != 0 {
		u, err = s.userRepo.GetByID(ctx, query.UserID)
		if err != nil {
			log.Error.Printf("GetUser: GetByID failed for userID=%d: %v", query.UserID, err)
			return nil, errors.ErrInternalServer
		}
	} else if query.Username != "" {
		// Fall back to Username if UserID is not provided
		username, err := userVO.NewUsername(query.Username)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}
		u, err = s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			log.Error.Printf("GetUser: GetByUsername failed for username=%s: %v", query.Username, err)
			return nil, errors.ErrInternalServer
		}
	} else {
		return nil, errors.ErrUserNotFound
	}

	if u == nil {
		log.Warn.Printf("GetUser: user not found, userID=%d, username=%s", query.UserID, query.Username)
		return nil, errors.ErrUserNotFound
	}

	return &UserResult{
		UserID:   u.ID,
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}, nil
}

func (s *UserQueryService) BatchGetUsers(ctx context.Context, query *BatchGetUsersQuery) ([]*UserResult, error) {
	if len(query.UserIDs) == 0 {
		return nil, nil
	}
	users, err := s.userRepo.GetByIDs(ctx, query.UserIDs)
	if err != nil {
		log.Error.Printf("BatchGetUsers: GetByIDs failed: %v", err)
		return nil, errors.ErrInternalServer
	}
	results := make([]*UserResult, 0, len(users))
	for _, u := range users {
		results = append(results, &UserResult{
			UserID:   u.ID,
			Username: u.Username.String(),
			Email:    u.Email.String(),
		})
	}
	return results, nil
}
