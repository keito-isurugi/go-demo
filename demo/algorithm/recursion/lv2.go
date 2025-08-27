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