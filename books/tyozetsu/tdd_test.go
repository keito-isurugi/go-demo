package tyozetsu

import (
	"fmt"
	"testing"
)

func TestMath_min(t *testing.T) {
	math := &Math{}
	var tests = []struct {
		a, b, want int
	}{
		{1, 2, 1},
		{10, 9, 9},
		{-1, -2, -2},
		{0, 0, 0},
		{100, 50, 50},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("min(%d,%d)", tt.a, tt.b), func(t *testing.T) {
			got := math.min(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("math.min(%d,%d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMath_max(t *testing.T) {
	math := &Math{}
	var tests = []struct {
		a, b, want int
	}{
		{1, 2, 2},
		{10, 9, 10},
		{-1, -2, -1},
		{0, 0, 0},
		{100, 50, 100},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("max(%d,%d)", tt.a, tt.b), func(t *testing.T) {
			got := math.max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("math.max(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMath_saturate(t *testing.T) {
	mathUtinl := &MathUtil{}
	var tests = []struct {
		value, minValue, maxValue, want int
	} {
		{2, 1, 3, 2},
		{0, 1, 3, 1},
		{4, 1, 3, 3},
		{1, 1, 3, 1},
		{3, 1, 3, 3},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("saturate(%d,%d,%d)", tt.value, tt.minValue, tt.maxValue), func(t *testing.T) {
			got := mathUtinl.saturate(tt.value, tt.minValue, tt.maxValue)
			if got != tt.want {
				t.Errorf("saturate(%d,%d,%d) = %d; want %d", tt.value, tt.minValue, tt.maxValue, got, tt.want)
			}
		})
	}
}