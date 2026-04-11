package redis

import (
	"context"
	"errors"
	"time"

	userRepo "vicomova/internal/user/domain/repository"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	"vicomova/internal/shared/pkg/log"
)

const emailCodePrefix = "email_code:"

var ErrCodeNotFound = errors.New("verification code not found")

type EmailCodeRepository struct{}

func NewEmailCodeRepository() userRepo.EmailCodeRepository {
	return &EmailCodeRepository{}
}

func (r *EmailCodeRepository) Store(ctx context.Context, email, code string, ttl time.Duration) error {
	key := emailCodePrefix + email
	if err := sharedRedis.GetClient().Set(ctx, key, code, ttl).Err(); err != nil {
		log.Error.Printf("EmailCodeRepository.Store: failed for %s: %v", email, err)
		return err
	}
	log.Info.Printf("EmailCodeRepository.Store: stored code for %s, ttl=%v", email, ttl)
	return nil
}

func (r *EmailCodeRepository) Verify(ctx context.Context, email, code string) (bool, error) {
	key := emailCodePrefix + email
	storedCode, err := sharedRedis.GetClient().Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			log.Warn.Printf("EmailCodeRepository.Verify: code not found for %s", email)
			return false, nil
		}
		log.Error.Printf("EmailCodeRepository.Verify: failed for %s: %v", email, err)
		return false, err
	}
	match := storedCode == code
	if !match {
		log.Warn.Printf("EmailCodeRepository.Verify: code mismatch for %s", email)
	}
	return match, nil
}

func (r *EmailCodeRepository) Delete(ctx context.Context, email string) error {
	key := emailCodePrefix + email
	if err := sharedRedis.GetClient().Del(ctx, key).Err(); err != nil {
		log.Error.Printf("EmailCodeRepository.Delete: failed for %s: %v", email, err)
		return err
	}
	log.Info.Printf("EmailCodeRepository.Delete: deleted code for %s", email)
	return nil
}