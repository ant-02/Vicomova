package redis

import (
	"context"
	"encoding/json"
	"time"

	redisClient "vicomova/internal/data/redis"
	"vicomova/internal/domain/user"

	"github.com/redis/go-redis/v9"
)

const refreshTokenPrefix = "refresh_token:"

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *user.RefreshToken) error {
	data := &RefreshTokenData{
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
	}
	ttl := rt.ExpiresAt.Sub(rt.CreatedAt)
	return SetRefreshToken(ctx, rt.Token, data, ttl)
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*user.RefreshToken, error) {
	data, err := GetRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &user.RefreshToken{
		UserID:    data.UserID,
		Token:     token,
		ExpiresAt: data.ExpiresAt,
		Revoked:   false,
	}, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	return DeleteRefreshToken(ctx, token)
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	return DeleteUserRefreshTokens(ctx, userID)
}

type RefreshTokenData struct {
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func SetRefreshToken(ctx context.Context, token string, data *RefreshTokenData, ttl time.Duration) error {
	key := refreshTokenPrefix + token
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisClient.GetClient().Set(ctx, key, jsonData, ttl).Err()
}

func GetRefreshToken(ctx context.Context, token string) (*RefreshTokenData, error) {
	key := refreshTokenPrefix + token
	val, err := redisClient.GetClient().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data RefreshTokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func DeleteRefreshToken(ctx context.Context, token string) error {
	key := refreshTokenPrefix + token
	return redisClient.GetClient().Del(ctx, key).Err()
}

func DeleteUserRefreshTokens(ctx context.Context, userID int64) error {
	pattern := refreshTokenPrefix + "*"
	iter := redisClient.GetClient().Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := redisClient.GetClient().Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var data RefreshTokenData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}
		if data.UserID == userID {
			redisClient.GetClient().Del(ctx, key)
		}
	}
	return iter.Err()
}
