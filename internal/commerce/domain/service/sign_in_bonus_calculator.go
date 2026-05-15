package service

import (
	"math"
)

// SignInBonusCalculator 签到积分计算器
type SignInBonusCalculator struct{}

func NewSignInBonusCalculator() *SignInBonusCalculator {
	return &SignInBonusCalculator{}
}

// Calculate 计算签到积分
// 规则：
//   - 连续 1-2 天：基础 10 积分
//   - 连续 3-7 天：基础 15 积分
//   - 连续 8-14 天：基础 20 积分
//   - 连续 15-30 天：基础 30 积分
//   - 连续 31+ 天：基础 50 积分
//   - 额外加成：每满 7 天加 5 积分
func (c *SignInBonusCalculator) Calculate(consecutiveDays int) int64 {
	baseBonus := c.getBaseBonus(consecutiveDays)
	extraBonus := int64(consecutiveDays / 7 * 5)
	return baseBonus + extraBonus
}

func (c *SignInBonusCalculator) getBaseBonus(consecutiveDays int) int64 {
	switch {
	case consecutiveDays <= 2:
		return 10
	case consecutiveDays <= 7:
		return 15
	case consecutiveDays <= 14:
		return 20
	case consecutiveDays <= 30:
		return 30
	default:
		return 50
	}
}

// GetBaseBonus 获取基础积分（用于显示）
func (c *SignInBonusCalculator) GetBaseBonus() int64 {
	return 10
}

// MaxBonus 最大积分上限
func (c *SignInBonusCalculator) MaxBonus() int64 {
	return 100
}

var _ interface {
	Calculate(consecutiveDays int) int64
	GetBaseBonus() int64
	MaxBonus() int64
} = (*SignInBonusCalculator)(nil)

var _ = math.Max // import check
