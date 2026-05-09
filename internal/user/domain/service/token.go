package service

import (
	"time"

	"vicomova/pkg/constants"
	"vicomova/pkg/utils"
)

type TokenService struct {
	jwtSecret string
}

func NewTokenService(jwtSecret string) *TokenService {
	return &TokenService{jwtSecret: jwtSecret}
}

func (s *TokenService) GenerateAccessToken(userID int64, username string) (string, error) {
	return utils.GenerateToken(userID, username, s.jwtSecret, constants.AccessTokenExpiry)
}

func (s *TokenService) ParseAccessToken(tokenString string) (map[string]interface{}, error) {
	return utils.ParseToken(tokenString, s.jwtSecret)
}

func (s *TokenService) GetAccessTokenExpiry() time.Duration {
	return constants.AccessTokenExpiry
}

func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return constants.RefreshTokenExpiry
}
