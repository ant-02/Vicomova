package mock

import (
	"context"

	userVO "vicomova/internal/user/domain/valueobject"
)

type MockRefreshTokenRepository struct {
	Tokens    map[string]*userVO.RefreshToken
	CreateErr error
	GetErr    error
	RevokeErr error
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{
		Tokens: make(map[string]*userVO.RefreshToken),
	}
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, rt *userVO.RefreshToken) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.Tokens[rt.Token] = rt
	return nil
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*userVO.RefreshToken, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Tokens[token], nil
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	if m.RevokeErr != nil {
		return m.RevokeErr
	}
	if rt, ok := m.Tokens[token]; ok {
		rt.Revoked = true
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	return nil
}
