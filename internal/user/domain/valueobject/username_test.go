package valueobject

import (
	"testing"
)

func TestNewUsername_Valid(t *testing.T) {
	validUsernames := []string{"abc", "user123", "test_user", "TestUser", "user1234567890123456789012345"}

	for _, u := range validUsernames {
		_, err := NewUsername(u)
		if err != nil {
			t.Errorf("expected valid username %q, got error: %v", u, err)
		}
	}
}

func TestNewUsername_TooShort(t *testing.T) {
	_, err := NewUsername("ab")
	if err == nil {
		t.Error("expected error for username < 3 chars")
	}
}

func TestNewUsername_TooLong(t *testing.T) {
	_, err := NewUsername("abcdefghijklmnopqrstuvwxyz12345") // 33 chars
	if err == nil {
		t.Error("expected error for username > 32 chars")
	}
}

func TestNewUsername_InvalidChars(t *testing.T) {
	invalid := []string{"user@name", "user name", "user-name", "user.name"}

	for _, u := range invalid {
		_, err := NewUsername(u)
		if err == nil {
			t.Errorf("expected error for invalid username %q", u)
		}
	}
}

func TestUsername_String(t *testing.T) {
	u, _ := NewUsername("testuser")
	if u.String() != "testuser" {
		t.Errorf("expected 'testuser', got %q", u.String())
	}
}

func TestUsername_Equal(t *testing.T) {
	u1, _ := NewUsername("user1")
	u2, _ := NewUsername("user1")
	u3, _ := NewUsername("user2")

	if !u1.Equal(u2) {
		t.Error("expected equal usernames to be equal")
	}
	if u1.Equal(u3) {
		t.Error("expected different usernames to not be equal")
	}
}