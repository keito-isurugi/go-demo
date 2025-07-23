package algorithm

import (
	"fmt"
	"reflect"
	"testing"
)

// TestInsertionSort は挿入ソートの正常系テストです
func TestInsertionSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "通常の配列",
			input:    []int{64, 34, 25, 12, 22, 11, 90},
			expected: []int{11, 12, 22, 25, 34, 64, 90},
		},
		{
			name:     "逆順の配列",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "既にソート済みの配列",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "重複要素を含む配列",
			input:    []int{3, 5, 3, 5, 1, 3},
			expected: []int{1, 3, 3, 3, 5, 5},
		},
		{
			name:     "単一要素の配列",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "空の配列",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "負の数を含む配列",
			input:    []int{3, -1, 4, -5, 9, -2},
			expected: []int{-5, -2, -1, 3, 4, 9},
		},
		{
			name:     "大きな値と小さな値",
			input:    []int{1000, 1, 500, 2, 999},
			expected: []int{1, 2, 500, 999, 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の配列のコピーを作成（関数が元の配列を変更しないことを確認）
			originalInput := make([]int, len(tt.input))
			copy(originalInput, tt.input)
			
			result := InsertionSort(tt.input)
			
			// 結果の検証
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertionSort() = %v, want %v", result, tt.expected)
			}
			
			// 元の配列が変更されていないことを確認
			if !reflect.DeepEqual(tt.input, originalInput) {
				t.Errorf("InsertionSort() modified the original array: got %v, want %v", tt.input, originalInput)
			}
		})
	}
}

// BenchmarkInsertionSort はパフォーマンスベンチマークです
func BenchmarkInsertionSort(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// テストデータの準備
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = size - i // 逆順のデータ（最悪ケース）
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				InsertionSort(data)
			}
		})
	}
}

// TestInsertionSortStability は安定ソートであることを確認します
func TestInsertionSortStability(t *testing.T) {
	// 構造体を使って安定性をテスト
	type Item struct {
		value int
		index int
	}
	
	items := []Item{
		{3, 0}, {1, 1}, {3, 2}, {2, 3}, {3, 4},
	}
	
	// int配列に変換してソート
	values := make([]int, len(items))
	for i, item := range items {
		values[i] = item.value
	}
	
	sorted := InsertionSort(values)
	expected := []int{1, 2, 3, 3, 3}
	
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Stability test failed: got %v, want %v", sorted, expected)
	}
}