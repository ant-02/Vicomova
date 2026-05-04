package query

import (
	"context"
	"testing"

	userEntity "vicomova/internal/user/domain/entity"
	"vicomova/internal/user/domain/valueobject"
	userRepo "vicomova/internal/user/domain/repository"
	"vicomova/internal/user/application/mock"
)

func newTestService(userRepo userRepo.UserRepository) *UserQueryService {
	return NewUserQueryService(userRepo)
}

func TestGetUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mock.NewMockUserRepository()

	username, _ := valueobject.NewUsername("testuser")
	email, _ := valueobject.NewEmail("test@example.com")
	user := &userEntity.User{
		ID:       1,
		Username: username,
		Email:    email,
	}
	userRepo.Users["testuser"] = user

	svc := newTestService(userRepo)

	result, err := svc.GetUser(ctx, &GetUserQuery{
		UserID:   1,
		Username: "testuser",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", result.Username)
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", result.Email)
	}
	if result.UserID != 1 {
		t.Errorf("expected userID 1, got %d", result.UserID)
	}
}

func TestGetUser_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := mock.NewMockUserRepository()
	svc := newTestService(userRepo)

	_, err := svc.GetUser(ctx, &GetUserQuery{
		UserID:   1,
		Username: "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error when user not found")
	}
}

func TestGetUser_InvalidUsername(t *testing.T) {
	ctx := context.Background()
	userRepo := mock.NewMockUserRepository()
	svc := newTestService(userRepo)

	_, err := svc.GetUser(ctx, &GetUserQuery{
		UserID:   1,
		Username: "", // Invalid empty username
	})
	if err == nil {
		t.Fatal("expected error for invalid username")
	}
}

func TestGetUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	userRepo := mock.NewMockUserRepository()
	userRepo.GetByUNErr = errInternal

	svc := newTestService(userRepo)

	_, err := svc.GetUser(ctx, &GetUserQuery{
		UserID:   1,
		Username: "testuser",
	})
	if err == nil {
		t.Fatal("expected error when repository fails")
	}
}

var errInternal = mockErr("internal error")

type mockErr string

func (e mockErr) Error() string { return string(e) }