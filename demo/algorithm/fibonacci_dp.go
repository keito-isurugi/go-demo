package main

import "fmt"

// FibonacciDP はボトムアップ動的計画法でフィボナッチ数列のn番目を計算
func FibonacciDP(n int) int {
	if n <= 1 {
		return n
	}

	// dpテーブルを作成（0からnまで）
	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1

	// ボトムアップで計算
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

// FibonacciDPWithSteps はステップごとに計算過程を表示
func FibonacciDPWithSteps(n int) int {
	fmt.Printf("\nFibonacci(%d) を動的計画法で計算:\n", n)
	fmt.Println("=================================")

	if n <= 1 {
		fmt.Printf("n <= 1 のため、結果は %d\n", n)
		return n
	}

	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1

	fmt.Printf("初期値: dp[0] = %d, dp[1] = %d\n", dp[0], dp[1])
	fmt.Println("\n計算過程:")

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
		fmt.Printf("dp[%d] = dp[%d] + dp[%d] = %d + %d = %d\n",
			i, i-1, i-2, dp[i-1], dp[i-2], dp[i])
	}

	fmt.Printf("\n最終的なdpテーブル: %v\n", dp)
	fmt.Printf("結果: Fibonacci(%d) = %d\n", n, dp[n])

	return dp[n]
}

// FibonacciDPOptimized は空間計算量を最適化した動的計画法
func FibonacciDPOptimized(n int) int {
	if n <= 1 {
		return n
	}

	// 直前の2つの値のみ保持（空間計算量O(1)）
	prev2, prev1 := 0, 1

	for i := 2; i <= n; i++ {
		current := prev1 + prev2
		prev2, prev1 = prev1, current
	}

	return prev1
}

// FibonacciDPOptimizedWithSteps は最適化版のステップ表示
func FibonacciDPOptimizedWithSteps(n int) int {
	fmt.Printf("\nFibonacci(%d) を最適化動的計画法で計算:\n", n)
	fmt.Println("=================================")

	if n <= 1 {
		fmt.Printf("n <= 1 のため、結果は %d\n", n)
		return n
	}

	prev2, prev1 := 0, 1
	fmt.Printf("初期値: prev2 = %d, prev1 = %d\n", prev2, prev1)
	fmt.Println("\n計算過程:")

	for i := 2; i <= n; i++ {
		current := prev1 + prev2
		fmt.Printf("i=%d: current = prev1 + prev2 = %d + %d = %d\n",
			i, prev1, prev2, current)
		prev2, prev1 = prev1, current
		fmt.Printf("      更新後: prev2 = %d, prev1 = %d\n", prev2, prev1)
	}

	fmt.Printf("\n結果: Fibonacci(%d) = %d\n", n, prev1)
	return prev1
}

// FibonacciSequenceDP は動的計画法でn番目までのフィボナッチ数列を生成
func FibonacciSequenceDP(n int) []int {
	if n < 0 {
		return []int{}
	}

	if n == 0 {
		return []int{0}
	}

	sequence := make([]int, n+1)
	sequence[0] = 0
	sequence[1] = 1

	for i := 2; i <= n; i++ {
		sequence[i] = sequence[i-1] + sequence[i-2]
	}

	return sequence
}

// DemoFibonacciDP はフィボナッチ数列（動的計画法）のデモを実行
func DemoFibonacciDP() {
	fmt.Println("\n1. 基本的な動的計画法:")
	fmt.Println("=================================")

	testCases := []int{0, 1, 5, 10, 15, 20}
	for _, n := range testCases {
		result := FibonacciDP(n)
		fmt.Printf("Fibonacci(%d) = %d\n", n, result)
	}

	fmt.Println("\n2. ステップごとの計算過程（ボトムアップDP）:")
	fmt.Println("=================================")
	FibonacciDPWithSteps(8)

	fmt.Println("\n3. 最適化版動的計画法（空間O(1)）:")
	fmt.Println("=================================")
	FibonacciDPOptimizedWithSteps(8)

	fmt.Println("\n4. メモ化（トップダウンDP）:")
	fmt.Println("=================================")
	FibonacciMemoizationWithSteps(8) // 既存のfibonacci.goの関数を使用

	fmt.Println("\n5. フィボナッチ数列の生成:")
	fmt.Println("=================================")
	sequence := FibonacciSequenceDP(15)
	fmt.Printf("Fibonacci数列（0-15）: %v\n", sequence)

	fmt.Println("\n6. 大きな値の計算:")
	fmt.Println("=================================")
	largeTests := []int{30, 40, 50}
	for _, n := range largeTests {
		result := FibonacciDPOptimized(n)
		fmt.Printf("Fibonacci(%d) = %d\n", n, result)
	}

	fmt.Println("\n7. 実装比較:")
	fmt.Println("=================================")
	n := 10
	dpResult := FibonacciDP(n)
	optimizedResult := FibonacciDPOptimized(n)
	memoResult := FibonacciMemoization(n) // 既存のfibonacci.goの関数を使用

	fmt.Printf("ボトムアップDP:   Fibonacci(%d) = %d\n", n, dpResult)
	fmt.Printf("最適化DP:         Fibonacci(%d) = %d\n", n, optimizedResult)
	fmt.Printf("メモ化:           Fibonacci(%d) = %d\n", n, memoResult)
}
