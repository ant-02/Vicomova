package valueobject

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrInvalidPassword = errors.New("invalid password")

type Password struct {
	hash string
}

func NewPassword(plain string) (*Password, error) {
	if err := validatePassword(plain); err != nil {
		return nil, err
	}
	return &Password{hash: hashPassword(plain)}, nil
}

func validatePassword(plain string) error {
	if len(plain) < 6 || len(plain) > 128 {
		return errors.New("password must be between 6 and 128 characters")
	}
	return nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (p *Password) Verify(plain string) bool {
	return p.hash == hashPassword(plain)
}

func (p *Password) Hash() string {
	return p.hash
}

func (p *Password) Equal(other *Password) bool {
	return p.hash == other.hash
}

// NewPasswordFromHash creates a Password from an existing hash (for repository layer).
func NewPasswordFromHash(hash string) *Password {
	return &Password{hash: hash}
}