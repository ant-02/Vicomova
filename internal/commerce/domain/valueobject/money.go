package valueobject

import "errors"

var ErrNegativeAmount = errors.New("amount must be positive")

// Money 金额（积分或钱包）
type Money struct {
	amount int64
	unit   string // "points" or "wallet"
}

func NewPoints(amount int64) (*Money, error) {
	if amount < 0 {
		return nil, ErrNegativeAmount
	}
	return &Money{amount: amount, unit: "points"}, nil
}

func NewWallet(amount float64) (*Money, error) {
	if amount < 0 {
		return nil, ErrNegativeAmount
	}
	milli := int64(amount * 1000)
	return &Money{amount: milli, unit: "wallet"}, nil
}

func (m *Money) Amount() int64 {
	return m.amount
}

func (m *Money) Unit() string {
	return m.unit
}

func (m *Money) IsPoints() bool {
	return m.unit == "points"
}

func (m *Money) IsWallet() bool {
	return m.unit == "wallet"
}

func (m *Money) Equal(other *Money) bool {
	return m.amount == other.amount && m.unit == other.unit
}
