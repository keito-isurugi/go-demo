package algorithm

import (
	"fmt"
)

// BitFullSearch はbit全探索を実行し、すべての部分集合を返す関数
// n: 要素数
// 戻り値: すべての部分集合のビット表現
func BitFullSearch(n int) []int {
	var result []int
	
	// 2^n 通りのパターンを全探索
	for bit := 0; bit < (1 << n); bit++ {
		result = append(result, bit)
	}
	
	return result
}

// PrintSubset は集合の要素とビット表現から部分集合を表示する関数
// items: 元の集合の要素
// bit: 部分集合のビット表現
func PrintSubset(items []string, bit int) {
	fmt.Printf("ビット表現: %b -> 部分集合: {", bit)
	
	first := true
	for i := 0; i < len(items); i++ {
		// i番目のビットが1かどうかチェック
		if (bit >> i) & 1 == 1 {
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(items[i])
			first = false
		}
	}
	fmt.Println("}")
}

// SubsetSumBitSearch は集合の中から合計が target になる部分集合を探す関数
// nums: 数値の集合
// target: 目標の合計値
// 戻り値: 条件を満たす部分集合のビット表現のスライス
func SubsetSumBitSearch(nums []int, target int) []int {
	n := len(nums)
	var result []int
	
	// 2^n 通りのパターンを全探索
	for bit := 0; bit < (1 << n); bit++ {
		sum := 0
		
		// 各ビットをチェックして、対応する要素を合計に加える
		for i := 0; i < n; i++ {
			if (bit >> i) & 1 == 1 {
				sum += nums[i]
			}
		}
		
		// 合計が目標値と一致する場合、結果に追加
		if sum == target {
			result = append(result, bit)
		}
	}
	
	return result
}

// CountSubsets は特定の条件を満たす部分集合の数を数える関数
// n: 要素数
// condition: 部分集合が満たすべき条件（ビット表現を受け取る関数）
// 戻り値: 条件を満たす部分集合の数
func CountSubsets(n int, condition func(int) bool) int {
	count := 0
	
	// 2^n 通りのパターンを全探索
	for bit := 0; bit < (1 << n); bit++ {
		if condition(bit) {
			count++
		}
	}
	
	return count
}

// GetSubsetElements は特定のビット表現から実際の要素を取得する関数
// items: 元の集合の要素
// bit: 部分集合のビット表現
// 戻り値: 部分集合の要素
func GetSubsetElements(items []interface{}, bit int) []interface{} {
	var subset []interface{}
	
	for i := 0; i < len(items); i++ {
		if (bit >> i) & 1 == 1 {
			subset = append(subset, items[i])
		}
	}
	
	return subset
}

// CountSubsetSum は集合の中から合計が target になる組み合わせの数を返す関数
// nums: 数値の集合
// target: 目標の合計値
// 戻り値: 条件を満たす組み合わせの数
func CountSubsetSum(nums []int, target int) int {
	n := len(nums)
	count := 0
	
	// 2^n 通りのパターンを全探索
	for bit := 0; bit < (1 << n); bit++ {
		sum := 0
		
		// 各ビットをチェックして、対応する要素を合計に加える
		for i := 0; i < n; i++ {
			if (bit >> i) & 1 == 1 {
				sum += nums[i]
			}
		}
		
		// 合計が目標値と一致する場合、カウントを増やす
		if sum == target {
			count++
		}
	}
	
	return count
}

// PrintSubsetSumDetail は合計が target になる組み合わせを詳細表示する関数
// nums: 数値の集合
// target: 目標の合計値
func PrintSubsetSumDetail(nums []int, target int) {
	n := len(nums)
	count := 0
	
	fmt.Printf("%d種類の数値 %v の内、合計が%dになる組み合わせ:\n", n, nums, target)
	
	// 2^n 通りのパターンを全探索
	for bit := 0; bit < (1 << n); bit++ {
		sum := 0
		var combination []int
		
		// 各ビットをチェックして、対応する要素を合計に加える
		for i := 0; i < n; i++ {
			if (bit >> i) & 1 == 1 {
				sum += nums[i]
				combination = append(combination, nums[i])
			}
		}
		
		// 合計が目標値と一致する場合、組み合わせを表示
		if sum == target {
			count++
			fmt.Printf("  組み合わせ%d: ", count)
			for i, num := range combination {
				if i > 0 {
					fmt.Print(" + ")
				}
				fmt.Print(num)
			}
			fmt.Printf(" = %d\n", target)
		}
	}
	
	fmt.Printf("\n合計%d通りの組み合わせが見つかりました。\n", count)
}