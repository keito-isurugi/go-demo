package algorithm

import "fmt"

// QuickSort はクイックソートアルゴリズムを実装します
// 配列を昇順にソートします
func QuickSort(arr []int) []int {
	// 配列のコピーを作成（元の配列を変更しない）
	result := make([]int, len(arr))
	copy(result, arr)
	
	// 実際のソート処理を実行
	quickSortHelper(result, 0, len(result)-1)
	
	return result
}

// quickSortHelper は再帰的にクイックソートを実行します
func quickSortHelper(arr []int, low, high int) {
	if low < high {
		// パーティション操作を実行し、ピボットの位置を取得
		pi := partition(arr, low, high)
		
		// ピボットの左側と右側をそれぞれソート
		quickSortHelper(arr, low, pi-1)
		quickSortHelper(arr, pi+1, high)
	}
}

// partition は配列を分割し、ピボットの最終位置を返します
func partition(arr []int, low, high int) int {
	// 最後の要素をピボットとして選択
	pivot := arr[high]
	
	// より小さい要素の位置を追跡
	i := low - 1
	
	for j := low; j < high; j++ {
		// 現在の要素がピボット以下の場合
		if arr[j] <= pivot {
			i++
			// 要素を交換
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	
	// ピボットを正しい位置に配置
	arr[i+1], arr[high] = arr[high], arr[i+1]
	
	return i + 1
}

// QuickSortWithSteps はクイックソートの各ステップを表示します
// デバッグや学習用途に使用できます
func QuickSortWithSteps(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	fmt.Printf("初期配列: %v\n", result)
	quickSortHelperWithSteps(result, 0, len(result)-1, 0)
	
	return result
}

// quickSortHelperWithSteps は各ステップを表示しながらソートします
func quickSortHelperWithSteps(arr []int, low, high, depth int) {
	if low < high {
		indent := getIndent(depth)
		fmt.Printf("%s範囲 [%d:%d] をソート中\n", indent, low, high)
		
		pi := partitionWithSteps(arr, low, high, depth)
		
		fmt.Printf("%sピボット位置: %d, 値: %d\n", indent, pi, arr[pi])
		fmt.Printf("%s配列の状態: %v\n\n", indent, arr)
		
		quickSortHelperWithSteps(arr, low, pi-1, depth+1)
		quickSortHelperWithSteps(arr, pi+1, high, depth+1)
	}
}

// partitionWithSteps は分割プロセスを表示します
func partitionWithSteps(arr []int, low, high, depth int) int {
	indent := getIndent(depth)
	pivot := arr[high]
	fmt.Printf("%sピボット: %d\n", indent, pivot)
	
	i := low - 1
	
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			if i != j {
				fmt.Printf("%s交換: arr[%d]=%d <-> arr[%d]=%d\n", indent, i, arr[i], j, arr[j])
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	
	arr[i+1], arr[high] = arr[high], arr[i+1]
	fmt.Printf("%sピボットを最終位置に配置: arr[%d]=%d\n", indent, i+1, arr[i+1])
	
	return i + 1
}

// getIndent は深さに応じたインデントを返します
func getIndent(depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	return indent
}

// RunQuickSortDemo はクイックソートのデモを実行します
func RunQuickSortDemo() {
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
		sorted := QuickSort(sample)
		fmt.Printf("ソート後: %v\n", sorted)
	}
	
	// ステップ表示のデモ
	fmt.Println("\n\n--- クイックソートの詳細なステップ表示 ---")
	demoArr := []int{3, 7, 8, 5, 2, 1, 9, 4}
	fmt.Println("デモ配列でのソート過程:")
	QuickSortWithSteps(demoArr)
}