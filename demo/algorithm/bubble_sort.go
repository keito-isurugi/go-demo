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
