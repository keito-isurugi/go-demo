package algorithm

import (
	"reflect"
	"testing"
	"fmt"
)

func TestHeapSort(t *testing.T) {
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
			
			result := HeapSort(tt.input)
			
			// 結果の検証
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HeapSort() = %v, want %v", result, tt.expected)
			}
			
			// 元の配列が変更されていないことを確認
			if !reflect.DeepEqual(tt.input, original) {
				t.Errorf("元の配列が変更されました: %v, 期待値: %v", tt.input, original)
			}
		})
	}
}

func TestBuildMaxHeap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		checkHeap bool
	}{
		{
			name:      "通常のケース",
			input:     []int{4, 10, 3, 5, 1},
			checkHeap: true,
		},
		{
			name:      "既にヒープ",
			input:     []int{10, 5, 3, 4, 1},
			checkHeap: true,
		},
		{
			name:      "単一要素",
			input:     []int{42},
			checkHeap: true,
		},
		{
			name:      "空配列",
			input:     []int{},
			checkHeap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildMaxHeap(tt.input)
			
			if tt.checkHeap && !IsMaxHeap(result) {
				t.Errorf("BuildMaxHeap() = %v は最大ヒープではありません", result)
			}
		})
	}
}

func TestIsMaxHeap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name:     "最大ヒープ",
			input:    []int{10, 5, 3, 4, 1},
			expected: true,
		},
		{
			name:     "最大ヒープではない",
			input:    []int{1, 5, 3, 4, 10},
			expected: false,
		},
		{
			name:     "単一要素",
			input:    []int{42},
			expected: true,
		},
		{
			name:     "空配列",
			input:    []int{},
			expected: true,
		},
		{
			name:     "2要素（ヒープ）",
			input:    []int{5, 3},
			expected: true,
		},
		{
			name:     "2要素（非ヒープ）",
			input:    []int{3, 5},
			expected: false,
		},
		{
			name:     "完全な最大ヒープ",
			input:    []int{100, 19, 36, 17, 12, 25, 5, 9, 15, 6, 11, 13, 8, 1, 4},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMaxHeap(tt.input)
			if result != tt.expected {
				t.Errorf("IsMaxHeap(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHeapExtractMax(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		expectedMax int
		expectedLen int
	}{
		{
			name:        "通常のケース",
			input:       []int{10, 5, 3, 4, 1},
			expectedMax: 10,
			expectedLen: 4,
		},
		{
			name:        "単一要素",
			input:       []int{42},
			expectedMax: 42,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ヒープにする
			heap := BuildMaxHeap(tt.input)
			
			max, newHeap := HeapExtractMax(heap)
			
			if max != tt.expectedMax {
				t.Errorf("HeapExtractMax() max = %v, want %v", max, tt.expectedMax)
			}
			
			if len(newHeap) != tt.expectedLen {
				t.Errorf("HeapExtractMax() len = %v, want %v", len(newHeap), tt.expectedLen)
			}
			
			// 残りのヒープが正しいかチェック
			if len(newHeap) > 0 && !IsMaxHeap(newHeap) {
				t.Errorf("HeapExtractMax() 結果 %v は最大ヒープではありません", newHeap)
			}
		})
	}
}

func TestHeapInsert(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		value    int
		expected int // 期待される長さ
	}{
		{
			name:     "通常のケース",
			input:    []int{10, 5, 3, 4, 1},
			value:    15,
			expected: 6,
		},
		{
			name:     "空ヒープに挿入",
			input:    []int{},
			value:    42,
			expected: 1,
		},
		{
			name:     "最小値を挿入",
			input:    []int{10, 5, 3},
			value:    1,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ヒープにする
			heap := BuildMaxHeap(tt.input)
			
			result := HeapInsert(heap, tt.value)
			
			if len(result) != tt.expected {
				t.Errorf("HeapInsert() len = %v, want %v", len(result), tt.expected)
			}
			
			// 結果がヒープかチェック
			if !IsMaxHeap(result) {
				t.Errorf("HeapInsert() 結果 %v は最大ヒープではありません", result)
			}
			
			// 挿入した値が含まれているかチェック
			found := false
			for _, v := range result {
				if v == tt.value {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("HeapInsert() 挿入した値 %v が見つかりません", tt.value)
			}
		})
	}
}

func TestHeapify(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		heapSize  int
		rootIndex int
		checkHeap bool
	}{
		{
			name:      "ルートのヒープ化",
			input:     []int{1, 10, 3, 5, 4},
			heapSize:  5,
			rootIndex: 0,
			checkHeap: true,
		},
		{
			name:      "部分的なヒープ化",
			input:     []int{10, 1, 3, 5, 4},
			heapSize:  5,
			rootIndex: 1,
			checkHeap: false, // 部分的なので全体はヒープにならない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]int, len(tt.input))
			copy(arr, tt.input)
			
			heapify(arr, tt.heapSize, tt.rootIndex)
			
			// 指定された部分木がヒープ条件を満たしているかチェック
			if tt.checkHeap {
				// 簡単なチェック：ルートから始まる部分木のヒープ条件
				checkSubtreeHeap := func(arr []int, size, index int) bool {
					left := 2*index + 1
					right := 2*index + 2
					
					if left < size && arr[index] < arr[left] {
						return false
					}
					if right < size && arr[index] < arr[right] {
						return false
					}
					return true
				}
				
				if !checkSubtreeHeap(arr, tt.heapSize, tt.rootIndex) {
					t.Errorf("heapify() 後の部分木がヒープ条件を満たしていません: %v", arr)
				}
			}
		})
	}
}

// ベンチマークテスト
func BenchmarkHeapSort(b *testing.B) {
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
				HeapSort(arr)
			}
		})
	}
}

func BenchmarkBuildMaxHeap(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = size - i
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				BuildMaxHeap(arr)
			}
		})
	}
}