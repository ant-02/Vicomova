package valueobject

import (
	"errors"
)

var ErrPasswordTooShort = errors.New("password must be between 6 and 128 characters")

type Password struct {
	hash string
}

// NewPassword creates a Password with a pre-hashed value (for repository layer).
func NewPassword(hash string) *Password {
	return &Password{hash: hash}
}

// NewPasswordFromPlain creates a Password from plain text with validation.
func NewPasswordFromPlain(plain string, hasher Hasher) (*Password, error) {
	if len(plain) < 6 || len(plain) > 128 {
		return nil, ErrPasswordTooShort
	}
	return &Password{hash: hasher.Hash(plain)}, nil
}

// Hasher interface for password hashing.
type Hasher interface {
	Hash(password string) string
	Verify(password, hash string) bool
}

func (p *Password) Hash() string {
	return p.hash
}

func (p *Password) Equal(other *Password) bool {
	return p.hash == other.hash
}
