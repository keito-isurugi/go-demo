package algorithm

import "fmt"

// GCD はユークリッドの互除法で最大公約数を求める
func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

// GCDIterative は反復版のユークリッドの互除法
func GCDIterative(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCDWithSteps はステップごとの計算過程を表示
func GCDWithSteps(a, b int) int {
	fmt.Printf("GCD(%d, %d)の計算:\n", a, b)
	step := 1

	for b != 0 {
		remainder := a % b
		fmt.Printf("Step %d: %d = %d × %d + %d\n", step, a, b, a/b, remainder)
		a, b = b, remainder
		step++
	}

	fmt.Printf("\n最大公約数: %d\n", a)
	return a
}

// ExtendedGCD は拡張ユークリッドの互除法（ax + by = gcd(a,b)のx,yも求める）
func ExtendedGCD(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}

	gcd, x1, y1 := ExtendedGCD(b, a%b)
	x = y1
	y = x1 - (a/b)*y1

	return gcd, x, y
}

// ExtendedGCDWithSteps は拡張ユークリッドの互除法をステップ表示
func ExtendedGCDWithSteps(a, b int) (gcd, x, y int) {
	fmt.Printf("\n拡張ユークリッドの互除法: %dx + %dy = gcd(%d, %d)\n", a, b, a, b)
	fmt.Println("計算過程:")

	// 商を保存するスライス
	quotients := []int{}
	origA, origB := a, b

	// 通常のユークリッドの互除法で商を記録
	for b != 0 {
		q := a / b
		quotients = append(quotients, q)
		fmt.Printf("%d = %d × %d + %d\n", a, b, q, a%b)
		a, b = b, a%b
	}

	gcd = a
	fmt.Printf("\nGCD = %d\n", gcd)

	// 逆算してx, yを求める
	x, y = 1, 0
	for i := len(quotients) - 1; i >= 0; i-- {
		x, y = y, x-quotients[i]*y
	}

	// bが元々負の場合の調整
	if origB < 0 {
		y = -y
	}
	if origA < 0 {
		x = -x
	}

	fmt.Printf("\n結果: %d × %d + %d × %d = %d\n", origA, x, origB, y, gcd)
	return gcd, x, y
}

// LCM は最小公倍数を求める（GCDを利用）
func LCM(a, b int) int {
	return (a * b) / GCD(a, b)
}

// GCDMultiple は複数の数の最大公約数を求める
func GCDMultiple(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result = GCD(result, nums[i])
		if result == 1 {
			break // 1になったらそれ以上小さくならない
		}
	}

	return result
}

// LCMMultiple は複数の数の最小公倍数を求める
func LCMMultiple(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result = LCM(result, nums[i])
	}

	return result
}

// IsCoprime は2つの数が互いに素かどうかを判定
func IsCoprime(a, b int) bool {
	return GCD(a, b) == 1
}

// SimplifyFraction は分数を最簡分数に約分
func SimplifyFraction(numerator, denominator int) (int, int) {
	if denominator == 0 {
		return numerator, denominator
	}

	gcd := GCD(abs(numerator), abs(denominator))
	return numerator / gcd, denominator / gcd
}

// abs は絶対値を返す
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}