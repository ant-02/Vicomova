package mock

import (
	"context"
	"time"
)

type MockEmailCodeRepository struct {
	Codes     map[string]string
	StoreErr  error
	VerifyErr error
	DeleteErr error
}

func NewMockEmailCodeRepository() *MockEmailCodeRepository {
	return &MockEmailCodeRepository{
		Codes: make(map[string]string),
	}
}

func (m *MockEmailCodeRepository) Store(ctx context.Context, email, code string, ttl time.Duration) error {
	if m.StoreErr != nil {
		return m.StoreErr
	}
	m.Codes[email] = code
	return nil
}

func (m *MockEmailCodeRepository) Verify(ctx context.Context, email, code string) (bool, error) {
	if m.VerifyErr != nil {
		return false, m.VerifyErr
	}
	stored, ok := m.Codes[email]
	if !ok {
		return false, nil
	}
	return stored == code, nil
}

func (m *MockEmailCodeRepository) Delete(ctx context.Context, email string) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.Codes, email)
	return nil
}
