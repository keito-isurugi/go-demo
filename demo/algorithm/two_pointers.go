package algorithm

import "fmt"

// TwoSumSorted は昇順配列から合計がtargetになる2要素のインデックスを返す（Two Pointers）
func TwoSumSorted(nums []int, target int) []int {
	left := 0
	right := len(nums) - 1

	for left < right {
		sum := nums[left] + nums[right]
		if sum == target {
			return []int{left, right}
		} else if sum < target {
			left++ // 合計を増やすために左端を右に移動
		} else {
			right-- // 合計を減らすために右端を左に移動
		}
	}

	return []int{} // 見つからない場合
}

// TwoSumSortedWithSteps はステップごとに計算過程を表示
func TwoSumSortedWithSteps(nums []int, target int) []int {
	fmt.Printf("\n配列: %v, 目標値: %d\n", nums, target)
	fmt.Println("=================================")

	left := 0
	right := len(nums) - 1

	fmt.Println("\nTwo Pointersによる探索:")
	step := 1

	for left < right {
		sum := nums[left] + nums[right]
		fmt.Printf("Step %d: left=%d, right=%d → nums[%d]=%d + nums[%d]=%d = %d",
			step, left, right, left, nums[left], right, nums[right], sum)

		if sum == target {
			fmt.Printf(" ✓ 見つかりました!\n")
			fmt.Printf("\n結果: インデックス [%d, %d]\n", left, right)
			return []int{left, right}
		} else if sum < target {
			fmt.Printf(" (合計が小さい → left++)\n")
			left++
		} else {
			fmt.Printf(" (合計が大きい → right--)\n")
			right--
		}
		step++
	}

	fmt.Println("\n結果: 見つかりませんでした")
	return []int{}
}

// RemoveDuplicates はソート済み配列から重複を削除し、ユニーク要素の数を返す
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	slow := 0 // ユニーク要素を配置する位置

	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}

	return slow + 1
}

// RemoveDuplicatesWithSteps はステップごとに計算過程を表示
func RemoveDuplicatesWithSteps(nums []int) int {
	fmt.Printf("\n配列: %v\n", nums)
	fmt.Println("=================================")

	if len(nums) == 0 {
		fmt.Println("結果: 0")
		return 0
	}

	slow := 0

	fmt.Println("\nTwo Pointersによる重複削除:")
	fmt.Printf("初期状態: slow=%d, nums=%v\n", slow, nums)

	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
			fmt.Printf("Step %d: fast=%d, nums[fast]=%d ≠ nums[slow-1]=%d → nums[%d] = %d\n",
				fast, fast, nums[fast], nums[slow-1], slow, nums[slow])
			fmt.Printf("        配列: %v\n", nums[:slow+1])
		} else {
			fmt.Printf("Step %d: fast=%d, nums[fast]=%d = nums[slow]=%d → スキップ\n",
				fast, fast, nums[fast], nums[slow])
		}
	}

	uniqueCount := slow + 1
	fmt.Printf("\n結果: ユニーク要素数 = %d\n", uniqueCount)
	fmt.Printf("最終配列: %v\n", nums[:uniqueCount])

	return uniqueCount
}

// ReverseString は文字列を反転する（Two Pointers）
func ReverseString(s []byte) {
	left := 0
	right := len(s) - 1

	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}

// ReverseStringWithSteps はステップごとに計算過程を表示
func ReverseStringWithSteps(s []byte) {
	fmt.Printf("\n文字列: \"%s\"\n", string(s))
	fmt.Println("=================================")

	left := 0
	right := len(s) - 1

	fmt.Println("\nTwo Pointersによる文字列反転:")
	step := 1

	for left < right {
		fmt.Printf("Step %d: swap s[%d]='%c' ↔ s[%d]='%c' → \"%s\"\n",
			step, left, s[left], right, s[right], string(s))

		s[left], s[right] = s[right], s[left]
		left++
		right--
		step++
	}

	fmt.Printf("\n結果: \"%s\"\n", string(s))
}

// IsPalindrome は文字列が回文かどうかを判定する（Two Pointers）
func IsPalindrome(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}

	return true
}

// IsPalindromeWithSteps はステップごとに計算過程を表示
func IsPalindromeWithSteps(s string) bool {
	fmt.Printf("\n文字列: \"%s\"\n", s)
	fmt.Println("=================================")

	left := 0
	right := len(s) - 1

	fmt.Println("\nTwo Pointersによる回文判定:")
	step := 1

	for left < right {
		fmt.Printf("Step %d: s[%d]='%c' vs s[%d]='%c'",
			step, left, s[left], right, s[right])

		if s[left] != s[right] {
			fmt.Printf(" ✗ 不一致 → 回文ではありません\n")
			return false
		}

		fmt.Printf(" ✓ 一致\n")
		left++
		right--
		step++
	}

	fmt.Println("\n結果: 回文です")
	return true
}

// ThreeSum は配列から合計が0になる3要素の組み合わせをすべて返す
func ThreeSum(nums []int) [][]int {
	result := [][]int{}
	n := len(nums)

	// まずソート
	sortSlice(nums)

	for i := 0; i < n-2; i++ {
		// 重複をスキップ
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		// Two Pointersで残り2要素を探す
		left := i + 1
		right := n - 1
		target := -nums[i]

		for left < right {
			sum := nums[left] + nums[right]

			if sum == target {
				result = append(result, []int{nums[i], nums[left], nums[right]})

				// 重複をスキップ
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}

				left++
				right--
			} else if sum < target {
				left++
			} else {
				right--
			}
		}
	}

	return result
}

// sortSlice はスライスを昇順にソート（バブルソートで実装）
func sortSlice(nums []int) {
	n := len(nums)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
}

// ThreeSumWithSteps はステップごとに計算過程を表示
func ThreeSumWithSteps(nums []int) [][]int {
	fmt.Printf("\n配列: %v\n", nums)
	fmt.Println("=================================")

	result := [][]int{}
	n := len(nums)

	// まずソート
	sortSlice(nums)
	fmt.Printf("ソート後: %v\n", nums)

	fmt.Println("\nThree Sum（固定 + Two Pointers）:")

	for i := 0; i < n-2; i++ {
		// 重複をスキップ
		if i > 0 && nums[i] == nums[i-1] {
			fmt.Printf("i=%d: nums[%d]=%d は重複のためスキップ\n", i, i, nums[i])
			continue
		}

		fmt.Printf("\ni=%d: nums[%d]=%d を固定、残り2要素の合計目標値=%d\n",
			i, i, nums[i], -nums[i])

		// Two Pointersで残り2要素を探す
		left := i + 1
		right := n - 1
		target := -nums[i]

		for left < right {
			sum := nums[left] + nums[right]
			fmt.Printf("  left=%d, right=%d → %d + %d = %d",
				left, right, nums[left], nums[right], sum)

			if sum == target {
				triplet := []int{nums[i], nums[left], nums[right]}
				result = append(result, triplet)
				fmt.Printf(" ✓ 見つかりました: %v\n", triplet)

				// 重複をスキップ
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}

				left++
				right--
			} else if sum < target {
				fmt.Printf(" (合計が小さい → left++)\n")
				left++
			} else {
				fmt.Printf(" (合計が大きい → right--)\n")
				right--
			}
		}
	}

	fmt.Printf("\n結果: %v\n", result)
	return result
}

// ContainerWithMostWater は最大の水を貯められる2つの垂直線を見つける
func ContainerWithMostWater(height []int) int {
	maxArea := 0
	left := 0
	right := len(height) - 1

	for left < right {
		// 面積 = 幅 × 高さ（低い方）
		width := right - left
		h := minInt(height[left], height[right])
		area := width * h

		if area > maxArea {
			maxArea = area
		}

		// 低い方のポインタを移動
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}

	return maxArea
}

// minInt は2つの整数の最小値を返す
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ContainerWithMostWaterWithSteps はステップごとに計算過程を表示
func ContainerWithMostWaterWithSteps(height []int) int {
	fmt.Printf("\n高さ配列: %v\n", height)
	fmt.Println("=================================")

	maxArea := 0
	maxLeft := 0
	maxRight := 0
	left := 0
	right := len(height) - 1

	fmt.Println("\nTwo Pointersによる最大面積探索:")
	step := 1

	for left < right {
		width := right - left
		h := minInt(height[left], height[right])
		area := width * h

		fmt.Printf("Step %d: left=%d(h=%d), right=%d(h=%d) → 幅=%d × 高さ=%d = 面積=%d",
			step, left, height[left], right, height[right], width, h, area)

		if area > maxArea {
			maxArea = area
			maxLeft = left
			maxRight = right
			fmt.Printf(" ← 最大面積更新!\n")
		} else {
			fmt.Println()
		}

		// 低い方のポインタを移動
		if height[left] < height[right] {
			fmt.Printf("        height[%d]=%d < height[%d]=%d → left++\n",
				left, height[left], right, height[right])
			left++
		} else {
			fmt.Printf("        height[%d]=%d >= height[%d]=%d → right--\n",
				left, height[left], right, height[right])
			right--
		}

		step++
	}

	fmt.Printf("\n結果: 最大面積 = %d (位置 [%d, %d], 高さ [%d, %d])\n",
		maxArea, maxLeft, maxRight, height[maxLeft], height[maxRight])

	return maxArea
}

// MoveZeroes は配列内の0をすべて末尾に移動する（相対順序は保持）
func MoveZeroes(nums []int) {
	slow := 0 // 非0要素を配置する位置

	// すべての非0要素を左に詰める
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow] = nums[fast]
			slow++
		}
	}

	// 残りを0で埋める
	for i := slow; i < len(nums); i++ {
		nums[i] = 0
	}
}

// MoveZeroesWithSteps はステップごとに計算過程を表示
func MoveZeroesWithSteps(nums []int) {
	fmt.Printf("\n配列: %v\n", nums)
	fmt.Println("=================================")

	slow := 0

	fmt.Println("\nTwo Pointersによる0の移動:")
	fmt.Println("Phase 1: 非0要素を左に詰める")

	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow] = nums[fast]
			fmt.Printf("Step %d: fast=%d, nums[fast]=%d → nums[%d] = %d\n",
				fast+1, fast, nums[fast], slow, nums[slow])
			fmt.Printf("        配列: %v\n", nums)
			slow++
		} else {
			fmt.Printf("Step %d: fast=%d, nums[fast]=0 → スキップ\n", fast+1, fast)
		}
	}

	fmt.Println("\nPhase 2: 残りを0で埋める")
	for i := slow; i < len(nums); i++ {
		nums[i] = 0
		fmt.Printf("        nums[%d] = 0\n", i)
	}

	fmt.Printf("\n結果: %v\n", nums)
}

// DemoTwoPointers はTwo Pointersのデモを実行
func DemoTwoPointers() {
	fmt.Println("\n1. Two Sum（ソート済み配列）:")
	fmt.Println("=================================")
	nums1 := []int{2, 7, 11, 15}
	target1 := 9
	TwoSumSortedWithSteps(nums1, target1)

	nums2 := []int{2, 3, 4}
	target2 := 6
	TwoSumSortedWithSteps(nums2, target2)

	fmt.Println("\n\n2. 重複の削除（ソート済み配列）:")
	fmt.Println("=================================")
	nums3 := []int{1, 1, 2}
	RemoveDuplicatesWithSteps(nums3)

	nums4 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	RemoveDuplicatesWithSteps(nums4)

	fmt.Println("\n\n3. 文字列の反転:")
	fmt.Println("=================================")
	s1 := []byte("hello")
	ReverseStringWithSteps(s1)

	s2 := []byte("Hannah")
	ReverseStringWithSteps(s2)

	fmt.Println("\n\n4. 回文判定:")
	fmt.Println("=================================")
	IsPalindromeWithSteps("racecar")
	IsPalindromeWithSteps("hello")

	fmt.Println("\n\n5. Three Sum（3要素の合計が0）:")
	fmt.Println("=================================")
	nums5 := []int{-1, 0, 1, 2, -1, -4}
	ThreeSumWithSteps(nums5)

	nums6 := []int{0, 1, 1}
	ThreeSumWithSteps(nums6)

	fmt.Println("\n\n6. Container With Most Water:")
	fmt.Println("=================================")
	height1 := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	ContainerWithMostWaterWithSteps(height1)

	height2 := []int{1, 1}
	ContainerWithMostWaterWithSteps(height2)

	fmt.Println("\n\n7. Move Zeroes（0を末尾に移動）:")
	fmt.Println("=================================")
	nums7 := []int{0, 1, 0, 3, 12}
	MoveZeroesWithSteps(nums7)

	nums8 := []int{0}
	MoveZeroesWithSteps(nums8)

	fmt.Println("\n\n8. 性能比較:")
	fmt.Println("=================================")
	largeNums := make([]int, 100000)
	for i := range largeNums {
		largeNums[i] = i
	}
	target := 199999

	result := TwoSumSorted(largeNums, target)
	fmt.Printf("配列サイズ: %d, 目標値: %d\n", len(largeNums), target)
	fmt.Printf("見つかった位置: %v\n", result)
	fmt.Println("Two Pointers: O(n) - 高速 ✓")
	fmt.Println("総当たり: O(n²) - 遅い")
}
