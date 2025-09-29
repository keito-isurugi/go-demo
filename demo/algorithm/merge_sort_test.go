package algorithm

import (
	"reflect"
	"testing"
	"fmt"
)

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "通常のケース",
			input:    []int{64, 34, 25, 12, 22, 11, 90},
			expected: []int{11, 12, 22, 25, 34, 64, 90},
		},
		{
			name:     "既にソート済み",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "逆順",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "重複要素あり",
			input:    []int{3, 5, 3, 5, 3},
			expected: []int{3, 3, 3, 5, 5},
		},
		{
			name:     "単一要素",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "空配列",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "負の数を含む",
			input:    []int{-3, 0, -1, 5, -2, 4},
			expected: []int{-3, -2, -1, 0, 4, 5},
		},
		{
			name:     "大きな配列",
			input:    []int{100, 23, 45, 67, 89, 12, 34, 56, 78, 90, 11, 22, 33, 44, 55},
			expected: []int{11, 12, 22, 23, 33, 34, 44, 45, 55, 56, 67, 78, 89, 90, 100},
		},
		{
			name:     "2要素",
			input:    []int{2, 1},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の配列のコピーを作成（元配列が変更されないことを確認）
			original := make([]int, len(tt.input))
			copy(original, tt.input)
			
			result := MergeSort(tt.input)
			
			// 結果の検証
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeSort() = %v, want %v", result, tt.expected)
			}
			
			// 元の配列が変更されていないことを確認
			if !reflect.DeepEqual(tt.input, original) {
				t.Errorf("元の配列が変更されました: %v, 期待値: %v", tt.input, original)
			}
		})
	}
}

func TestMergeSortBottomUp(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "通常のケース",
			input:    []int{38, 27, 43, 3, 9, 82, 10},
			expected: []int{3, 9, 10, 27, 38, 43, 82},
		},
		{
			name:     "2のべき乗サイズ",
			input:    []int{8, 7, 6, 5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "奇数サイズ",
			input:    []int{5, 3, 1, 4, 2},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]int, len(tt.input))
			copy(original, tt.input)
			
			result := MergeSortBottomUp(tt.input)
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeSortBottomUp() = %v, want %v", result, tt.expected)
			}
			
			if !reflect.DeepEqual(tt.input, original) {
				t.Errorf("元の配列が変更されました: %v, 期待値: %v", tt.input, original)
			}
		})
	}
}

// ベンチマークテスト
func BenchmarkMergeSort(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// テスト用の配列を生成
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = size - i // 逆順の配列
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MergeSort(arr)
			}
		})
	}
}

func BenchmarkMergeSortBottomUp(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = size - i
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MergeSortBottomUp(arr)
			}
		})
	}
}

// マージ関数のテスト
func TestMerge(t *testing.T) {
	tests := []struct {
		name          string
		arr           []int
		left          int
		mid           int
		right         int
		expectedArray []int
	}{
		{
			name:          "基本的なケース",
			arr:           []int{38, 27, 43, 3, 9, 82, 10},
			left:          0,
			mid:           2,
			right:         6,
			expectedArray: []int{3, 9, 10, 27, 38, 43, 82},
		},
		{
			name:          "部分配列のマージ",
			arr:           []int{1, 5, 8, 2, 4, 6, 9},
			left:          1,
			mid:           2,
			right:         5,
			expectedArray: []int{1, 2, 4, 5, 6, 8, 9},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			
			// まず左右をソート
			mergeSortHelper(arr, tt.left, tt.mid)
			mergeSortHelper(arr, tt.mid+1, tt.right)
			
			// マージを実行
			merge(arr, tt.left, tt.mid, tt.right)
			
			if !reflect.DeepEqual(arr, tt.expectedArray) {
				t.Errorf("merge() 結果 = %v, want %v", arr, tt.expectedArray)
			}
		})
	}
}