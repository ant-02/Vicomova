package redis

import (
	"context"
	"encoding/json"
	"time"

	userRepo "vicomova/internal/user/domain/repository"
	userVO "vicomova/internal/user/domain/valueobject"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	"vicomova/internal/shared/pkg/log"

	"github.com/redis/go-redis/v9"
)

const refreshTokenPrefix = "refresh_token:"

type RefreshTokenRepository struct {
	redis *sharedRedis.Client
}

func NewRefreshTokenRepository(redisClient *sharedRedis.Client) userRepo.RefreshTokenRepository {
	return &RefreshTokenRepository{redis: redisClient}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *userVO.RefreshToken) error {
	data := &RefreshTokenData{
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
	}
	ttl := rt.ExpiresAt.Sub(rt.CreatedAt)
	if err := SetRefreshToken(ctx, r.redis, rt.Token, data, ttl); err != nil {
		log.Error.Printf("RefreshTokenRepository.Create: failed for userID=%d: %v", rt.UserID, err)
		return err
	}
	log.Info.Printf("RefreshTokenRepository.Create: created token for userID=%d", rt.UserID)
	return nil
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*userVO.RefreshToken, error) {
	data, err := GetRefreshToken(ctx, r.redis, token)
	if err != nil {
		log.Error.Printf("RefreshTokenRepository.GetByToken: failed: %v", err)
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &userVO.RefreshToken{
		UserID:    data.UserID,
		Token:     token,
		ExpiresAt: data.ExpiresAt,
		Revoked:   false,
	}, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	if err := DeleteRefreshToken(ctx, r.redis, token); err != nil {
		log.Error.Printf("RefreshTokenRepository.Revoke: failed: %v", err)
		return err
	}
	log.Info.Printf("RefreshTokenRepository.Revoke: revoked token")
	return nil
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	if err := DeleteUserRefreshTokens(ctx, r.redis, userID); err != nil {
		log.Error.Printf("RefreshTokenRepository.RevokeAllForUser: failed for userID=%d: %v", userID, err)
		return err
	}
	log.Info.Printf("RefreshTokenRepository.RevokeAllForUser: revoked all tokens for userID=%d", userID)
	return nil
}

type RefreshTokenData struct {
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func SetRefreshToken(ctx context.Context, client *sharedRedis.Client, token string, data *RefreshTokenData, ttl time.Duration) error {
	key := refreshTokenPrefix + token
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, jsonData, ttl).Err()
}

func GetRefreshToken(ctx context.Context, client *sharedRedis.Client, token string) (*RefreshTokenData, error) {
	key := refreshTokenPrefix + token
	val, err := client.Get(ctx, key).Result()
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

func DeleteRefreshToken(ctx context.Context, client *sharedRedis.Client, token string) error {
	key := refreshTokenPrefix + token
	return client.Del(ctx, key).Err()
}

func DeleteUserRefreshTokens(ctx context.Context, client *sharedRedis.Client, userID int64) error {
	pattern := refreshTokenPrefix + "*"
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var data RefreshTokenData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}
		if data.UserID == userID {
			client.Del(ctx, key)
		}
	}
	return iter.Err()
}