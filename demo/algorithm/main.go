package main

import (
	"fmt"

	"algorithm/recursion"
)

func main() {
	// arr := []int{5, 3, 8, 4, 2}
	// target := 8

	// ls := LinearSearch(arr, target)
	// fmt.Printf("Linear Search: %d\n", ls)

	// bas := BubbleAscSort(arr)
	// fmt.Printf("Bubble Asc Sort: %v\n", bas)

	// bds := BubbleDescSort(arr)
	// fmt.Printf("Bubble Desc Sort: %v\n", bds)
	
	// bs := BinarySearch(bas, 8)
	// fmt.Printf("Binary Search: %d\n", bs)

	// fmt.Println("\n--- Bit全探索のデモ ---")
	
	// 合計が特定の値になる組み合わせの数を求める
	// fmt.Println("\n2. 合計が特定の値になる組み合わせの数:")
	
	// 3種類の数値から合計4になる組み合わせ
	// nums1 := []int{1, 2, 3}
	// target1 := 4
	// count1 := CountSubsetSum(nums1, target1)
	// fmt.Println(count1, "通り")

	// fmt.Println("\n--- 挿入ソートのデモ ---")
	// RunInsertionSortDemo()
	
	// fmt.Println("\n--- クイックソートのデモ ---")
	// RunQuickSortDemo()
	
	// fmt.Println("\n--- マージソートのデモ ---")
	// RunMergeSortDemo()
	
	// fmt.Println("\n--- ヒープソートのデモ ---")
	// RunHeapSortDemo()
	
	// fmt.Println("\n--- フィボナッチ数列のデモ ---")
	// RunFibonacciDemo()
	
	// fmt.Println("\n--- 階乗計算のデモ ---")
	// FactorialDemo()

	fmt.Println("\n--- 再帰関数のデモ ---")
	// result := recursion.Factorial(5)
	// fmt.Println(result)

	// result2 := recursion.FactorialWithFor(5)
	// fmt.Println(result2)

	// result3 := recursion.CountDown(5)
	// fmt.Println(result3)

	// recursion.CountDownWithFor(5)

	// result4 := recursion.Sum(5)
	// fmt.Println(result4)

	// result5 := recursion.SumWithFor(5)
	// fmt.Println(result5)

	// arr := []int{1, 2, 3, 4, 5}
	// _ = recursion.PrintArray(arr, 0)
	// recursion.PrintArrayWithFor(arr)

	// fmt.Println(arr[:2])
	// arr2 := append(arr[:0], arr[3:]...)
	// fmt.Println(arr2)

	// 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144
	// result6 := recursion.Fibonacci(10)
	// fmt.Println(result6)

	// result7 := recursion.FibonacciWithFor(10)
	// fmt.Println(result7)

	// arr := []int{1, 2, 4, 10}
	// result8 := recursion.SumArray(arr, 0)
	// fmt.Println(result8)

	// result9 := recursion.SumArrayWithFor(arr)
	// fmt.Println(result9)

	// num := 5 / 2
	// fmt.Println(num)

	// arr2 := []int{1, 2, 3, 4, 5, 6}

	// for _, v := range arr2 {
	// 	result10 := recursion.BinarySearchWithFor(arr2, v)
	// 	fmt.Println(result10)
	// }

	// for _, v := range arr2 {
	// 	result10 := recursion.BinarySearchWithFor(arr2, v)
	// 	fmt.Println(result10)
	// }

	// recursion.Tree()
	
	// ハノイの塔のデモ
	// fmt.Println("\n--- ハノイの塔（再帰）のデモ ---")
	
	// // 3枚のディスクでシミュレーション
	// fmt.Println("■ 3枚のディスクでシミュレーション:")
	// recursion.SimulateTowerOfHanoi(3)
	
	// // ステップ表示版（再帰呼び出しの様子を表示）
	// fmt.Println("\n■ 再帰呼び出しの詳細（3枚）:")
	// recursion.TowerOfHanoiWithSteps(3, "A", "C", "B", 0)
	
	// // 移動回数の計算
	// fmt.Println("\n■ ディスク枚数と移動回数:")
	// for n := 1; n <= 10; n++ {
	// 	moves := recursion.CountHanoiMoves(n)
	// 	fmt.Printf("  %2d枚: %d回\n", n, moves)
	// }

	// recursion.TreeWithFor()
	
	// 配列の逆順（再帰）のデモ
	fmt.Println("\n--- 配列の逆順（再帰）のデモ ---")
	recursion.ReverseArrayDemo()
}