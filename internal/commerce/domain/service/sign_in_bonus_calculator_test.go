package service

import (
	"testing"
)

func TestSignInBonusCalculator_Calculate(t *testing.T) {
	calc := NewSignInBonusCalculator()

	tests := []struct {
		name            string
		consecutiveDays int
		wantBase        int64
		wantExtra       int64
		wantTotal       int64
	}{
		{"Day 1", 1, 10, 0, 10},
		{"Day 2", 2, 10, 0, 10},
		{"Day 3", 3, 15, 0, 15},
		{"Day 7", 7, 15, 5, 20}, // 7/7*5 = 5 extra bonus
		{"Day 8", 8, 20, 5, 25},
		{"Day 14", 14, 20, 10, 30},
		{"Day 15", 15, 30, 10, 40},
		{"Day 30", 30, 30, 20, 50},
		{"Day 31", 31, 50, 20, 70},
		{"Day 35", 35, 50, 25, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Calculate(tt.consecutiveDays)
			if got != tt.wantTotal {
				t.Errorf("Calculate(%d) = %d, want %d", tt.consecutiveDays, got, tt.wantTotal)
			}
		})
	}
}

func TestSignInBonusCalculator_GetBaseBonus(t *testing.T) {
	calc := NewSignInBonusCalculator()
	if got := calc.GetBaseBonus(); got != 10 {
		t.Errorf("GetBaseBonus() = %d, want 10", got)
	}
}

func TestSignInBonusCalculator_MaxBonus(t *testing.T) {
	calc := NewSignInBonusCalculator()
	if got := calc.MaxBonus(); got != 100 {
		t.Errorf("MaxBonus() = %d, want 100", got)
	}
}
