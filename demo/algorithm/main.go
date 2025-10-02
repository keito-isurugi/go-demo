package algorithm

import (
	"fmt"

	// "algorithm/recursion"
)

func main() {
	// arr := []int{5, 3, 8, 4, 2}
	// target := 8

	// ls := LinearSearch(arr, target)
	// fmt.Printf("Linear Search: %d\n", ls)

	// bas := BubbleAscSort(arr)
	// fmt.Printf("Bubble Asc Sort: %v\n", bas)

	// bds := BubbleDescSort(arr)
	// fmt.Printf("Bubble Desc Sort: %v\n", bds)

	// bs := BinarySearch(bas, 8)
	// fmt.Printf("Binary Search: %d\n", bs)

	// fmt.Println("\n--- Bit全探索のデモ ---")

	// 合計が特定の値になる組み合わせの数を求める
	// fmt.Println("\n2. 合計が特定の値になる組み合わせの数:")

	// 3種類の数値から合計4になる組み合わせ
	// nums1 := []int{1, 2, 3}
	// target1 := 4
	// count1 := CountSubsetSum(nums1, target1)
	// fmt.Println(count1, "通り")

	// fmt.Println("\n--- 挿入ソートのデモ ---")
	// RunInsertionSortDemo()

	// fmt.Println("\n--- クイックソートのデモ ---")
	// RunQuickSortDemo()

	// fmt.Println("\n--- マージソートのデモ ---")
	// RunMergeSortDemo()

	// fmt.Println("\n--- ヒープソートのデモ ---")
	// RunHeapSortDemo()

	// fmt.Println("\n--- フィボナッチ数列のデモ ---")
	// RunFibonacciDemo()

	// fmt.Println("\n--- 階乗計算のデモ ---")
	// FactorialDemo()

	// fmt.Println("\n--- 再帰関数のデモ ---")
	// result := recursion.Factorial(5)
	// fmt.Println(result)

	// result2 := recursion.FactorialWithFor(5)
	// fmt.Println(result2)

	// result3 := recursion.CountDown(5)
	// fmt.Println(result3)

	// recursion.CountDownWithFor(5)

	// result4 := recursion.Sum(5)
	// fmt.Println(result4)

	// result5 := recursion.SumWithFor(5)
	// fmt.Println(result5)

	// arr := []int{1, 2, 3, 4, 5}
	// _ = recursion.PrintArray(arr, 0)
	// recursion.PrintArrayWithFor(arr)

	// fmt.Println(arr[:2])
	// arr2 := append(arr[:0], arr[3:]...)
	// fmt.Println(arr2)

	// 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144
	// result6 := recursion.Fibonacci(10)
	// fmt.Println(result6)

	// result7 := recursion.FibonacciWithFor(10)
	// fmt.Println(result7)

	// arr := []int{1, 2, 4, 10}
	// result8 := recursion.SumArray(arr, 0)
	// fmt.Println(result8)

	// result9 := recursion.SumArrayWithFor(arr)
	// fmt.Println(result9)

	// num := 5 / 2
	// fmt.Println(num)

	// arr2 := []int{1, 2, 3, 4, 5, 6}

	// for _, v := range arr2 {
	// 	result10 := recursion.BinarySearchWithFor(arr2, v)
	// 	fmt.Println(result10)
	// }

	// for _, v := range arr2 {
	// 	result10 := recursion.BinarySearchWithFor(arr2, v)
	// 	fmt.Println(result10)
	// }

	// recursion.Tree()

	// ハノイの塔のデモ
	// fmt.Println("\n--- ハノイの塔（再帰）のデモ ---")

	// // 3枚のディスクでシミュレーション
	// fmt.Println("■ 3枚のディスクでシミュレーション:")
	// recursion.SimulateTowerOfHanoi(3)

	// // ステップ表示版（再帰呼び出しの様子を表示）
	// fmt.Println("\n■ 再帰呼び出しの詳細（3枚）:")
	// recursion.TowerOfHanoiWithSteps(3, "A", "C", "B", 0)

	// // 移動回数の計算
	// fmt.Println("\n■ ディスク枚数と移動回数:")
	// for n := 1; n <= 10; n++ {
	// 	moves := recursion.CountHanoiMoves(n)
	// 	fmt.Printf("  %2d枚: %d回\n", n, moves)
	// }

	// recursion.TreeWithFor()

	// 配列の逆順（再帰）のデモ
	// fmt.Println("\n--- 配列の逆順（再帰）のデモ ---")
	// recursion.ReverseArrayDemo()

	// 最大公約数（ユークリッドの互除法）のデモ
	// fmt.Println("\n--- 最大公約数（ユークリッドの互除法）のデモ ---")
	// RunGCDDemo()

	// 最小公倍数（LCM）のデモ
	// fmt.Println("\n--- 最小公倍数（LCM）のデモ ---")
	// DemoLCM()

	// フィボナッチ数列（動的計画法）のデモ
	// fmt.Println("\n--- フィボナッチ数列（動的計画法）のデモ ---")
	// DemoFibonacciDP()

	// 部分和問題（動的計画法）のデモ
	// fmt.Println("\n--- 部分和問題（動的計画法）のデモ ---")
	// DemoSubsetSum()

	// 深さ優先探索（DFS）のデモ
	// fmt.Println("\n--- 深さ優先探索（DFS）のデモ ---")
	// DFSDemo()

	// スライディングウィンドウのデモ
	fmt.Println("\n--- スライディングウィンドウのデモ ---")
	DemoSlidingWindow()
}

// RunGCDDemo は最大公約数（ユークリッドの互除法）のデモを実行
func RunGCDDemo() {
	fmt.Println("\n1. 基本的な最大公約数の計算:")
	fmt.Println("=================================")

	// 基本例
	testCases := []struct {
		a, b int
	}{
		{48, 18},
		{100, 40},
		{1071, 462},
		{13, 17}, // 互いに素
	}

	for _, tc := range testCases {
		gcd := GCD(tc.a, tc.b)
		fmt.Printf("GCD(%d, %d) = %d\n", tc.a, tc.b, gcd)
	}

	fmt.Println("\n2. ステップごとの計算過程:")
	fmt.Println("=================================")
	GCDWithSteps(48, 18)

	fmt.Println("\n3. 拡張ユークリッドの互除法:")
	fmt.Println("=================================")
	gcd, x, y := ExtendedGCDWithSteps(48, 18)
	fmt.Printf("\n検証: 48 × %d + 18 × %d = %d\n", x, y, gcd)

	fmt.Println("\n4. 最小公倍数の計算:")
	fmt.Println("=================================")
	lcmTests := []struct {
		a, b int
	}{
		{12, 18},
		{15, 25},
		{7, 13},
	}

	for _, tc := range lcmTests {
		lcm := LCM(tc.a, tc.b)
		gcd := GCD(tc.a, tc.b)
		fmt.Printf("LCM(%d, %d) = %d (GCD = %d)\n", tc.a, tc.b, lcm, gcd)
	}

	fmt.Println("\n5. 複数の数の最大公約数:")
	fmt.Println("=================================")
	nums := []int{48, 64, 16, 32}
	gcdMultiple := GCDMultiple(nums...)
	fmt.Printf("GCD(%v) = %d\n", nums, gcdMultiple)

	nums2 := []int{12, 18, 24}
	gcdMultiple2 := GCDMultiple(nums2...)
	fmt.Printf("GCD(%v) = %d\n", nums2, gcdMultiple2)

	fmt.Println("\n6. 複数の数の最小公倍数:")
	fmt.Println("=================================")
	nums3 := []int{4, 6, 8}
	lcmMultiple := LCMMultiple(nums3...)
	fmt.Printf("LCM(%v) = %d\n", nums3, lcmMultiple)

	fmt.Println("\n7. 分数の約分:")
	fmt.Println("=================================")
	fractions := []struct {
		num, den int
	}{
		{24, 36},
		{15, 25},
		{100, 250},
	}

	for _, f := range fractions {
		num, den := SimplifyFraction(f.num, f.den)
		fmt.Printf("%d/%d = %d/%d\n", f.num, f.den, num, den)
	}

	fmt.Println("\n8. 互いに素の判定:")
	fmt.Println("=================================")
	primeTests := []struct {
		a, b int
	}{
		{15, 28},
		{13, 17},
		{24, 35},
		{7, 11},
	}

	for _, tc := range primeTests {
		if IsCoprime(tc.a, tc.b) {
			fmt.Printf("%d と %d は互いに素です\n", tc.a, tc.b)
		} else {
			fmt.Printf("%d と %d は互いに素ではありません (GCD = %d)\n", tc.a, tc.b, GCD(tc.a, tc.b))
		}
	}

	fmt.Println("\n9. 性能比較（再帰 vs 反復）:")
	fmt.Println("=================================")
	a, b := 1234567890, 987654321

	// 再帰版
	start := fmt.Sprintf("GCD (再帰): ")
	gcdRec := GCD(a, b)
	fmt.Printf("%s GCD(%d, %d) = %d\n", start, a, b, gcdRec)

	// 反復版
	start2 := fmt.Sprintf("GCD (反復): ")
	gcdIter := GCDIterative(a, b)
	fmt.Printf("%s GCD(%d, %d) = %d\n", start2, a, b, gcdIter)
}
