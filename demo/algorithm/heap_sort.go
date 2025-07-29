package main

import "fmt"

// HeapSort は配列を昇順にソート
func HeapSort(arr []int) []int {
	// 配列のコピーを作成（元の配列を変更しない）
	result := make([]int, len(arr))
	copy(result, arr)
	
	heapSortHelper(result)
	return result
}

// heapSortHelper は実際のヒープソート処理
func heapSortHelper(arr []int) {
	n := len(arr)
	
	// ヒープを構築（最大ヒープ）
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}
	
	// 一つずつ要素をヒープから取り出す
	for i := n - 1; i > 0; i-- {
		// 現在のルート（最大値）を末尾に移動
		arr[0], arr[i] = arr[i], arr[0]
		
		// 縮小されたヒープで再度ヒープ化
		heapify(arr, i, 0)
	}
}

// heapify は部分木をヒープ化する
func heapify(arr []int, heapSize, rootIndex int) {
	largest := rootIndex
	leftChild := 2*rootIndex + 1
	rightChild := 2*rootIndex + 2
	
	// 左の子が存在し、ルートより大きい場合
	if leftChild < heapSize && arr[leftChild] > arr[largest] {
		largest = leftChild
	}
	
	// 右の子が存在し、現在の最大値より大きい場合
	if rightChild < heapSize && arr[rightChild] > arr[largest] {
		largest = rightChild
	}
	
	// ルートが最大値でない場合
	if largest != rootIndex {
		arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
		
		// 影響を受けた部分木を再帰的にヒープ化
		heapify(arr, heapSize, largest)
	}
}

// HeapSortWithSteps はステップごとの状態を表示
func HeapSortWithSteps(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	fmt.Printf("初期: %v\n", result)
	
	n := len(result)
	
	// ヒープ構築フェーズ
	fmt.Println("\n--- ヒープ構築フェーズ ---")
	for i := n/2 - 1; i >= 0; i-- {
		fmt.Printf("ヒープ化: インデックス %d から\n", i)
		heapifyWithSteps(result, n, i)
		fmt.Printf("結果: %v\n\n", result)
	}
	
	fmt.Println("最大ヒープ完成:", result)
	
	// ソートフェーズ
	fmt.Println("\n--- ソートフェーズ ---")
	for i := n - 1; i > 0; i-- {
		fmt.Printf("ステップ %d: 最大値 %d を位置 %d に移動\n", n-i, result[0], i)
		
		// ルートを末尾に移動
		result[0], result[i] = result[i], result[0]
		fmt.Printf("交換後: %v\n", result)
		
		// 縮小されたヒープで再ヒープ化
		fmt.Printf("ヒープサイズ %d で再ヒープ化\n", i)
		heapifyWithSteps(result, i, 0)
		fmt.Printf("ヒープ化後: %v\n\n", result)
	}
	
	fmt.Printf("ソート完了: %v\n", result)
	return result
}

// heapifyWithSteps はヒープ化の過程を表示
func heapifyWithSteps(arr []int, heapSize, rootIndex int) {
	largest := rootIndex
	leftChild := 2*rootIndex + 1
	rightChild := 2*rootIndex + 2
	
	fmt.Printf("  ノード[%d]=%d の子を確認", rootIndex, arr[rootIndex])
	
	if leftChild < heapSize {
		fmt.Printf(" 左[%d]=%d", leftChild, arr[leftChild])
		if arr[leftChild] > arr[largest] {
			largest = leftChild
		}
	}
	
	if rightChild < heapSize {
		fmt.Printf(" 右[%d]=%d", rightChild, arr[rightChild])
		if arr[rightChild] > arr[largest] {
			largest = rightChild
		}
	}
	
	fmt.Printf("\n")
	
	if largest != rootIndex {
		fmt.Printf("  交換: [%d]=%d ↔ [%d]=%d\n", 
			rootIndex, arr[rootIndex], largest, arr[largest])
		arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
		
		// 再帰的にヒープ化
		heapifyWithSteps(arr, heapSize, largest)
	} else {
		fmt.Printf("  ヒープ条件満たしている\n")
	}
}

// BuildMaxHeap は配列を最大ヒープに変換
func BuildMaxHeap(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	n := len(result)
	for i := n/2 - 1; i >= 0; i-- {
		heapify(result, n, i)
	}
	
	return result
}

// IsMaxHeap は配列が最大ヒープかどうかを確認
func IsMaxHeap(arr []int) bool {
	n := len(arr)
	
	for i := 0; i < n/2; i++ {
		leftChild := 2*i + 1
		rightChild := 2*i + 2
		
		// 左の子が存在し、親より大きい
		if leftChild < n && arr[i] < arr[leftChild] {
			return false
		}
		
		// 右の子が存在し、親より大きい
		if rightChild < n && arr[i] < arr[rightChild] {
			return false
		}
	}
	
	return true
}

// HeapExtractMax は最大ヒープから最大値を取り出す
func HeapExtractMax(heap []int) (int, []int) {
	if len(heap) == 0 {
		return 0, heap
	}
	
	max := heap[0]
	
	// 最後の要素をルートに移動
	heap[0] = heap[len(heap)-1]
	heap = heap[:len(heap)-1]
	
	// ヒープ性質を回復
	if len(heap) > 0 {
		heapify(heap, len(heap), 0)
	}
	
	return max, heap
}

// HeapInsert は最大ヒープに新しい要素を挿入
func HeapInsert(heap []int, value int) []int {
	// 要素を末尾に追加
	heap = append(heap, value)
	
	// ボトムアップでヒープ性質を回復
	index := len(heap) - 1
	for index > 0 {
		parent := (index - 1) / 2
		
		if heap[parent] >= heap[index] {
			break
		}
		
		heap[parent], heap[index] = heap[index], heap[parent]
		index = parent
	}
	
	return heap
}

// RunHeapSortDemo はヒープソートのデモを実行
func RunHeapSortDemo() {
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
		sorted := HeapSort(sample)
		fmt.Printf("ソート後: %v\n", sorted)
	}
	
	// ヒープ操作のデモ
	fmt.Println("\n\n--- ヒープ操作のデモ ---")
	fmt.Println("配列をヒープに変換:")
	original := []int{4, 10, 3, 5, 1}
	fmt.Printf("元の配列: %v\n", original)
	
	heap := BuildMaxHeap(original)
	fmt.Printf("最大ヒープ: %v\n", heap)
	fmt.Printf("ヒープかどうか: %v\n", IsMaxHeap(heap))
	
	// 要素の挿入
	fmt.Println("\n要素の挿入:")
	heap = HeapInsert(heap, 15)
	fmt.Printf("15を挿入: %v\n", heap)
	
	// 最大値の取り出し
	fmt.Println("\n最大値の取り出し:")
	for len(heap) > 0 {
		max, newHeap := HeapExtractMax(heap)
		fmt.Printf("取り出し: %d, 残り: %v\n", max, newHeap)
		heap = newHeap
	}
	
	// ステップ表示のデモ
	fmt.Println("\n\n--- ヒープソートの詳細なステップ表示 ---")
	demoArr := []int{12, 11, 13, 5, 6, 7}
	fmt.Println("デモ配列でのソート過程:")
	HeapSortWithSteps(demoArr)
}