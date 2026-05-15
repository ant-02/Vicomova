package service

import "context"

// PaymentService 第三方支付服务接口（预留）
type PaymentService interface {
	// CreatePaymentLink 创建支付链接，返回支付链接地址
	CreatePaymentLink(ctx context.Context, userID int64, amount float64) (string, error)
	// VerifyPayment 验证支付结果，返回是否成功
	VerifyPayment(ctx context.Context, paymentID string) (bool, error)
}

// MockPaymentService Mock 实现，用于开发测试
type MockPaymentService struct{}

func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{}
}

func (s *MockPaymentService) CreatePaymentLink(ctx context.Context, userID int64, amount float64) (string, error) {
	// 预留接口，初期直接返回 mock 链接
	return "https://mock-payment.example.com/pay", nil
}

func (s *MockPaymentService) VerifyPayment(ctx context.Context, paymentID string) (bool, error) {
	// 预留接口，初期直接返回 true
	return true, nil
}
