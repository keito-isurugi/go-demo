package main

import (
	"fmt"
)

func main() {
	arr := []int{5, 3, 8, 4, 2}
	target := 8

	ls := LinearSearch(arr, target)
	fmt.Printf("Linear Search: %d\n", ls)

	bas := BubbleAscSort(arr)
	fmt.Printf("Bubble Asc Sort: %v\n", bas)

	bds := BubbleDescSort(arr)
	fmt.Printf("Bubble Desc Sort: %v\n", bds)
	
	bs := BinarySearch(bas, 8)
	fmt.Printf("Binary Search: %d\n", bs)
}
