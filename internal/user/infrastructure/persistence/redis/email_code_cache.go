package redis

import (
	"context"
	"errors"
	"time"

	userRepo "vicomova/internal/user/domain/repository"
	sharedRedis "vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

var ErrCodeNotFound = errors.New("verification code not found")

type EmailCodeRepository struct {
	redis *sharedRedis.Client
}

func NewEmailCodeRepository(redisClient *sharedRedis.Client) userRepo.EmailCodeRepository {
	return &EmailCodeRepository{redis: redisClient}
}

func (r *EmailCodeRepository) Store(ctx context.Context, email, code string, ttl time.Duration) error {
	if code == "" {
		log.Error.Printf("EmailCodeRepository.Store: code is empty for %s", email)
		return errors.New("code cannot be empty")
	}
	key := EmailCodePrefix + email
	if err := r.redis.Set(ctx, key, code, ttl).Err(); err != nil {
		log.Error.Printf("EmailCodeRepository.Store: failed for %s: %v", email, err)
		return err
	}
	log.Info.Printf("EmailCodeRepository.Store: stored code=%s for %s, ttl=%v", code, email, ttl)
	return nil
}

func (r *EmailCodeRepository) Verify(ctx context.Context, email, code string) (bool, error) {
	key := EmailCodePrefix + email
	storedCode, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			log.Warn.Printf("EmailCodeRepository.Verify: code not found for %s", email)
			return false, nil
		}
		log.Error.Printf("EmailCodeRepository.Verify: failed for %s: %v", email, err)
		return false, err
	}
	log.Info.Printf("EmailCodeRepository.Verify: stored=%s, input=%s for %s", storedCode, code, email)
	match := storedCode == code
	if !match {
		log.Warn.Printf("EmailCodeRepository.Verify: code mismatch for %s", email)
	}
	return match, nil
}

func (r *EmailCodeRepository) Delete(ctx context.Context, email string) error {
	key := EmailCodePrefix + email
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		log.Error.Printf("EmailCodeRepository.Delete: failed for %s: %v", email, err)
		return err
	}
	log.Info.Printf("EmailCodeRepository.Delete: deleted code for %s", email)
	return nil
}
