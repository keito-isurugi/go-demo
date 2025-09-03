package recursion

import "fmt"

// ReverseArray は配列を再帰的に逆順にする
func ReverseArray(arr []int, start, end int) []int {
	// ベースケース: 配列の中央に到達したら終了
	if start >= end {
		return arr
	}
	
	// 要素を交換
	arr[start], arr[end] = arr[end], arr[start]
	
	// 再帰呼び出し: 範囲を狭めて継続
	return ReverseArray(arr, start+1, end-1)
}

// ReverseArrayWithSlice はスライス操作で配列を逆順にする（再帰）
func ReverseArrayWithSlice(arr []int) []int {
	// ベースケース: 要素が1つ以下なら逆順にする必要なし
	if len(arr) <= 1 {
		return arr
	}
	
	// 最初の要素を取り出し、残りを再帰的に逆順にしてから結合
	return append(ReverseArrayWithSlice(arr[1:]), arr[0])
}

// ReverseArrayWithFor は反復処理で配列を逆順にする（比較用）
func ReverseArrayWithFor(arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	copy(result, arr)
	
	for i := 0; i < n/2; i++ {
		result[i], result[n-1-i] = result[n-1-i], result[i]
	}
	
	return result
}

// ReverseArrayWithSteps は再帰呼び出しの過程を表示
func ReverseArrayWithSteps(arr []int, start, end, depth int) []int {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	fmt.Printf("%sReverseArray([%v], start=%d, end=%d)\n", indent, arr, start, end)
	
	if start >= end {
		fmt.Printf("%s  ベースケース到達: 終了\n", indent)
		return arr
	}
	
	fmt.Printf("%s  交換: arr[%d]=%d <-> arr[%d]=%d\n", 
		indent, start, arr[start], end, arr[end])
	arr[start], arr[end] = arr[end], arr[start]
	fmt.Printf("%s  交換後: %v\n", indent, arr)
	
	return ReverseArrayWithSteps(arr, start+1, end-1, depth+1)
}

// ReverseArrayDemo は配列逆順のデモを実行
func ReverseArrayDemo() {
	fmt.Println("■ 配列の逆順（再帰）デモ")
	fmt.Println("==================================================")
	
	// 基本的な例
	fmt.Println("\n1. 基本的な配列の逆順:")
	arr1 := []int{1, 2, 3, 4, 5}
	fmt.Printf("元の配列: %v\n", arr1)
	
	// コピーを作成して逆順にする
	reversed1 := make([]int, len(arr1))
	copy(reversed1, arr1)
	reversed1 = ReverseArray(reversed1, 0, len(reversed1)-1)
	fmt.Printf("逆順後:   %v\n", reversed1)
	
	// スライス操作版
	fmt.Println("\n2. スライス操作での逆順:")
	arr2 := []int{10, 20, 30, 40, 50, 60}
	fmt.Printf("元の配列: %v\n", arr2)
	reversed2 := ReverseArrayWithSlice(arr2)
	fmt.Printf("逆順後:   %v\n", reversed2)
	
	// ステップ表示版
	fmt.Println("\n3. 再帰呼び出しの詳細:")
	arr3 := []int{1, 2, 3, 4, 5}
	fmt.Printf("初期配列: %v\n", arr3)
	fmt.Println("再帰呼び出しの過程:")
	reversed3 := make([]int, len(arr3))
	copy(reversed3, arr3)
	reversed3 = ReverseArrayWithSteps(reversed3, 0, len(reversed3)-1, 0)
	fmt.Printf("最終結果: %v\n", reversed3)
	
	// パフォーマンス比較
	fmt.Println("\n4. パフォーマンス比較:")
	testSizes := []int{10, 100, 1000}
	
	for _, size := range testSizes {
		// テストデータ生成
		testArr := make([]int, size)
		for i := 0; i < size; i++ {
			testArr[i] = i + 1
		}
		
		fmt.Printf("\n配列サイズ %d の場合:\n", size)
		
		// 再帰版（インプレース）
		arr := make([]int, size)
		copy(arr, testArr)
		_ = ReverseArray(arr, 0, len(arr)-1)
		fmt.Printf("  再帰版（インプレース）: 完了\n")
		
		// 再帰版（スライス）
		_ = ReverseArrayWithSlice(testArr)
		fmt.Printf("  再帰版（スライス）: 完了\n")
		
		// 反復版
		_ = ReverseArrayWithFor(testArr)
		fmt.Printf("  反復版: 完了\n")
	}
	
	// 特殊ケース
	fmt.Println("\n5. 特殊ケースのテスト:")
	
	// 空配列
	empty := []int{}
	fmt.Printf("空配列: %v -> %v\n", empty, ReverseArrayWithSlice(empty))
	
	// 要素1つ
	single := []int{42}
	fmt.Printf("要素1つ: %v -> %v\n", single, ReverseArrayWithSlice(single))
	
	// 要素2つ
	two := []int{1, 2}
	twoReversed := make([]int, 2)
	copy(twoReversed, two)
	twoReversed = ReverseArray(twoReversed, 0, 1)
	fmt.Printf("要素2つ: %v -> %v\n", two, twoReversed)
}