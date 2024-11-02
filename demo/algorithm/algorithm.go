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