package main

import "fmt"

// SubsetSum は動的計画法で部分和問題を解く
// nums: 数値の配列
// target: 目標とする合計値
// 戻り値: targetを作れる場合はtrue、作れない場合はfalse
func SubsetSum(nums []int, target int) bool {
	n := len(nums)

	// dp[i][j]: nums[0..i-1]の要素を使って合計jを作れるか
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
	}

	// 初期条件: 合計0は常に作れる（何も選ばない）
	for i := 0; i <= n; i++ {
		dp[i][0] = true
	}

	// ボトムアップで計算
	for i := 1; i <= n; i++ {
		for j := 0; j <= target; j++ {
			// i番目の要素を使わない場合
			dp[i][j] = dp[i-1][j]

			// i番目の要素を使う場合（値がjより小さい場合のみ）
			if j >= nums[i-1] {
				dp[i][j] = dp[i][j] || dp[i-1][j-nums[i-1]]
			}
		}
	}

	return dp[n][target]
}

// SubsetSumWithSteps は計算過程を表示しながら部分和問題を解く
func SubsetSumWithSteps(nums []int, target int) bool {
	fmt.Printf("\n部分和問題の計算過程:\n")
	fmt.Printf("配列: %v\n", nums)
	fmt.Printf("目標合計: %d\n", target)
	fmt.Println("=================================")

	n := len(nums)
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
	}

	// 初期条件
	for i := 0; i <= n; i++ {
		dp[i][0] = true
	}
	fmt.Println("\n初期条件: 合計0は常に作れる（何も選ばない）")

	// dpテーブルの構築
	fmt.Println("\n計算過程:")
	for i := 1; i <= n; i++ {
		fmt.Printf("\n要素 nums[%d] = %d を考慮:\n", i-1, nums[i-1])
		for j := 0; j <= target; j++ {
			// i番目の要素を使わない場合
			dp[i][j] = dp[i-1][j]

			// i番目の要素を使う場合
			if j >= nums[i-1] {
				useCurrent := dp[i-1][j-nums[i-1]]
				dp[i][j] = dp[i][j] || useCurrent

				if useCurrent && !dp[i-1][j] {
					fmt.Printf("  合計%d: 要素%dを追加 (前の状態から%dを作れる)\n",
						j, nums[i-1], j-nums[i-1])
				}
			}
		}
	}

	// dpテーブルの表示
	fmt.Println("\nDPテーブル (◯: 作れる, ×: 作れない):")
	fmt.Print("     ")
	for j := 0; j <= target; j++ {
		fmt.Printf("%3d", j)
	}
	fmt.Println()

	for i := 0; i <= n; i++ {
		if i == 0 {
			fmt.Print("初期 ")
		} else {
			fmt.Printf("%3d  ", nums[i-1])
		}
		for j := 0; j <= target; j++ {
			if dp[i][j] {
				fmt.Print("  ◯")
			} else {
				fmt.Print("  ×")
			}
		}
		fmt.Println()
	}

	result := dp[n][target]
	fmt.Printf("\n結果: 合計%dは", target)
	if result {
		fmt.Println("作れます ✓")
	} else {
		fmt.Println("作れません ✗")
	}

	return result
}

// FindSubsetSum は部分和を作る実際の要素を見つける
func FindSubsetSum(nums []int, target int) []int {
	n := len(nums)
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
	}

	for i := 0; i <= n; i++ {
		dp[i][0] = true
	}

	for i := 1; i <= n; i++ {
		for j := 0; j <= target; j++ {
			dp[i][j] = dp[i-1][j]
			if j >= nums[i-1] {
				dp[i][j] = dp[i][j] || dp[i-1][j-nums[i-1]]
			}
		}
	}

	// 解が存在しない場合
	if !dp[n][target] {
		return nil
	}

	// バックトラック: 実際に使った要素を復元
	result := []int{}
	i, j := n, target

	for i > 0 && j > 0 {
		// i番目の要素を使った場合
		if !dp[i-1][j] && j >= nums[i-1] && dp[i-1][j-nums[i-1]] {
			result = append(result, nums[i-1])
			j -= nums[i-1]
		}
		i--
	}

	return result
}

// FindSubsetSumWithSteps はバックトラックの過程を表示
func FindSubsetSumWithSteps(nums []int, target int) []int {
	n := len(nums)
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
	}

	for i := 0; i <= n; i++ {
		dp[i][0] = true
	}

	for i := 1; i <= n; i++ {
		for j := 0; j <= target; j++ {
			dp[i][j] = dp[i-1][j]
			if j >= nums[i-1] {
				dp[i][j] = dp[i][j] || dp[i-1][j-nums[i-1]]
			}
		}
	}

	if !dp[n][target] {
		fmt.Println("\n解は存在しません")
		return nil
	}

	fmt.Println("\nバックトラック（解の復元）:")
	fmt.Println("=================================")

	result := []int{}
	i, j := n, target

	for i > 0 && j > 0 {
		if !dp[i-1][j] && j >= nums[i-1] && dp[i-1][j-nums[i-1]] {
			fmt.Printf("要素 nums[%d] = %d を使用 (残り合計: %d → %d)\n",
				i-1, nums[i-1], j, j-nums[i-1])
			result = append(result, nums[i-1])
			j -= nums[i-1]
		} else {
			fmt.Printf("要素 nums[%d] = %d は使用しない\n", i-1, nums[i-1])
		}
		i--
	}

	fmt.Printf("\n選択した要素: %v\n", result)

	// 検証
	sum := 0
	for _, v := range result {
		sum += v
	}
	fmt.Printf("検証: 合計 = %d\n", sum)

	return result
}

// CountSubsetSums は目標値を作る方法の数を数える
func CountSubsetSums(nums []int, target int) int {
	n := len(nums)

	// dp[i][j]: nums[0..i-1]を使って合計jを作る方法の数
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, target+1)
	}

	// 初期条件: 合計0を作る方法は1通り（何も選ばない）
	for i := 0; i <= n; i++ {
		dp[i][0] = 1
	}

	// ボトムアップで計算
	for i := 1; i <= n; i++ {
		for j := 0; j <= target; j++ {
			// i番目の要素を使わない場合
			dp[i][j] = dp[i-1][j]

			// i番目の要素を使う場合
			if j >= nums[i-1] {
				dp[i][j] += dp[i-1][j-nums[i-1]]
			}
		}
	}

	return dp[n][target]
}

// SubsetSumOptimized は空間計算量を最適化した部分和問題
func SubsetSumOptimized(nums []int, target int) bool {
	// dp[j]: 合計jを作れるか
	dp := make([]bool, target+1)
	dp[0] = true

	for _, num := range nums {
		// 逆順で更新（同じ要素を複数回使わないため）
		for j := target; j >= num; j-- {
			if dp[j-num] {
				dp[j] = true
			}
		}
	}

	return dp[target]
}

// DemoSubsetSum は部分和問題のデモを実行
func DemoSubsetSum() {
	fmt.Println("\n部分和問題（動的計画法）のデモ")
	fmt.Println("=================================")

	fmt.Println("\n1. 基本的な部分和問題:")
	fmt.Println("=================================")

	testCases := []struct {
		nums   []int
		target int
	}{
		{[]int{3, 34, 4, 12, 5, 2}, 9},
		{[]int{3, 34, 4, 12, 5, 2}, 30},
		{[]int{1, 2, 3, 7}, 6},
		{[]int{1, 2, 7, 1, 5}, 10},
		{[]int{1, 5, 11, 5}, 11},
	}

	for _, tc := range testCases {
		result := SubsetSum(tc.nums, tc.target)
		fmt.Printf("配列: %v, 目標: %d → ", tc.nums, tc.target)
		if result {
			fmt.Println("可能 ✓")
		} else {
			fmt.Println("不可能 ✗")
		}
	}

	fmt.Println("\n2. ステップごとの計算過程:")
	fmt.Println("=================================")
	SubsetSumWithSteps([]int{3, 4, 5, 2}, 9)

	fmt.Println("\n3. 実際の要素を見つける:")
	fmt.Println("=================================")
	nums := []int{3, 34, 4, 12, 5, 2}
	target := 9

	subset := FindSubsetSum(nums, target)
	if subset != nil {
		fmt.Printf("配列: %v, 目標: %d\n", nums, target)
		fmt.Printf("解: %v\n", subset)
	}

	fmt.Println("\n4. バックトラックの詳細:")
	fmt.Println("=================================")
	FindSubsetSumWithSteps([]int{3, 4, 5, 2}, 9)

	fmt.Println("\n5. 解の個数を数える:")
	fmt.Println("=================================")
	countCases := []struct {
		nums   []int
		target int
	}{
		{[]int{1, 1, 2, 3}, 4},
		{[]int{1, 2, 3, 3}, 6},
		{[]int{2, 3, 5, 6, 8, 10}, 10},
	}

	for _, tc := range countCases {
		count := CountSubsetSums(tc.nums, tc.target)
		fmt.Printf("配列: %v, 目標: %d → %d通り\n", tc.nums, tc.target, count)
	}

	fmt.Println("\n6. 最適化版（空間O(target)）:")
	fmt.Println("=================================")
	optimizedCases := []struct {
		nums   []int
		target int
	}{
		{[]int{1, 5, 11, 5}, 11},
		{[]int{1, 2, 3, 7}, 6},
		{[]int{3, 34, 4, 12, 5, 2}, 30},
	}

	for _, tc := range optimizedCases {
		result := SubsetSumOptimized(tc.nums, tc.target)
		fmt.Printf("配列: %v, 目標: %d → ", tc.nums, tc.target)
		if result {
			fmt.Println("可能 ✓")
		} else {
			fmt.Println("不可能 ✗")
		}
	}

	fmt.Println("\n7. 実用例: 支払い問題:")
	fmt.Println("=================================")
	coins := []int{1, 5, 10, 25, 50, 100}
	prices := []int{87, 150, 63}

	for _, price := range prices {
		result := SubsetSum(coins, price)
		fmt.Printf("%d円を払えるか: ", price)
		if result {
			subset := FindSubsetSum(coins, price)
			fmt.Printf("可能 (使用硬貨: %v円)\n", subset)
		} else {
			fmt.Println("不可能")
		}
	}
}
