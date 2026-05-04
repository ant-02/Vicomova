package valueobject

import (
	"testing"
)

func TestNewEmail_Valid(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.org",
		"user+tag@example.com",
		"a@b.co",
	}

	for _, e := range validEmails {
		_, err := NewEmail(e)
		if err != nil {
			t.Errorf("expected valid email %q, got error: %v", e, err)
		}
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	invalidEmails := []string{
		"notanemail",
		"missing@",
		"@nodomain.com",
		"spaces in@email.com",
		"invalid@.com",
	}

	for _, e := range invalidEmails {
		_, err := NewEmail(e)
		if err == nil {
			t.Errorf("expected error for invalid email %q", e)
		}
	}
}

func TestEmail_String(t *testing.T) {
	e, _ := NewEmail("test@example.com")
	if e.String() != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %q", e.String())
	}
}

func TestEmail_Equal(t *testing.T) {
	e1, _ := NewEmail("a@test.com")
	e2, _ := NewEmail("a@test.com")
	e3, _ := NewEmail("b@test.com")

	if !e1.Equal(e2) {
		t.Error("expected equal emails to be equal")
	}
	if e1.Equal(e3) {
		t.Error("expected different emails to not be equal")
	}
}