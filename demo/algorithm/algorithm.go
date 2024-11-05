package algorithm

import "fmt"

func LinearSearch(slice []int, t int) (int, error) {
	for _, v := range slice {
		if v == t {
			return v, nil
 		}
	}
	return 0, fmt.Errorf("該当する値が存在しません")
}

func BubbleAscSort(slice []int) []int {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if slice[j] > slice[j+1] {
				tmp := slice[j]
				slice[j] = slice[j+1]
				slice[j+1] = tmp
			}
		} 
	}
	return slice
}

func BubbleDescSort(slice []int) []int {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if slice[j] < slice[j+1] {
				tmp := slice[j]
				slice[j] = slice[j+1]
				slice[j+1] = tmp
			}
		} 
	}
	return slice
}