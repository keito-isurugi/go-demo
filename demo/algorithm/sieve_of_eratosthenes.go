package algorithm

import "fmt"

// SieveOfEratosthenes はエラトステネスの篩を使ってn以下の素数をすべて求める
func SieveOfEratosthenes(n int) []int {
	if n < 2 {
		return []int{}
	}

	// すべての数を素数候補とする（true = 素数）
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	// 2からsqrt(n)まで処理
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			// iの倍数をすべて合成数としてマーク
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	// 素数を収集
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// SieveOfEratosthenesWithSteps はステップごとに計算過程を表示
func SieveOfEratosthenesWithSteps(n int) []int {
	fmt.Printf("\nn以下の素数を求める (n=%d)\n", n)
	fmt.Println("=================================")

	if n < 2 {
		fmt.Println("結果: 素数なし")
		return []int{}
	}

	// すべての数を素数候補とする
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	fmt.Printf("\n初期状態: 2〜%dをすべて素数候補とする\n", n)
	printSieveState(isPrime, n)

	fmt.Println("\nエラトステネスの篩を開始:")
	fmt.Println("=================================")

	// 2からsqrt(n)まで処理
	step := 1
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			fmt.Printf("\nStep %d: %dは素数\n", step, i)
			fmt.Printf("        %dの倍数を合成数としてマーク: ", i)

			// iの倍数をすべて合成数としてマーク
			markedNums := []int{}
			for j := i * i; j <= n; j += i {
				if isPrime[j] {
					isPrime[j] = false
					markedNums = append(markedNums, j)
				}
			}

			if len(markedNums) > 0 {
				fmt.Printf("%v\n", markedNums)
			} else {
				fmt.Println("なし")
			}

			step++
		}
	}

	fmt.Println("\n最終状態:")
	printSieveState(isPrime, n)

	// 素数を収集
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	fmt.Printf("\n結果: %d個の素数が見つかりました\n", len(primes))
	fmt.Printf("素数: %v\n", primes)

	return primes
}

// printSieveState は篩の状態を表示
func printSieveState(isPrime []bool, n int) {
	fmt.Print("        ")
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()
}

// CountPrimes はn以下の素数の個数を返す
func CountPrimes(n int) int {
	if n < 2 {
		return 0
	}

	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	count := 0
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			count++
		}
	}

	return count
}

// IsPrime は素数判定を行う（エラトステネスの篩を使用）
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}

	primes := SieveOfEratosthenes(n)
	for _, p := range primes {
		if p == n {
			return true
		}
	}

	return false
}

// IsPrimeSimple は単純な素数判定（試し割り法）
func IsPrimeSimple(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// SieveOfEratosthenesRange は範囲[low, high]の素数を求める
func SieveOfEratosthenesRange(low, high int) []int {
	if high < 2 {
		return []int{}
	}

	// まず√highまでの素数を求める
	limit := int(float64(high) * 0.5)
	if limit < 2 {
		limit = 2
	}

	smallPrimes := SieveOfEratosthenes(limit)

	// 範囲[low, high]の素数候補
	isPrime := make([]bool, high-low+1)
	for i := range isPrime {
		isPrime[i] = true
	}

	// 小さな素数の倍数を除外
	for _, p := range smallPrimes {
		// pの倍数の開始位置を計算
		start := ((low + p - 1) / p) * p
		if start < p*p {
			start = p * p
		}

		for j := start; j <= high; j += p {
			isPrime[j-low] = false
		}
	}

	// 素数を収集
	primes := []int{}
	start := low
	if start < 2 {
		start = 2
	}

	for i := start; i <= high; i++ {
		if isPrime[i-low] {
			primes = append(primes, i)
		}
	}

	return primes
}

// DemoSieveOfEratosthenes はエラトステネスの篩のデモを実行
func DemoSieveOfEratosthenes() {
	fmt.Println("\n1. 基本的なエラトステネスの篩:")
	fmt.Println("=================================")
	SieveOfEratosthenesWithSteps(30)

	fmt.Println("\n\n2. より大きな範囲での素数探索:")
	fmt.Println("=================================")
	n := 100
	primes := SieveOfEratosthenes(n)
	fmt.Printf("n=%d以下の素数: %d個\n", n, len(primes))
	fmt.Printf("素数: %v\n", primes)

	fmt.Println("\n\n3. 素数の個数を数える:")
	fmt.Println("=================================")
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		count := CountPrimes(tc)
		fmt.Printf("n=%d以下の素数: %d個\n", tc, count)
	}

	fmt.Println("\n\n4. 素数判定:")
	fmt.Println("=================================")
	testNums := []int{2, 3, 4, 17, 18, 97, 100}
	for _, num := range testNums {
		if IsPrimeSimple(num) {
			fmt.Printf("%d は素数です\n", num)
		} else {
			fmt.Printf("%d は合成数です\n", num)
		}
	}

	fmt.Println("\n\n5. 範囲指定での素数探索:")
	fmt.Println("=================================")
	low, high := 50, 100
	rangePrimes := SieveOfEratosthenesRange(low, high)
	fmt.Printf("[%d, %d]の素数: %d個\n", low, high, len(rangePrimes))
	fmt.Printf("素数: %v\n", rangePrimes)

	fmt.Println("\n\n6. 性能比較（エラトステネスの篩 vs 試し割り法）:")
	fmt.Println("=================================")
	n2 := 10000
	fmt.Printf("n=%d以下の素数を求める:\n", n2)

	// エラトステネスの篩
	primes2 := SieveOfEratosthenes(n2)
	fmt.Printf("エラトステネスの篩: %d個の素数を発見\n", len(primes2))

	// 試し割り法
	count := 0
	for i := 2; i <= n2; i++ {
		if IsPrimeSimple(i) {
			count++
		}
	}
	fmt.Printf("試し割り法: %d個の素数を発見\n", count)

	fmt.Println("\n計算量:")
	fmt.Println("  エラトステネスの篩: O(n log log n) - 高速 ✓")
	fmt.Println("  試し割り法: O(n√n) - 遅い")

	fmt.Println("\n\n7. 双子素数の探索:")
	fmt.Println("=================================")
	primes3 := SieveOfEratosthenes(100)
	fmt.Println("100以下の双子素数（差が2の素数ペア）:")
	for i := 0; i < len(primes3)-1; i++ {
		if primes3[i+1]-primes3[i] == 2 {
			fmt.Printf("(%d, %d)\n", primes3[i], primes3[i+1])
		}
	}

	fmt.Println("\n\n8. 素数の分布:")
	fmt.Println("=================================")
	ranges := []struct {
		low, high int
	}{
		{1, 10},
		{11, 20},
		{21, 30},
		{31, 40},
		{41, 50},
		{51, 60},
		{61, 70},
		{71, 80},
		{81, 90},
		{91, 100},
	}

	for _, r := range ranges {
		rangePrimes2 := SieveOfEratosthenesRange(r.low, r.high)
		fmt.Printf("[%2d-%3d]: %d個 - %v\n", r.low, r.high, len(rangePrimes2), rangePrimes2)
	}
}
