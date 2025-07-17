package main

import "fmt"

// SelectionSort は、与えられた整数スライスを選択ソートアルゴリズムで昇順にソートします。
// 選択ソートは、配列の中から最小値を見つけて先頭と交換することを繰り返すアルゴリズムです。
func SelectionSort(arr []int) {
	n := len(arr)
	
	// 配列の各要素について処理を行う
	for i := 0; i < n-1; i++ {
		// i番目以降の要素の中で最小値のインデックスを見つける
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		
		// 最小値をi番目の要素と交換する
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
}

// SelectionSortDescending は、与えられた整数スライスを選択ソートアルゴリズムで降順にソートします。
func SelectionSortDescending(arr []int) {
	n := len(arr)
	
	// 配列の各要素について処理を行う
	for i := 0; i < n-1; i++ {
		// i番目以降の要素の中で最大値のインデックスを見つける
		maxIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] > arr[maxIndex] {
				maxIndex = j
			}
		}
		
		// 最大値をi番目の要素と交換する
		if maxIndex != i {
			arr[i], arr[maxIndex] = arr[maxIndex], arr[i]
		}
	}
}

// SelectionSortWithSteps は、選択ソートの各ステップを表示しながらソートを実行します。
// デバッグや学習目的で使用できます。
func SelectionSortWithSteps(arr []int) {
	n := len(arr)
	fmt.Printf("初期状態: %v\n", arr)
	
	for i := 0; i < n-1; i++ {
		minIndex := i
		fmt.Printf("\nステップ %d: %d番目の要素の位置を決定\n", i+1, i)
		
		// 最小値を探す
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		
		// 交換前の状態を表示
		fmt.Printf("  最小値 %d (インデックス: %d) を位置 %d と交換\n", arr[minIndex], minIndex, i)
		
		// 最小値を交換
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
			fmt.Printf("  交換後: %v\n", arr)
		} else {
			fmt.Printf("  交換不要: %v\n", arr)
		}
	}
	
	fmt.Printf("\nソート完了: %v\n", arr)
}

