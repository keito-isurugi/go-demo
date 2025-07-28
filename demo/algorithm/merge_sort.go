package main

import "fmt"

// MergeSort は配列を昇順にソート
func MergeSort(arr []int) []int {
	// 配列のコピーを作成（元の配列を変更しない）
	result := make([]int, len(arr))
	copy(result, arr)
	
	// 実際のソート処理を実行
	mergeSortHelper(result, 0, len(result)-1)
	
	return result
}

// mergeSortHelper は再帰的にソート
func mergeSortHelper(arr []int, left, right int) {
	if left < right {
		// 中央を計算
		mid := left + (right-left)/2
		
		// 左半分をソート
		mergeSortHelper(arr, left, mid)
		
		// 右半分をソート
		mergeSortHelper(arr, mid+1, right)
		
		// ソート済みの配列をマージ
		merge(arr, left, mid, right)
	}
}

// merge は2つのソート済み配列をマージ
func merge(arr []int, left, mid, right int) {
	// 左右の配列のサイズを計算
	n1 := mid - left + 1
	n2 := right - mid
	
	// 一時配列を作成
	leftArr := make([]int, n1)
	rightArr := make([]int, n2)
	
	// データをコピー
	for i := 0; i < n1; i++ {
		leftArr[i] = arr[left+i]
	}
	for j := 0; j < n2; j++ {
		rightArr[j] = arr[mid+1+j]
	}
	
	// マージ処理
	i, j, k := 0, 0, left
	
	for i < n1 && j < n2 {
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}
	
	// 残りの要素をコピー
	for i < n1 {
		arr[k] = leftArr[i]
		i++
		k++
	}
	
	for j < n2 {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

// MergeSortWithSteps はステップごとの状態を表示
func MergeSortWithSteps(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	fmt.Printf("初期: %v\n", result)
	mergeSortHelperWithSteps(result, 0, len(result)-1, 0)
	
	return result
}

// mergeSortHelperWithSteps は各ステップを表示しながらソート
func mergeSortHelperWithSteps(arr []int, left, right, depth int) {
	if left < right {
		indent := getMergeIndent(depth)
		mid := left + (right-left)/2
		
		fmt.Printf("%s分割: [%d:%d] → [%d:%d] + [%d:%d]\n", 
			indent, left, right, left, mid, mid+1, right)
		
		// 左半分をソート
		mergeSortHelperWithSteps(arr, left, mid, depth+1)
		
		// 右半分をソート
		mergeSortHelperWithSteps(arr, mid+1, right, depth+1)
		
		// マージ
		mergeWithSteps(arr, left, mid, right, depth)
	}
}

// mergeWithSteps はマージ過程を表示
func mergeWithSteps(arr []int, left, mid, right, depth int) {
	indent := getMergeIndent(depth)
	
	// マージ前の状態を表示
	fmt.Printf("%sマージ前: ", indent)
	for i := left; i <= right; i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Printf("\n")
	
	// 通常のマージ処理
	merge(arr, left, mid, right)
	
	// マージ後の状態を表示
	fmt.Printf("%sマージ後: ", indent)
	for i := left; i <= right; i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Printf("\n\n")
}

// getMergeIndent は深さに応じたインデントを返す
func getMergeIndent(depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	return indent
}

// RunMergeSortDemo はマージソートのデモを実行
func RunMergeSortDemo() {
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
		sorted := MergeSort(sample)
		fmt.Printf("ソート後: %v\n", sorted)
	}
	
	// ステップ表示のデモ
	fmt.Println("\n\n--- マージソートの詳細なステップ表示 ---")
	demoArr := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Println("デモ配列でのソート過程:")
	MergeSortWithSteps(demoArr)
}

// ボトムアップ版マージソート（非再帰）
func MergeSortBottomUp(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)
	
	// サブ配列のサイズを2倍ずつ増やしていく
	for size := 1; size < n; size *= 2 {
		// 現在のサイズで配列全体を処理
		for start := 0; start < n-size; start += 2*size {
			// マージする範囲を計算
			mid := start + size - 1
			end := min(start+2*size-1, n-1)
			
			// マージ実行
			merge(result, start, mid, end)
		}
	}
	
	return result
}

// min は2つの整数の小さい方を返す
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}