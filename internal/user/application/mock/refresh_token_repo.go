package mock

import (
	"context"
	"time"
)

type MockRefreshTokenRepository struct {
	Tokens    map[string]int64 // token -> userID
	CreateErr error
	GetErr    error
	RevokeErr error
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{
		Tokens: make(map[string]int64),
	}
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.Tokens[token] = userID
	return nil
}

func (m *MockRefreshTokenRepository) GetUserID(ctx context.Context, token string) (int64, error) {
	if m.GetErr != nil {
		return 0, m.GetErr
	}
	return m.Tokens[token], nil
}

func (m *MockRefreshTokenRepository) Exists(ctx context.Context, token string) (bool, error) {
	_, ok := m.Tokens[token]
	return ok, nil
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	if m.RevokeErr != nil {
		return m.RevokeErr
	}
	delete(m.Tokens, token)
	return nil
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	return nil
}