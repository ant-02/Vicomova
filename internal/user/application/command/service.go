package command

import (
	userRepo "vicomova/internal/user/domain/repository"
	"vicomova/internal/user/domain/service"
	infraEmail "vicomova/internal/user/infrastructure/external/email"
	sharedHasher "vicomova/internal/shared/pkg/hasher"
)

type UserCommandService struct {
	userRepo         userRepo.UserRepository
	refreshTokenRepo userRepo.RefreshTokenRepository
	emailCodeRepo    userRepo.EmailCodeRepository
	emailService     infraEmail.EmailService
	tokenService     *service.TokenService
	passwordHasher   sharedHasher.Hasher
}

func NewUserCommandService(
	userRepo userRepo.UserRepository,
	refreshTokenRepo userRepo.RefreshTokenRepository,
	emailCodeRepo userRepo.EmailCodeRepository,
	emailService infraEmail.EmailService,
	tokenService *service.TokenService,
	passwordHasher sharedHasher.Hasher,
) *UserCommandService {
	return &UserCommandService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		emailCodeRepo:    emailCodeRepo,
		emailService:     emailService,
		tokenService:     tokenService,
		passwordHasher:   passwordHasher,
	}
}
