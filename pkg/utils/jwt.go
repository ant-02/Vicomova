package utils

import (
	"time"

	errorsPkg "vicomova/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64, username string, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(expiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errorsPkg.ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errorsPkg.ErrInvalidToken
	}

	return claims, nil
}

func GetUserIDFromClaims(claims jwt.MapClaims) (int64, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errorsPkg.ErrInvalidToken
	}
	return int64(userID), nil
}
