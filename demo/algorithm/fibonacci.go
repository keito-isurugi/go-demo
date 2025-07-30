package main

import (
	"fmt"
	"time"
)

// FibonacciRecursive は再帰でフィボナッチ数列のn番目を計算
func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

// FibonacciRecursiveWithSteps は再帰の呼び出し過程を表示
func FibonacciRecursiveWithSteps(n int, depth int) int {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	fmt.Printf("%sFib(%d) を計算開始\n", indent, n)
	
	if n <= 1 {
		fmt.Printf("%sFib(%d) = %d (ベースケース)\n", indent, n, n)
		return n
	}
	
	fmt.Printf("%sFib(%d) = Fib(%d) + Fib(%d)\n", indent, n, n-1, n-2)
	
	left := FibonacciRecursiveWithSteps(n-1, depth+1)
	right := FibonacciRecursiveWithSteps(n-2, depth+1)
	
	result := left + right
	fmt.Printf("%sFib(%d) = %d + %d = %d\n", indent, n, left, right, result)
	
	return result
}

// FibonacciMemoization はメモ化を使った再帰実装
func FibonacciMemoization(n int) int {
	memo := make(map[int]int)
	return fibMemo(n, memo)
}

// fibMemo はメモ化のヘルパー関数
func fibMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	
	// メモに結果があれば使用
	if val, exists := memo[n]; exists {
		return val
	}
	
	// 計算してメモに保存
	result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
	memo[n] = result
	return result
}

// FibonacciMemoizationWithSteps はメモ化の過程を表示
func FibonacciMemoizationWithSteps(n int) int {
	memo := make(map[int]int)
	fmt.Printf("メモ化を使ったフィボナッチ数列の計算: Fib(%d)\n", n)
	return fibMemoWithSteps(n, memo, 0)
}

// fibMemoWithSteps はメモ化の過程を表示するヘルパー関数
func fibMemoWithSteps(n int, memo map[int]int, depth int) int {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	if n <= 1 {
		fmt.Printf("%sFib(%d) = %d (ベースケース)\n", indent, n, n)
		return n
	}
	
	// メモに結果があるかチェック
	if val, exists := memo[n]; exists {
		fmt.Printf("%sFib(%d) = %d (メモから取得)\n", indent, n, val)
		return val
	}
	
	fmt.Printf("%sFib(%d) を計算 = Fib(%d) + Fib(%d)\n", indent, n, n-1, n-2)
	
	left := fibMemoWithSteps(n-1, memo, depth+1)
	right := fibMemoWithSteps(n-2, memo, depth+1)
	
	result := left + right
	memo[n] = result
	
	fmt.Printf("%sFib(%d) = %d + %d = %d (メモに保存)\n", indent, n, left, right, result)
	
	return result
}

// FibonacciIterative は反復による実装（比較用）
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// FibonacciSequence はn番目までのフィボナッチ数列を生成
func FibonacciSequence(n int) []int {
	if n < 0 {
		return []int{}
	}
	
	sequence := make([]int, n+1)
	for i := 0; i <= n; i++ {
		sequence[i] = FibonacciRecursive(i)
	}
	return sequence
}

// FibonacciSequenceIterative は反復でn番目までの数列を生成
func FibonacciSequenceIterative(n int) []int {
	if n < 0 {
		return []int{}
	}
	
	sequence := make([]int, n+1)
	for i := 0; i <= n; i++ {
		sequence[i] = FibonacciIterative(i)
	}
	return sequence
}

// measureTime は関数の実行時間を測定
func measureTime(name string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%s の実行時間: %v\n", name, duration)
}

// RunFibonacciDemo はフィボナッチ数列のデモを実行
func RunFibonacciDemo() {
	fmt.Println("フィボナッチ数列の基本計算:")
	
	// 基本的な計算例
	for i := 0; i <= 10; i++ {
		result := FibonacciRecursive(i)
		fmt.Printf("Fib(%d) = %d\n", i, result)
	}
	
	// 再帰の呼び出し過程を表示
	fmt.Println("\n--- 再帰の呼び出し過程（Fib(5)）---")
	FibonacciRecursiveWithSteps(5, 0)
	
	// メモ化の効果を表示
	fmt.Println("\n--- メモ化の効果（Fib(6)）---")
	FibonacciMemoizationWithSteps(6)
	
	// パフォーマンス比較
	fmt.Println("\n--- パフォーマンス比較 ---")
	n := 35
	
	fmt.Printf("n = %d でのパフォーマンス比較:\n", n)
	
	measureTime("再帰版", func() {
		result := FibonacciRecursive(n)
		fmt.Printf("再帰版 Fib(%d) = %d\n", n, result)
	})
	
	measureTime("メモ化版", func() {
		result := FibonacciMemoization(n)
		fmt.Printf("メモ化版 Fib(%d) = %d\n", n, result)
	})
	
	measureTime("反復版", func() {
		result := FibonacciIterative(n)
		fmt.Printf("反復版 Fib(%d) = %d\n", n, result)
	})
	
	// 数列の生成
	fmt.Println("\n--- フィボナッチ数列（最初の15項）---")
	sequence := FibonacciSequenceIterative(14)
	fmt.Printf("数列: %v\n", sequence)
	
	// 大きな値での比較（メモ化と反復のみ）
	fmt.Println("\n--- 大きな値での計算 ---")
	largeN := 50
	
	measureTime(fmt.Sprintf("メモ化版 Fib(%d)", largeN), func() {
		result := FibonacciMemoization(largeN)
		fmt.Printf("メモ化版 Fib(%d) = %d\n", largeN, result)
	})
	
	measureTime(fmt.Sprintf("反復版 Fib(%d)", largeN), func() {
		result := FibonacciIterative(largeN)
		fmt.Printf("反復版 Fib(%d) = %d\n", largeN, result)
	})
	
	// 計算量の違いを示すデモ
	fmt.Println("\n--- 計算量の違いを実感 ---")
	fmt.Println("再帰版は指数的に遅くなるため、小さな値での比較:")
	
	for testN := 30; testN <= 35; testN += 5 {
		fmt.Printf("\nFib(%d) の計算時間:\n", testN)
		
		measureTime("再帰版", func() {
			FibonacciRecursive(testN)
		})
		
		measureTime("メモ化版", func() {
			FibonacciMemoization(testN)
		})
		
		measureTime("反復版", func() {
			FibonacciIterative(testN)
		})
	}
}