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
	for i := range slice {
		if slice[i] > slice[i - 1] {
			tmp := slice[i]
			slice[i] = slice[i - 1]
			slice[i - 1] = tmp
		}
	}
	return slice
}