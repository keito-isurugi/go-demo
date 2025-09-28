package main

import (
	"fmt"
	"sort"
)

func SortDemo() {
	// sort.Slice
	arr1 := []int{5, 3, 8, 4, 2}
	fmt.Println(arr1)
	// 昇順
	sort.Slice(arr1, func(i, j int) bool {
		return arr1[i] < arr1[j]
	})
	fmt.Println(arr1)
	// 降順
	sort.Slice(arr1, func(i, j int) bool {
		return arr1[i] > arr1[j]
	})
	fmt.Println(arr1)

	// sort.Ints
	// 昇順
	sort.Ints(arr1)
	fmt.Println(arr1)
}