package algorithm

import (
	"testing"
	"reflect"
)

func TestFibonacciRecursive(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "Fib(0)",
			input:    0,
			expected: 0,
		},
		{
			name:     "Fib(1)",
			input:    1,
			expected: 1,
		},
		{
			name:     "Fib(2)",
			input:    2,
			expected: 1,
		},
		{
			name:     "Fib(3)",
			input:    3,
			expected: 2,
		},
		{
			name:     "Fib(4)",
			input:    4,
			expected: 3,
		},
		{
			name:     "Fib(5)",
			input:    5,
			expected: 5,
		},
		{
			name:     "Fib(6)",
			input:    6,
			expected: 8,
		},
		{
			name:     "Fib(10)",
			input:    10,
			expected: 55,
		},
		{
			name:     "Fib(15)",
			input:    15,
			expected: 610,
		},
		{
			name:     "Fib(20)",
			input:    20,
			expected: 6765,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciRecursive(tt.input)
			if result != tt.expected {
				t.Errorf("FibonacciRecursive(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciMemoization(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "Fib(0)",
			input:    0,
			expected: 0,
		},
		{
			name:     "Fib(1)",
			input:    1,
			expected: 1,
		},
		{
			name:     "Fib(10)",
			input:    10,
			expected: 55,
		},
		{
			name:     "Fib(20)",
			input:    20,
			expected: 6765,
		},
		{
			name:     "Fib(30)",
			input:    30,
			expected: 832040,
		},
		{
			name:     "Fib(40)",
			input:    40,
			expected: 102334155,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciMemoization(tt.input)
			if result != tt.expected {
				t.Errorf("FibonacciMemoization(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciIterative(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "Fib(0)",
			input:    0,
			expected: 0,
		},
		{
			name:     "Fib(1)",
			input:    1,
			expected: 1,
		},
		{
			name:     "Fib(10)",
			input:    10,
			expected: 55,
		},
		{
			name:     "Fib(20)",
			input:    20,
			expected: 6765,
		},
		{
			name:     "Fib(50)",
			input:    50,
			expected: 12586269025,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciIterative(tt.input)
			if result != tt.expected {
				t.Errorf("FibonacciIterative(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []int
	}{
		{
			name:     "最初の0項",
			input:    0,
			expected: []int{0},
		},
		{
			name:     "最初の1項",
			input:    1,
			expected: []int{0, 1},
		},
		{
			name:     "最初の5項",
			input:    5,
			expected: []int{0, 1, 1, 2, 3, 5},
		},
		{
			name:     "最初の10項",
			input:    10,
			expected: []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55},
		},
		{
			name:     "負の値",
			input:    -1,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciSequence(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FibonacciSequence(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciSequenceIterative(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []int
	}{
		{
			name:     "最初の5項",
			input:    5,
			expected: []int{0, 1, 1, 2, 3, 5},
		},
		{
			name:     "最初の10項",
			input:    10,
			expected: []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciSequenceIterative(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FibonacciSequenceIterative(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// 各実装が同じ結果を返すことを確認
func TestFibonacciConsistency(t *testing.T) {
	for n := 0; n <= 20; n++ {
		recursive := FibonacciRecursive(n)
		memoized := FibonacciMemoization(n)
		iterative := FibonacciIterative(n)

		if recursive != memoized {
			t.Errorf("Fib(%d): recursive=%d, memoized=%d", n, recursive, memoized)
		}

		if recursive != iterative {
			t.Errorf("Fib(%d): recursive=%d, iterative=%d", n, recursive, iterative)
		}
	}
}

// ベンチマークテスト
func BenchmarkFibonacciRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciRecursive(20)
	}
}

func BenchmarkFibonacciMemoization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciMemoization(20)
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}

// 大きな値でのベンチマーク（再帰は除く）
func BenchmarkFibonacciMemoizationLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciMemoization(100)
	}
}

func BenchmarkFibonacciIterativeLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(100)
	}
}

// 異なるサイズでのベンチマーク比較
func BenchmarkFibonacciComparison10(b *testing.B) {
	b.Run("Recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciRecursive(10)
		}
	})
	
	b.Run("Memoization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciMemoization(10)
		}
	})
	
	b.Run("Iterative", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciIterative(10)
		}
	})
}

func BenchmarkFibonacciComparison30(b *testing.B) {
	b.Run("Recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciRecursive(30)
		}
	})
	
	b.Run("Memoization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciMemoization(30)
		}
	})
	
	b.Run("Iterative", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FibonacciIterative(30)
		}
	})
}

// メモ化の効果を測定するテスト
func TestMemoizationEffectiveness(t *testing.T) {
	// メモ化の効果を確認するため、同じ値を複数回計算
	memo := make(map[int]int)
	
	// 最初の計算
	result1 := fibMemo(10, memo)
	
	// メモが正しく機能していることを確認
	if len(memo) == 0 {
		t.Error("メモ化が機能していません")
	}
	
	// 2回目の計算（メモから取得されるはず）
	result2 := fibMemo(10, memo)
	
	if result1 != result2 {
		t.Errorf("メモ化の結果が一致しません: %d != %d", result1, result2)
	}
	
	// より大きな値でもメモが使われることを確認
	result3 := fibMemo(15, memo)
	expected := 610
	
	if result3 != expected {
		t.Errorf("fibMemo(15) = %d, want %d", result3, expected)
	}
}