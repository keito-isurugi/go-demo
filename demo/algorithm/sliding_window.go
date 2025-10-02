package algorithm

import "fmt"

// MaxSumSubarray は固定長kの部分配列の最大合計を求める（スライディングウィンドウ）
func MaxSumSubarray(arr []int, k int) int {
	n := len(arr)
	if n < k {
		return 0
	}

	// 最初のウィンドウの合計を計算
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}

	maxSum := windowSum

	// ウィンドウをスライドさせながら最大値を更新
	for i := k; i < n; i++ {
		windowSum = windowSum + arr[i] - arr[i-k]
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum
}

// MaxSumSubarrayWithSteps はステップごとに計算過程を表示
func MaxSumSubarrayWithSteps(arr []int, k int) int {
	n := len(arr)
	fmt.Printf("\n配列: %v, ウィンドウサイズ: %d\n", arr, k)
	fmt.Println("=================================")

	if n < k {
		fmt.Println("配列の長さがウィンドウサイズより小さいため、計算できません")
		return 0
	}

	// 最初のウィンドウの合計を計算
	windowSum := 0
	fmt.Printf("\n最初のウィンドウ [0:%d] を計算:\n", k)
	for i := 0; i < k; i++ {
		windowSum += arr[i]
		fmt.Printf("  arr[%d] = %d, 累積和 = %d\n", i, arr[i], windowSum)
	}

	maxSum := windowSum
	maxIndex := 0
	fmt.Printf("初期ウィンドウ合計: %d\n", windowSum)

	// ウィンドウをスライドさせながら最大値を更新
	fmt.Println("\nウィンドウをスライド:")
	for i := k; i < n; i++ {
		oldVal := arr[i-k]
		newVal := arr[i]
		windowSum = windowSum + newVal - oldVal

		fmt.Printf("  [%d:%d] 削除: arr[%d]=%d, 追加: arr[%d]=%d → 合計=%d",
			i-k+1, i+1, i-k, oldVal, i, newVal, windowSum)

		if windowSum > maxSum {
			maxSum = windowSum
			maxIndex = i - k + 1
			fmt.Printf(" ← 最大値更新!\n")
		} else {
			fmt.Println()
		}
	}

	fmt.Printf("\n結果: 最大合計 = %d (位置 [%d:%d])\n", maxSum, maxIndex, maxIndex+k)
	return maxSum
}

// MinSubarrayLen は合計がtarget以上になる最小長の部分配列を求める（可変長スライディングウィンドウ)
func MinSubarrayLen(arr []int, target int) int {
	n := len(arr)
	minLen := n + 1
	windowSum := 0
	left := 0

	for right := 0; right < n; right++ {
		windowSum += arr[right]

		// 条件を満たす間、左端を縮める
		for windowSum >= target {
			currentLen := right - left + 1
			if currentLen < minLen {
				minLen = currentLen
			}
			windowSum -= arr[left]
			left++
		}
	}

	if minLen == n+1 {
		return 0 // 見つからない
	}
	return minLen
}

// MinSubarrayLenWithSteps はステップごとに計算過程を表示
func MinSubarrayLenWithSteps(arr []int, target int) int {
	n := len(arr)
	fmt.Printf("\n配列: %v, 目標値: %d\n", arr, target)
	fmt.Println("=================================")

	minLen := n + 1
	windowSum := 0
	left := 0

	fmt.Println("\nウィンドウの拡大と縮小:")
	for right := 0; right < n; right++ {
		windowSum += arr[right]
		fmt.Printf("  右端拡大: arr[%d]=%d → [%d:%d] 合計=%d",
			right, arr[right], left, right+1, windowSum)

		// 条件を満たす間、左端を縮める
		if windowSum >= target {
			fmt.Printf(" ✓ (目標達成)\n")

			for windowSum >= target {
				currentLen := right - left + 1
				if currentLen < minLen {
					minLen = currentLen
					fmt.Printf("    最小長更新: %d (位置 [%d:%d])\n", minLen, left, right+1)
				}
				fmt.Printf("    左端縮小: arr[%d]=%d を削除 → [%d:%d] 合計=%d\n",
					left, arr[left], left+1, right+1, windowSum-arr[left])
				windowSum -= arr[left]
				left++
			}
		} else {
			fmt.Println()
		}
	}

	if minLen == n+1 {
		fmt.Println("\n結果: 条件を満たす部分配列は見つかりませんでした")
		return 0
	}
	fmt.Printf("\n結果: 最小長 = %d\n", minLen)
	return minLen
}

// LongestSubstringKDistinct は最大k種類の異なる文字を含む最長部分文字列の長さを求める
func LongestSubstringKDistinct(s string, k int) int {
	if k == 0 || len(s) == 0 {
		return 0
	}

	charCount := make(map[byte]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		charCount[s[right]]++

		// 異なる文字の種類がkを超えたら左端を縮める
		for len(charCount) > k {
			charCount[s[left]]--
			if charCount[s[left]] == 0 {
				delete(charCount, s[left])
			}
			left++
		}

		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
		}
	}

	return maxLen
}

// LongestSubstringKDistinctWithSteps はステップごとに計算過程を表示
func LongestSubstringKDistinctWithSteps(s string, k int) int {
	fmt.Printf("\n文字列: \"%s\", 最大種類数: %d\n", s, k)
	fmt.Println("=================================")

	if k == 0 || len(s) == 0 {
		fmt.Println("結果: 0")
		return 0
	}

	charCount := make(map[byte]int)
	maxLen := 0
	left := 0
	maxStart := 0

	fmt.Println("\nウィンドウの拡大と縮小:")
	for right := 0; right < len(s); right++ {
		charCount[s[right]]++

		fmt.Printf("  右端拡大: '%c' → [%d:%d] \"%s\" 文字種類=%v",
			s[right], left, right+1, s[left:right+1], getCharSet(charCount))

		// 異なる文字の種類がkを超えたら左端を縮める
		if len(charCount) > k {
			fmt.Printf(" ✗ (種類数超過)\n")

			for len(charCount) > k {
				charCount[s[left]]--
				if charCount[s[left]] == 0 {
					fmt.Printf("    左端縮小: '%c' を削除 → ", s[left])
					delete(charCount, s[left])
				} else {
					fmt.Printf("    左端縮小: '%c' のカウント減少 → ", s[left])
				}
				left++
				fmt.Printf("[%d:%d] \"%s\" 文字種類=%v\n",
					left, right+1, s[left:right+1], getCharSet(charCount))
			}
		} else {
			fmt.Println()
		}

		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
			maxStart = left
			fmt.Printf("    最長更新: %d (位置 [%d:%d] \"%s\")\n",
				maxLen, left, right+1, s[left:right+1])
		}
	}

	fmt.Printf("\n結果: 最長部分文字列 = \"%s\" (長さ %d)\n", s[maxStart:maxStart+maxLen], maxLen)
	return maxLen
}

// getCharSet はマップのキーを文字列として返す
func getCharSet(m map[byte]int) string {
	chars := ""
	for k := range m {
		chars += string(k)
	}
	return "{" + chars + "}"
}

// MaxSlidingWindow は固定長kのウィンドウをスライドさせて各位置での最大値を求める
func MaxSlidingWindow(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}

	result := make([]int, 0, n-k+1)

	// 各ウィンドウの最大値を計算
	for i := 0; i <= n-k; i++ {
		max := nums[i]
		for j := i + 1; j < i+k; j++ {
			if nums[j] > max {
				max = nums[j]
			}
		}
		result = append(result, max)
	}

	return result
}

// DemoSlidingWindow はスライディングウィンドウのデモを実行
func DemoSlidingWindow() {
	fmt.Println("\n1. 固定長部分配列の最大合計:")
	fmt.Println("=================================")
	arr1 := []int{2, 1, 5, 1, 3, 2}
	k1 := 3
	MaxSumSubarrayWithSteps(arr1, k1)

	arr2 := []int{2, 3, 4, 1, 5}
	k2 := 2
	MaxSumSubarrayWithSteps(arr2, k2)

	fmt.Println("\n\n2. 合計が目標値以上になる最小長部分配列:")
	fmt.Println("=================================")
	arr3 := []int{2, 3, 1, 2, 4, 3}
	target1 := 7
	MinSubarrayLenWithSteps(arr3, target1)

	arr4 := []int{1, 4, 4}
	target2 := 4
	MinSubarrayLenWithSteps(arr4, target2)

	fmt.Println("\n\n3. 最大k種類の異なる文字を含む最長部分文字列:")
	fmt.Println("=================================")
	s1 := "eceba"
	k3 := 2
	LongestSubstringKDistinctWithSteps(s1, k3)

	s2 := "aa"
	k4 := 1
	LongestSubstringKDistinctWithSteps(s2, k4)

	fmt.Println("\n\n4. スライディングウィンドウの最大値:")
	fmt.Println("=================================")
	arr5 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k5 := 3
	result := MaxSlidingWindow(arr5, k5)
	fmt.Printf("配列: %v, ウィンドウサイズ: %d\n", arr5, k5)
	fmt.Printf("各ウィンドウの最大値: %v\n", result)

	// 詳細表示
	fmt.Println("\n詳細:")
	for i := 0; i <= len(arr5)-k5; i++ {
		window := arr5[i : i+k5]
		fmt.Printf("  [%d:%d] %v → 最大値 = %d\n", i, i+k5, window, result[i])
	}

	fmt.Println("\n\n5. 性能比較（固定長最大合計）:")
	fmt.Println("=================================")

	// 大きな配列でのテスト
	largeArr := make([]int, 100000)
	for i := range largeArr {
		largeArr[i] = i % 100
	}
	k := 1000

	maxSum := MaxSumSubarray(largeArr, k)
	fmt.Printf("配列サイズ: %d, ウィンドウサイズ: %d\n", len(largeArr), k)
	fmt.Printf("最大合計: %d\n", maxSum)
	fmt.Println("スライディングウィンドウ: O(n) - 高速 ✓")
	fmt.Println("総当たり: O(n×k) - 遅い")
}
