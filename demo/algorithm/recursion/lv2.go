package recursion

// import "fmt"

// n番目のフィボナッチ数を返す
// fib(5) = fib(4) + fib(3) = 3 + 2 = 5
// fib(4) = fib(3) + fib(2) = 2 + 1 = 3
// fib(3) = fib(2) + fib(1) = 1 + 1 = 2
// fib(2) = fib(1) + fib(0) = 1 + 0 = 1
// fib(1) = 1 = 1
// fib(0) = 0 = 0
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return Fibonacci(n - 1) + Fibonacci(n - 2)
}

func FibonacciWithFor(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a + b
	}
	return b
}

// 配列の合計
// [1, 2, 4, 10] = 17
// 0 = 1 + sum(arr, 1) = 1 + 16 = 17
// 1 = 2 + sum(arr, 2) = 2 + 14 = 16
// 2 = 4 + sum(arr, 3) = 4 + 10 = 14
// 3 = 10 + sum(arr, 4) = 10 + 0 = 10 
// 4 = 0 
func SumArray(arr []int, n int) int {
	if n >= len(arr){
		return 0
	}
	return arr[n] + SumArray(arr, n + 1)
}

func SumArrayWithFor(arr []int) int {
	if len(arr) <= 0 {
		return 0
	}
	
	sum := 0
	for i := 0; i < len(arr); i++ {
		sum = sum + arr[i]
	}
	return sum
}
