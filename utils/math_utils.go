package utils

// Add は2つの整数を足し算する
func Add(a, b int) int {
	return a + b
}

// Subtract は2つの整数を引き算する
func Subtract(a, b int) int {
	return a - b
}

// Multiply は2つの整数を掛け算する
func Multiply(a, b int) int {
	return a * b
}

// Divide は2つの整数を割り算する（ゼロ除算チェック付き）
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// IsEven は偶数かどうかを判定する
func IsEven(n int) bool {
	return n%2 == 0
}

// IsPositive は正の数かどうかを判定する
func IsPositive(n int) bool {
	return n > 0
}

// Max は2つの整数の大きい方を返す
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min は2つの整数の小さい方を返す
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
