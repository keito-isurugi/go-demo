package main

import "fmt"

func main() {
	arr := []int{5, 3, 8, 4, 2}
	target := 8

	ls := LinearSearch(arr, target)
	fmt.Printf("Linear Search: %d\n", ls)
}
