package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hasher defines the interface for password hashing implementations.
type Hasher interface {
	Hash(password string) string
	Verify(password, hash string) bool
}

// SHA256Hasher implements Hasher using SHA256.
type SHA256Hasher struct{}

func (h *SHA256Hasher) Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (h *SHA256Hasher) Verify(password, hash string) bool {
	return h.Hash(password) == hash
}
