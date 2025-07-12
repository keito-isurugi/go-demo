package main

import (
	"fmt"
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

	fmt.Println("\n--- Bit全探索のデモ ---")
	
	// 例1: 集合 {"A", "B", "C"} のすべての部分集合を表示
	// fmt.Println("\n1. すべての部分集合:")
	// items := []string{"A", "B", "C"}
	// patterns := BitFullSearch(len(items))
	// for _, bit := range patterns {
	// 	PrintSubset(items, bit)
	// }
	
	// 例2: 合計が特定の値になる組み合わせの数を求める
	fmt.Println("\n2. 合計が特定の値になる組み合わせの数:")
	
	// ケース1: 3種類の数値から合計4になる組み合わせ
	nums1 := []int{1, 2, 3}
	target1 := 4
	count1 := CountSubsetSum(nums1, target1)
	fmt.Println(count1, "通り")
	// fmt.Printf("%d種類の数値 %v の内、合計が%dになる組み合わせの数: %d通り\n", 
	// 	len(nums1), nums1, target1, count1)
	
	// // ケース2: 6種類の数値から合計20になる組み合わせ
	// nums2 := []int{2, 4, 6, 8, 10, 12}
	// target2 := 20
	// count2 := CountSubsetSum(nums2, target2)
	// fmt.Printf("%d種類の数値 %v の内、合計が%dになる組み合わせの数: %d通り\n", 
	// 	len(nums2), nums2, target2, count2)
	
	// // 例3: 詳細表示（組み合わせも表示）
	// fmt.Println("\n3. 組み合わせの詳細表示:")
	// nums3 := []int{1, 2, 3, 4, 5}
	// target3 := 10
	// PrintSubsetSumDetail(nums3, target3)
	
	// // 例2: 部分和問題
	// fmt.Println("\n2. 部分和問題 (合計が10になる組み合わせ):")
	// nums := []int{2, 3, 5, 7}
	// targetSum := 10
	// solutions := SubsetSum(nums, targetSum)
	// for _, bit := range solutions {
	// 	fmt.Printf("解: ")
	// 	first := true
	// 	for i := 0; i < len(nums); i++ {
	// 		if (bit >> i) & 1 == 1 {
	// 			if !first {
	// 				fmt.Print(" + ")
	// 			}
	// 			fmt.Print(nums[i])
	// 			first = false
	// 		}
	// 	}
	// 	fmt.Printf(" = %d\n", targetSum)
	// }
	
	// // 例3: 偶数個の要素を持つ部分集合の数を数える
	// fmt.Println("\n3. 偶数個の要素を持つ部分集合の数:")
	// n := 4
	// evenCount := CountSubsets(n, func(bit int) bool {
	// 	count := 0
	// 	for i := 0; i < n; i++ {
	// 		if (bit >> i) & 1 == 1 {
	// 			count++
	// 		}
	// 	}
	// 	return count%2 == 0
	// })
	// fmt.Printf("4要素の集合で、偶数個の要素を持つ部分集合: %d個\n", evenCount)
}
