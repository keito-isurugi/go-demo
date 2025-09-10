package main

import "fmt"

func ConditionalBranch() {
	num1 := 10
	if num1 > 1 {
		fmt.Println("over 1")
	} else {
		fmt.Println("under 0")
	}

	switch num1 {
	case 1:
		fmt.Println("num1 is 1")
	case 10:
		fmt.Println("num1 is 10")
	default:
		fmt.Println("default")
	}

	switch num1 {
	case 1:
		fmt.Println("num1 is 1")
	case 10, 20:
		fmt.Println("num1 is 10 or 20")
	default:
		fmt.Println("default")
	}
}