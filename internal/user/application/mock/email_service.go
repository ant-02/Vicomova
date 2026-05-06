package mock

import (
	"context"
)

type MockEmailService struct {
	SendErr    error
	SendCalled bool
	SentTo     string
	SentCode   string
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (m *MockEmailService) SendVerificationCode(ctx context.Context, toEmail, code string) error {
	m.SendCalled = true
	m.SentTo = toEmail
	m.SentCode = code
	return m.SendErr
}
