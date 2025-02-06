package main



func BubbleAscSort(array []int) []int {
	n := len(array)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if array[j] > array[j+1] {
				tmp := array[j]
				array[j] = array[j+1]
				array[j+1] = tmp
			}
		} 
	}
	return array
}

func BubbleDescSort(array []int) []int {
	n := len(array)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if array[j] < array[j+1] {
				tmp := array[j]
				array[j] = array[j+1]
				array[j+1] = tmp
			}
		} 
	}
	return array
}

func BinarySearch(array []int, target int) int {
	left, right := 0, len(array) - 1

	for left <= right {
		mid := left + (right - left) / 2
		
		if array[mid] == target {
			return mid
		} else if array[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}