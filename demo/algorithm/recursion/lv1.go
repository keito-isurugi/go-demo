package recursion

import "fmt"

// 階乗
func Factorial(n int) int {
	if n <= 0 {
		return 1
	}
	return n * Factorial(n - 1)
}

func FactorialWithFor(n int) int {
	if n <= 0 {
		return 1
	}
	result := 1
	for i := n; i > 0; i-- {
		result = result * i
	}
	return result
}

// nから0までカウントダウン
func CountDown(n int) int {
	if n <= 0 {
		return 0
	}
	count := n
	fmt.Println(count)
	return CountDown(count - 1)
}

func CountDownWithFor(n int) {
	if n <= 0 {
		fmt.Println(0)
	}
	for i := n; i >= 0; i-- {
		fmt.Println(i)	
	}
}

// 1からnまでの合計(負の値は0とする)
func Sum(n int) int {
	if n <= 0 {
		return 0
	}
	return n + Sum(n - 1)
}

func SumWithFor(n int) int {
	result := 0
	for i := n; i >= 0; i-- {
		result = result + i
	}
	return result
}

// 配列の要素を1つずつ表示
func PrintArray(arr []int, n int) int {
	if n >= len(arr) {
		return 0
	}
	fmt.Println(arr[n])
	return PrintArray(arr, n + 1)
}

func PrintArrayWithFor(arr []int) {
	if len(arr) <= 0 {
		fmt.Println("array is empty")
	}
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])	
	}
}