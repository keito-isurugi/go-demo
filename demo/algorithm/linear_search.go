package main

func LinearSearch(array []int, target int) int {
	for i, v := range array {
		if v == target {
			return i
		}
	}
	return -1
}