package algorithm

import (
	"fmt"
	"reflect"
	"testing"
)

// TestSelectionSort は、SelectionSort関数の正常動作をテストします
func TestSelectionSort(t *testing.T) {
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
			name:     "すでにソート済みの配列",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "逆順の配列",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "重複要素を含む配列",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5},
			expected: []int{1, 1, 2, 3, 4, 5, 5, 6, 9},
		},
		{
			name:     "単一要素の配列",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "2要素の配列",
			input:    []int{2, 1},
			expected: []int{1, 2},
		},
		{
			name:     "空の配列",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "負の数を含む配列",
			input:    []int{-5, 3, -1, 0, 2, -3},
			expected: []int{-5, -3, -1, 0, 2, 3},
		},
		{
			name:     "すべて同じ要素の配列",
			input:    []int{7, 7, 7, 7, 7},
			expected: []int{7, 7, 7, 7, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の配列を保持するためにコピーを作成
			arr := make([]int, len(tt.input))
			copy(arr, tt.input)
			
			SelectionSort(arr)
			
			if !reflect.DeepEqual(arr, tt.expected) {
				t.Errorf("SelectionSort(%v) = %v, want %v", tt.input, arr, tt.expected)
			}
		})
	}
}

// TestSelectionSortDescending は、SelectionSortDescending関数の正常動作をテストします
func TestSelectionSortDescending(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "通常の配列",
			input:    []int{64, 34, 25, 12, 22, 11, 90},
			expected: []int{90, 64, 34, 25, 22, 12, 11},
		},
		{
			name:     "すでに降順の配列",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "昇順の配列",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "重複要素を含む配列",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5},
			expected: []int{9, 6, 5, 5, 4, 3, 2, 1, 1},
		},
		{
			name:     "負の数を含む配列",
			input:    []int{-5, 3, -1, 0, 2, -3},
			expected: []int{3, 2, 0, -1, -3, -5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の配列を保持するためにコピーを作成
			arr := make([]int, len(tt.input))
			copy(arr, tt.input)
			
			SelectionSortDescending(arr)
			
			if !reflect.DeepEqual(arr, tt.expected) {
				t.Errorf("SelectionSortDescending(%v) = %v, want %v", tt.input, arr, tt.expected)
			}
		})
	}
}

// TestSelectionSortStability は、選択ソートが不安定ソートであることを確認します
// （同じ値の要素の相対的な順序が保持されない可能性がある）
func TestSelectionSortStability(t *testing.T) {
	// カスタム構造体を使用して安定性をテスト
	type Item struct {
		value int
		index int // 元のインデックス
	}
	
	items := []Item{
		{value: 3, index: 0},
		{value: 1, index: 1},
		{value: 3, index: 2},
		{value: 2, index: 3},
		{value: 3, index: 4},
	}
	
	// 値のみの配列を作成してソート
	values := make([]int, len(items))
	for i, item := range items {
		values[i] = item.value
	}
	
	SelectionSort(values)
	
	// 選択ソートは不安定なので、同じ値の要素の順序は保証されない
	expected := []int{1, 2, 3, 3, 3}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("SelectionSort stability test failed: got %v, want %v", values, expected)
	}
}

// BenchmarkSelectionSort は、SelectionSort関数のパフォーマンスをベンチマークします
func BenchmarkSelectionSort(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// テストデータを準備
			original := make([]int, size)
			for i := 0; i < size; i++ {
				original[i] = size - i // 逆順のデータ
			}
			
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				// 各イテレーションで新しいコピーを作成
				arr := make([]int, len(original))
				copy(arr, original)
				SelectionSort(arr)
			}
		})
	}
}

// BenchmarkSelectionSortBestCase は、最良の場合（すでにソート済み）のパフォーマンスをテストします
func BenchmarkSelectionSortBestCase(b *testing.B) {
	size := 1000
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		temp := make([]int, len(arr))
		copy(temp, arr)
		SelectionSort(temp)
	}
}

// BenchmarkSelectionSortWorstCase は、最悪の場合（逆順）のパフォーマンスをテストします
func BenchmarkSelectionSortWorstCase(b *testing.B) {
	size := 1000
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = size - i
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		temp := make([]int, len(arr))
		copy(temp, arr)
		SelectionSort(temp)
	}
}