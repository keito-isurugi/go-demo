package recursion

import "fmt"

// 二分探索
// [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
// 5
// 0, 9
// 5
// 0, 4
// 2
// 0, 1
// 1
// 1, 1
// [1, 2, 3, 4, 5, 6], 6
// 要素数 = 6, 目標 = 6
// ary[6 / 2] = 3
// 6 > 3

// func BinarySearch(arr []int, target int) int {
// }

func BinarySearchWithFor(arr []int, target int) int {
	if len(arr) <= 0 {
		fmt.Println("array is empty")
		return -1
	}
	
	left, right := 0, len(arr) -1
	
	for left <= right {
		mid := (left + right) / 2
		
		if target == arr[mid] {
			return mid
		}
		
		if target > arr[mid] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	fmt.Println("target not found")
	return -1
}