package main

func BubbleAscSort(array []int) []int {
	sortedArray := append([]int(nil), array...)

	n := len(sortedArray)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if sortedArray[j] > sortedArray[j+1] {
				tmp := sortedArray[j]
				sortedArray[j] = sortedArray[j+1]
				sortedArray[j+1] = tmp
			}
		} 
	}

	return sortedArray
}

func BubbleDescSort(array []int) []int {
	sortedArray := append([]int(nil), array...)

	n := len(sortedArray)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if sortedArray[j] < sortedArray[j+1] {
				tmp := sortedArray[j]
				sortedArray[j] = sortedArray[j+1]
				sortedArray[j+1] = tmp
			}
		} 
	}

	return sortedArray
}
