package query

import (
	userRepo "vicomova/internal/user/domain/repository"
)

type UserQueryService struct {
	userRepo userRepo.UserRepository
}

func NewUserQueryService(userRepo userRepo.UserRepository) *UserQueryService {
	return &UserQueryService{userRepo: userRepo}
}
