package command

import (
	userRepo "vicomova/internal/user/domain/repository"
	"vicomova/internal/user/domain/service"
	"vicomova/pkg/utils"
)

type UserCommandService struct {
	userRepo         userRepo.UserRepository
	refreshTokenRepo userRepo.RefreshTokenRepository
	emailCodeRepo    userRepo.EmailCodeRepository
	emailService     service.EmailService
	tokenService     *service.TokenService
	passwordHasher   utils.Hasher
}

func NewUserCommandService(
	userRepo userRepo.UserRepository,
	refreshTokenRepo userRepo.RefreshTokenRepository,
	emailCodeRepo userRepo.EmailCodeRepository,
	emailService service.EmailService,
	tokenService *service.TokenService,
	passwordHasher utils.Hasher,
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
