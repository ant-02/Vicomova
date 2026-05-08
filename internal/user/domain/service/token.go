package service

import (
	"time"

	"vicomova/pkg/constants"
	errorsPkg "vicomova/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	jwtSecret string
}

func NewTokenService(jwtSecret string) *TokenService {
	return &TokenService{jwtSecret: jwtSecret}
}

type AccessTokenClaims struct {
	UserID   int64
	Username string
}

func (s *TokenService) GenerateAccessToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(constants.AccessTokenExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *TokenService) ParseAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorsPkg.ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errorsPkg.ErrInvalidToken
}

func (s *TokenService) GetAccessTokenExpiry() time.Duration {
	return constants.AccessTokenExpiry
}

func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return constants.RefreshTokenExpiry
}
