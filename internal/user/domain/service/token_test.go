package service

import (
	"testing"

	"vicomova/pkg/constants"
)

func TestGenerateAccessToken(t *testing.T) {
	svc := NewTokenService("test-secret")

	token, err := svc.GenerateAccessToken(123, "testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Verify we can parse it
	claims, err := svc.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("expected to parse generated token, got %v", err)
	}
	if claims["user_id"] != float64(123) {
		t.Errorf("expected user_id 123, got %v", claims["user_id"])
	}
	if claims["username"] != "testuser" {
		t.Errorf("expected username testuser, got %v", claims["username"])
	}
}

func TestParseAccessToken_Valid(t *testing.T) {
	svc := NewTokenService("test-secret")

	token, _ := svc.GenerateAccessToken(456, "alice")

	claims, err := svc.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims["user_id"] != float64(456) {
		t.Errorf("expected user_id 456, got %v", claims["user_id"])
	}
}

func TestParseAccessToken_Invalid(t *testing.T) {
	svc := NewTokenService("test-secret")

	_, err := svc.ParseAccessToken("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestParseAccessToken_WrongSecret(t *testing.T) {
	svc1 := NewTokenService("secret1")
	svc2 := NewTokenService("secret2")

	token, _ := svc1.GenerateAccessToken(1, "user")

	_, err := svc2.ParseAccessToken(token)
	if err == nil {
		t.Error("expected error when parsing with wrong secret")
	}
}

func TestGetAccessTokenExpiry(t *testing.T) {
	svc := NewTokenService("test-secret")

	expiry := svc.GetAccessTokenExpiry()
	if expiry != constants.AccessTokenExpiry {
		t.Errorf("expected %v, got %v", constants.AccessTokenExpiry, expiry)
	}
}

func TestGetRefreshTokenExpiry(t *testing.T) {
	svc := NewTokenService("test-secret")

	expiry := svc.GetRefreshTokenExpiry()
	if expiry != constants.RefreshTokenExpiry {
		t.Errorf("expected %v, got %v", constants.RefreshTokenExpiry, expiry)
	}
}
