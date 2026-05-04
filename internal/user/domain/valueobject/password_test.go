package valueobject

import (
	"testing"
)

type mockHasher struct {
	hash       string
	verifyPass bool
}

func (m *mockHasher) Hash(password string) string {
	return m.hash
}

func (m *mockHasher) Verify(password, hash string) bool {
	return m.verifyPass
}

func TestNewPasswordFromPlain_Valid(t *testing.T) {
	hasher := &mockHasher{hash: "hashed_password"}
	p, err := NewPasswordFromPlain("password123", hasher)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Hash() != "hashed_password" {
		t.Errorf("expected hash 'hashed_password', got %q", p.Hash())
	}
}

func TestNewPasswordFromPlain_TooShort(t *testing.T) {
	hasher := &mockHasher{hash: "hashed"}
	_, err := NewPasswordFromPlain("12345", hasher)
	if err != ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestNewPasswordFromPlain_MaxLength(t *testing.T) {
	hasher := &mockHasher{hash: "hashed"}
	// 128 chars should be valid
	validPass := "aaaaa" + string(make([]byte, 123))
	_, err := NewPasswordFromPlain(validPass, hasher)
	if err != nil {
		t.Errorf("expected valid for 128 char password, got %v", err)
	}
}

func TestNewPasswordFromPlain_TooLong(t *testing.T) {
	hasher := &mockHasher{hash: "hashed"}
	// 129 chars should be invalid
	invalidPass := "a" + string(make([]byte, 129))
	_, err := NewPasswordFromPlain(invalidPass, hasher)
	if err != ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestPassword_Hash(t *testing.T) {
	hasher := &mockHasher{hash: "secret_hash"}
	p, _ := NewPasswordFromPlain("password", hasher)
	if p.Hash() != "secret_hash" {
		t.Errorf("expected 'secret_hash', got %q", p.Hash())
	}
}

func TestPassword_Equal(t *testing.T) {
	hasher := &mockHasher{hash: "hash1"}
	p1, _ := NewPasswordFromPlain("password1", hasher) // valid 9-char password

	p3 := NewPassword("hash1")
	p4 := NewPassword("hash2")

	if !p1.Equal(p3) {
		t.Error("expected same hash to be equal")
	}
	if p1.Equal(p4) {
		t.Error("expected different hash to not be equal")
	}
}

func TestNewPassword(t *testing.T) {
	p := NewPassword("some_hash")
	if p.Hash() != "some_hash" {
		t.Errorf("expected 'some_hash', got %q", p.Hash())
	}
}