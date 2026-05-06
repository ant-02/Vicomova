package command

import (
	"context"
	"errors"
	"testing"

	"vicomova/internal/shared/pkg/constants"
	sharedHasher "vicomova/internal/shared/pkg/hasher"
	"vicomova/internal/user/application/mock"
	userEntity "vicomova/internal/user/domain/entity"
	userRepo "vicomova/internal/user/domain/repository"
	"vicomova/internal/user/domain/service"
	"vicomova/internal/user/domain/valueobject"
	infraEmail "vicomova/internal/user/infrastructure/external/email"
)

// Helper to create test service
func newTestService(
	userRepo userRepo.UserRepository,
	emailCodeRepo userRepo.EmailCodeRepository,
	emailService infraEmail.EmailService,
) *UserCommandService {
	tokenSvc := service.NewTokenService("test-secret")
	hasher := &sharedHasher.SHA256Hasher{}
	return NewUserCommandService(
		userRepo,
		mock.NewMockRefreshTokenRepository(),
		emailCodeRepo,
		emailService,
		tokenSvc,
		hasher,
	)
}

// ===== Tests for SendVerificationCode =====

func TestSendVerificationCode_Success(t *testing.T) {
	ctx := context.Background()
	svc := newTestService(
		mock.NewMockUserRepository(),
		mock.NewMockEmailCodeRepository(),
		mock.NewMockEmailService(),
	)

	cmd := &SendVerificationCodeCommand{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	result, err := svc.SendVerificationCode(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
}

func TestSendVerificationCode_UserAlreadyExists(t *testing.T) {
	ctx := context.Background()
	userRepo := mock.NewMockUserRepository()
	username, _ := valueobject.NewUsername("testuser")
	email, _ := valueobject.NewEmail("existing@example.com")
	userRepo.Users["testuser"] = &userEntity.User{
		Username: username,
		Email:    email,
	}

	svc := newTestService(
		userRepo,
		mock.NewMockEmailCodeRepository(),
		mock.NewMockEmailService(),
	)

	cmd := &SendVerificationCodeCommand{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	result, err := svc.SendVerificationCode(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Success {
		t.Error("expected success=false for existing user")
	}
	if result.Message != "username already exists" {
		t.Errorf("expected message 'username already exists', got '%s'", result.Message)
	}
}

func TestSendVerificationCode_EmailServiceError(t *testing.T) {
	ctx := context.Background()
	emailSvc := mock.NewMockEmailService()
	emailSvc.SendErr = errors.New("email service unavailable")

	svc := newTestService(
		mock.NewMockUserRepository(),
		mock.NewMockEmailCodeRepository(),
		emailSvc,
	)

	cmd := &SendVerificationCodeCommand{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	_, err := svc.SendVerificationCode(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when email service fails")
	}
}

func TestSendVerificationCode_CodeStoredInRedis(t *testing.T) {
	ctx := context.Background()
	emailCodeRepo := mock.NewMockEmailCodeRepository()

	svc := newTestService(
		mock.NewMockUserRepository(),
		emailCodeRepo,
		mock.NewMockEmailService(),
	)

	cmd := &SendVerificationCodeCommand{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	_, err := svc.SendVerificationCode(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify code was stored (6 digits)
	storedCode, ok := emailCodeRepo.Codes["test@example.com"]
	if !ok {
		t.Fatal("code was not stored in repository")
	}
	if len(storedCode) != constants.EmailCodeLength {
		t.Errorf("expected code length %d, got %d", constants.EmailCodeLength, len(storedCode))
	}
}

// ===== Tests for VerifyAndRegister =====

func TestVerifyAndRegister_Success(t *testing.T) {
	ctx := context.Background()
	emailCodeRepo := mock.NewMockEmailCodeRepository()
	emailCodeRepo.Codes["test@example.com"] = "123456"

	svc := newTestService(
		mock.NewMockUserRepository(),
		emailCodeRepo,
		mock.NewMockEmailService(),
	)

	cmd := &VerifyAndRegisterCommand{
		Username: "newuser",
		Password: "password123",
		Email:    "test@example.com",
		Code:     "123456",
	}

	result, err := svc.VerifyAndRegister(ctx, cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", result.Username)
	}
	if result.UserID == 0 {
		t.Error("expected non-zero user ID")
	}

	// Verify code was deleted after use
	if _, ok := emailCodeRepo.Codes["test@example.com"]; ok {
		t.Error("code was not deleted after successful registration")
	}
}

func TestVerifyAndRegister_InvalidCode(t *testing.T) {
	ctx := context.Background()
	emailCodeRepo := mock.NewMockEmailCodeRepository()
	emailCodeRepo.Codes["test@example.com"] = "123456"

	svc := newTestService(
		mock.NewMockUserRepository(),
		emailCodeRepo,
		mock.NewMockEmailService(),
	)

	cmd := &VerifyAndRegisterCommand{
		Username: "newuser",
		Password: "password123",
		Email:    "test@example.com",
		Code:     "000000", // Wrong code
	}

	_, err := svc.VerifyAndRegister(ctx, cmd)
	if err == nil {
		t.Fatal("expected error for invalid code")
	}
}

func TestVerifyAndRegister_CodeNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newTestService(
		mock.NewMockUserRepository(),
		mock.NewMockEmailCodeRepository(),
		mock.NewMockEmailService(),
	)

	cmd := &VerifyAndRegisterCommand{
		Username: "newuser",
		Password: "password123",
		Email:    "nonexistent@example.com",
		Code:     "123456",
	}

	_, err := svc.VerifyAndRegister(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when code not found")
	}
}

func TestVerifyAndRegister_UserAlreadyExists(t *testing.T) {
	ctx := context.Background()
	emailCodeRepo := mock.NewMockEmailCodeRepository()
	emailCodeRepo.Codes["test@example.com"] = "123456"

	userRepo := mock.NewMockUserRepository()
	existingUsername, _ := valueobject.NewUsername("existinguser")
	existingEmail, _ := valueobject.NewEmail("other@example.com")
	userRepo.Users["existinguser"] = &userEntity.User{
		Username: existingUsername,
		Email:    existingEmail,
	}

	svc := newTestService(
		userRepo,
		emailCodeRepo,
		mock.NewMockEmailService(),
	)

	cmd := &VerifyAndRegisterCommand{
		Username: "existinguser",
		Password: "password123",
		Email:    "test@example.com",
		Code:     "123456",
	}

	_, err := svc.VerifyAndRegister(ctx, cmd)
	if err == nil {
		t.Fatal("expected error when user already exists")
	}
}

func TestVerifyAndRegister_CodeDeletedAfterSuccess(t *testing.T) {
	ctx := context.Background()
	emailCodeRepo := mock.NewMockEmailCodeRepository()
	emailCodeRepo.Codes["test@example.com"] = "123456"

	svc := newTestService(
		mock.NewMockUserRepository(),
		emailCodeRepo,
		mock.NewMockEmailService(),
	)

	cmd := &VerifyAndRegisterCommand{
		Username: "newuser",
		Password: "password123",
		Email:    "test@example.com",
		Code:     "123456",
	}

	_, err := svc.VerifyAndRegister(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Code should be deleted
	if _, exists := emailCodeRepo.Codes["test@example.com"]; exists {
		t.Error("code should be deleted after successful registration")
	}
}
