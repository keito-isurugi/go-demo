package algorithm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 元の配列のコピーを作成（元配列が変更されないことを確認）
			original := make([]int, len(tt.input))
			copy(original, tt.input)
			
			result := QuickSort(tt.input)
			
			// 結果の検証
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("QuickSort() = %v, want %v", result, tt.expected)
			}
			
			// 元の配列が変更されていないことを確認
			if !reflect.DeepEqual(tt.input, original) {
				t.Errorf("元の配列が変更されました: %v, 期待値: %v", tt.input, original)
			}
		})
	}
}

// ベンチマークテスト
func BenchmarkQuickSort(b *testing.B) {
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
				QuickSort(arr)
			}
		})
	}
}

// パーティション関数のテスト
func TestPartition(t *testing.T) {
	tests := []struct {
		name           string
		arr            []int
		low            int
		high           int
		expectedPivot  int
		expectedArray  []int
	}{
		{
			name:           "基本的なケース",
			arr:            []int{3, 7, 8, 5, 2, 1, 9, 4},
			low:            0,
			high:           7,
			expectedPivot:  3,
			expectedArray:  []int{3, 2, 1, 4, 7, 8, 9, 5},
		},
		{
			name:           "部分配列",
			arr:            []int{1, 2, 9, 8, 7, 6, 5},
			low:            2,
			high:           6,
			expectedPivot:  2,
			expectedArray:  []int{1, 2, 5, 8, 7, 6, 9},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.arr))
			copy(arr, tt.arr)
			
			pivotIndex := partition(arr, tt.low, tt.high)
			
			if pivotIndex != tt.expectedPivot {
				t.Errorf("partition() ピボット位置 = %v, want %v", pivotIndex, tt.expectedPivot)
			}
			
			// ピボットより左の要素がすべてピボット以下であることを確認
			pivot := arr[pivotIndex]
			for i := tt.low; i < pivotIndex; i++ {
				if arr[i] > pivot {
					t.Errorf("ピボット左側の要素 arr[%d]=%d がピボット %d より大きい", i, arr[i], pivot)
				}
			}
			
			// ピボットより右の要素がすべてピボットより大きいことを確認
			for i := pivotIndex + 1; i <= tt.high; i++ {
				if arr[i] <= pivot {
					t.Errorf("ピボット右側の要素 arr[%d]=%d がピボット %d 以下", i, arr[i], pivot)
				}
			}
		})
	}
}