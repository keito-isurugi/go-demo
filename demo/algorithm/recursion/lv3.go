package recursion

import (
	"fmt"
	"os"
	"path/filepath"
)

// 二分探索
// [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
// 5
// 0, 9
// 5
// 0, 4
// 2
// 0, 1
// 1
// 1, 1
// [1, 2, 3, 4, 5, 6], 6
// 要素数 = 6, 目標 = 6
// ary[6 / 2] = 3
// 6 > 3

func BinarySearch(arr []int, target, left, right int) int {
	if left > right {
        return -1
    }
	
	mid := (left + right) / 2
	
	if target == arr[mid] {
		return mid
	}
	
	if target > arr[mid] {
		return BinarySearch(arr, target, mid + 1, right)
	} else {
		return BinarySearch(arr, target, left, mid - 1)
	}
}


func BinarySearchWithFor(arr []int, target int) int {
	if len(arr) <= 0 {
		fmt.Println("array is empty")
		return -1
	}
	
	left, right := 0, len(arr) -1
	
	for left <= right {
		mid := (left + right) / 2
		
		if target == arr[mid] {
			return mid
		}
		
		if target > arr[mid] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	fmt.Println("target not found")
	return -1
}

func Tree() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Current Directory:", currentDir)

	dirName := filepath.Base(currentDir)
	fmt.Println("Directory Name:", dirName)

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Entries:", entries)
	
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("├──", entry.Name())
		} else {
			fmt.Println("├──", entry.Name())
		}
	}
}

func TreeWithFor() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Current Directory:", currentDir)
	
	dirName := filepath.Base(currentDir)
	fmt.Println("Directory Name:", dirName)
	
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Entries:", entries)
	
	entriesNum := len(entries) -1
	fmt.Println(".")
	for i, entry := range entries {
		if entry.IsDir() {
			if i == entriesNum {
				fmt.Println("└──", entry.Name())
			} else {
				fmt.Println("├──", entry.Name())
			}
		} else {
			if i == entriesNum {
				fmt.Println("└──", entry.Name())
			} else {
				fmt.Println("├──", entry.Name())
			}
		}
	}
}