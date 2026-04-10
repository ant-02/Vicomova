package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiry  = 30 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
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
		"exp":     time.Now().Add(AccessTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *TokenService) ParseAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *TokenService) GetAccessTokenExpiry() time.Duration {
	return AccessTokenExpiry
}

func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return RefreshTokenExpiry
}
