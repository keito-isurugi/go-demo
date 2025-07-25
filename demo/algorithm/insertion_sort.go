package main

import "fmt"

// InsertionSort は挿入ソートアルゴリズムを実装します
// 配列を昇順にソートします
func InsertionSort(arr []int) []int {
	// 配列のコピーを作成（元の配列を変更しない）
	result := make([]int, len(arr))
	copy(result, arr)
	
	// インデックス1から開始（0番目は既にソート済みとみなす）
	for i := 1; i < len(result); i++ {
		// 現在の要素を保存
		key := result[i]
		j := i - 1
		
		// keyより大きい要素を右にシフト
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}
		
		// 適切な位置にkeyを挿入
		result[j+1] = key
	}
	
	return result
}

// InsertionSortWithSteps は挿入ソートの各ステップを表示します
// デバッグや学習用途に使用できます
func InsertionSortWithSteps(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	fmt.Printf("初期配列: %v\n", result)
	
	for i := 1; i < len(result); i++ {
		key := result[i]
		j := i - 1
		
		fmt.Printf("\nステップ %d: 要素 %d を挿入\n", i, key)
		
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}
		
		result[j+1] = key
		fmt.Printf("配列の状態: %v\n", result)
	}
	
	return result
}

// RunInsertionSortDemo は挿入ソートのデモを実行します
func RunInsertionSortDemo() {
	samples := [][]int{
		{64, 34, 25, 12, 22, 11, 90},
		{5, 2, 4, 6, 1, 3},
		{1},
		{3, 3, 3, 3},
		{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	}
	
	for i, sample := range samples {
		fmt.Printf("\n例 %d:\n", i+1)
		fmt.Printf("ソート前: %v\n", sample)
		sorted := InsertionSort(sample)
		fmt.Printf("ソート後: %v\n", sorted)
	}
}