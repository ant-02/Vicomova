package valueobject

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	EmailCodeTTL = 10 * time.Minute
	EmailCodeLength = 6
)

type EmailCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

func GenerateEmailCode(email string) *EmailCode {
	code := generateRandomCode()
	return &EmailCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(EmailCodeTTL),
	}
}

func generateRandomCode() string {
	code := make([]byte, EmailCodeLength)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = byte('0' + n.Int64())
	}
	return fmt.Sprintf("%s", code)
}

func (e *EmailCode) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}