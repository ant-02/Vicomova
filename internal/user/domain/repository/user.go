package repository

import (
	"context"
	"time"

	"vicomova/internal/user/domain/entity"
	"vicomova/internal/user/domain/valueobject"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUsername(ctx context.Context, username *valueobject.Username) (*entity.User, error)
	Update(ctx context.Context, u *entity.User) error
	Delete(ctx context.Context, id int64) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, token string, ttl time.Duration) error
	GetUserID(ctx context.Context, token string) (int64, error)
	Exists(ctx context.Context, token string) (bool, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID int64) error
}

type EmailCodeRepository interface {
	Store(ctx context.Context, email, code string, ttl time.Duration) error
	Verify(ctx context.Context, email, code string) (bool, error)
	Delete(ctx context.Context, email string) error
}