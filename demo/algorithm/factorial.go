package algorithm

import (
	"fmt"
	"math/big"
)

func factorial(n int) uint64 {
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	if n == 0 || n == 1 {
		return 1
	}
	return uint64(n) * factorial(n-1)
}

func factorialWithMemo(n int, memo map[int]uint64) uint64 {
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	if n == 0 || n == 1 {
		return 1
	}
	if val, exists := memo[n]; exists {
		return val
	}
	result := uint64(n) * factorialWithMemo(n-1, memo)
	memo[n] = result
	return result
}

func factorialBigInt(n int) *big.Int {
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	if n == 0 || n == 1 {
		return big.NewInt(1)
	}
	result := big.NewInt(int64(n))
	return result.Mul(result, factorialBigInt(n-1))
}

func FactorialDemo() {
	fmt.Println("===== Factorial Demo =====")
	
	fmt.Println("\n基本的な再帰による階乗計算:")
	for i := 0; i <= 10; i++ {
		result := factorial(i)
		fmt.Printf("%d! = %d\n", i, result)
	}
	
	fmt.Println("\nメモ化を使った階乗計算:")
	memo := make(map[int]uint64)
	for i := 0; i <= 20; i++ {
		result := factorialWithMemo(i, memo)
		fmt.Printf("%d! = %d\n", i, result)
	}
	
	fmt.Println("\n大きな数の階乗計算 (big.Int使用):")
	for _, n := range []int{25, 30, 50, 100} {
		result := factorialBigInt(n)
		fmt.Printf("%d! = %s\n", n, result.String())
	}
	
	fmt.Println("\n計算過程の可視化 (5!):")
	visualizeFactorial(5, 0)
}

func visualizeFactorial(n int, depth int) uint64 {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	if n == 0 || n == 1 {
		fmt.Printf("%sfactorial(%d) = 1\n", indent, n)
		return 1
	}
	
	fmt.Printf("%sfactorial(%d) = %d * factorial(%d)\n", indent, n, n, n-1)
	result := uint64(n) * visualizeFactorial(n-1, depth+1)
	fmt.Printf("%sfactorial(%d) = %d\n", indent, n, result)
	return result
}