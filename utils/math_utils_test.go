package utils

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"with zero", 5, 0, 5},
		{"both zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 2},
		{"negative numbers", -5, -3, -2},
		{"mixed signs", 5, -3, 8},
		{"with zero", 5, 0, 5},
		{"result negative", 3, 5, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 3, 4, 12},
		{"negative numbers", -3, -4, 12},
		{"mixed signs", -3, 4, -12},
		{"with zero", 5, 0, 0},
		{"with one", 5, 1, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		expected    int
		expectError bool
	}{
		{"positive numbers", 12, 3, 4, false},
		{"negative numbers", -12, -3, 4, false},
		{"mixed signs", -12, 3, -4, false},
		{"division by zero", 5, 0, 0, true},
		{"zero divided by number", 0, 5, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if tt.expectError {
				if err == nil {
					t.Errorf("Divide(%d, %d) expected error but got none", tt.a, tt.b)
				}
				if err != ErrDivisionByZero {
					t.Errorf("Divide(%d, %d) expected ErrDivisionByZero but got %v", tt.a, tt.b, err)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{"even positive", 4, true},
		{"odd positive", 5, false},
		{"even negative", -4, true},
		{"odd negative", -5, false},
		{"zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEven(tt.n)
			if result != tt.expected {
				t.Errorf("IsEven(%d) = %t; expected %t", tt.n, result, tt.expected)
			}
		})
	}
}

func TestIsPositive(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{"positive number", 5, true},
		{"negative number", -5, false},
		{"zero", 0, false},
		{"large positive", 1000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPositive(tt.n)
			if result != tt.expected {
				t.Errorf("IsPositive(%d) = %t; expected %t", tt.n, result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"first larger", 5, 3, 5},
		{"second larger", 3, 5, 5},
		{"equal", 5, 5, 5},
		{"negative numbers", -3, -5, -3},
		{"mixed signs", -3, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Max(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"first smaller", 3, 5, 3},
		{"second smaller", 5, 3, 3},
		{"equal", 5, 5, 5},
		{"negative numbers", -3, -5, -5},
		{"mixed signs", -3, 5, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// ベンチマークテスト
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(100, 200)
	}
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Divide(100, 20)
	}
}
