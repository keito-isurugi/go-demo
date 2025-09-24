## 概要

- Goでフィボナッチ数列を動的計画法で実装
- 動的計画法（DP）は部分問題の解を保存し再利用することで、再帰の指数的計算量を線形に改善する手法。フィボナッチ数列はDPの典型例。

### 特徴

- **効率的計算**: 時間計算量O(n)で実用的
- **メモリトレードオフ**: 空間をO(n)使って時間を節約
- **2つのアプローチ**: トップダウン（メモ化）とボトムアップ
- **最適化可能**: 空間計算量をO(1)に削減可能
- **DPの典型例**: 動的計画法の学習に最適

## 動的計画法とは

動的計画法の基本原理：

1. **部分問題への分割**: 大きな問題を小さな部分問題に分解
2. **重複の排除**: 同じ部分問題を何度も解かない
3. **結果の保存**: 計算結果をテーブルやメモに保存
4. **再利用**: 保存した結果を使って効率的に解を求める

### 2つのアプローチ

#### トップダウン（メモ化）
- 再帰的に計算しながら結果をメモ
- 必要な部分だけ計算
- 実装が直感的

#### ボトムアップ
- 小さい問題から順に解いていく
- 全ての部分問題を計算
- 反復で実装、スタックオーバーフローなし

## サンプルコード

### ボトムアップ動的計画法

```go
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

```

### ステップ表示版

```go
// FibonacciDPWithSteps はステップごとに計算過程を表示
func FibonacciDPWithSteps(n int) int {
	fmt.Printf("\\nFibonacci(%d) を動的計画法で計算:\\n", n)
	fmt.Println("=================================")

	if n <= 1 {
		fmt.Printf("n <= 1 のため、結果は %d\\n", n)
		return n
	}

	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1

	fmt.Printf("初期値: dp[0] = %d, dp[1] = %d\\n", dp[0], dp[1])
	fmt.Println("\\n計算過程:")

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
		fmt.Printf("dp[%d] = dp[%d] + dp[%d] = %d + %d = %d\\n",
			i, i-1, i-2, dp[i-1], dp[i-2], dp[i])
	}

	fmt.Printf("\\n最終的なdpテーブル: %v\\n", dp)
	fmt.Printf("結果: Fibonacci(%d) = %d\\n", n, dp[n])

	return dp[n]
}

```

### 空間最適化版（O(1)）

```go
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

```

### 数列生成版

```go
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

```

## 動作原理

### ボトムアップDPの計算過程

`Fibonacci(8)` の計算：

```
初期値:
dp[0] = 0
dp[1] = 1

計算過程:
dp[2] = dp[1] + dp[0] = 1 + 0 = 1
dp[3] = dp[2] + dp[1] = 1 + 1 = 2
dp[4] = dp[3] + dp[2] = 2 + 1 = 3
dp[5] = dp[4] + dp[3] = 3 + 2 = 5
dp[6] = dp[5] + dp[4] = 5 + 3 = 8
dp[7] = dp[6] + dp[5] = 8 + 5 = 13
dp[8] = dp[7] + dp[6] = 13 + 8 = 21

dpテーブル: [0, 1, 1, 2, 3, 5, 8, 13, 21]
結果: Fibonacci(8) = 21
```

### トップダウン（メモ化）との比較

#### トップダウン（メモ化）
```
Fib(5) の計算:
├─ Fib(4) を計算
│  ├─ Fib(3) を計算
│  │  ├─ Fib(2) を計算
│  │  │  ├─ Fib(1) = 1
│  │  │  └─ Fib(0) = 0
│  │  └─ Fib(1) = 1 (既に計算済み)
│  └─ Fib(2) = 1 (メモから取得)
└─ Fib(3) = 2 (メモから取得)
```

#### ボトムアップ（DP）
```
Fib(5) の計算:
dp[0] = 0
dp[1] = 1
dp[2] = dp[1] + dp[0] = 1
dp[3] = dp[2] + dp[1] = 2
dp[4] = dp[3] + dp[2] = 3
dp[5] = dp[4] + dp[3] = 5
```

## 計算量

### 時間計算量

| 実装方法 | 時間計算量 | 説明 |
| --- | --- | --- |
| 純粋再帰 | O(2^n) | 指数的増加、重複計算多数 |
| メモ化 | O(n) | 各値を1回だけ計算 |
| ボトムアップDP | O(n) | 0からnまで線形に計算 |
| 最適化DP | O(n) | 同上 |

### 空間計算量

| 実装方法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 純粋再帰 | O(n) | 再帰スタックの深さ |
| メモ化 | O(n) | メモテーブル + スタック |
| ボトムアップDP | O(n) | dpテーブル |
| 最適化DP | O(1) | 2つの変数のみ |

### 実行時間比較（n=40）

```
純粋再帰:    約2秒
メモ化:      0.00001秒
ボトムアップDP: 0.00001秒
最適化DP:    0.00001秒
```

## 使いどころ

### 向いてる場面

- 実用的なフィボナッチ数の計算
- 大きなn（n > 30）でも高速に計算
- 動的計画法の学習
- 最適化の比較（空間vs時間）

### 実例：フィボナッチ数列の生成

```go
func main() {
	// 最初の20項を生成
	sequence := FibonacciSequenceDP(20)
	fmt.Printf("フィボナッチ数列: %v\\n", sequence)
	// 出力: [0 1 1 2 3 5 8 13 21 34 55 89 144 233 377 610 987 1597 2584 4181 6765]

	// 大きな値でも高速
	fib100 := FibonacciDPOptimized(100)
	fmt.Printf("Fibonacci(100) = %d\\n", fib100)
}

```

## 実装の比較

| アプローチ | 時間 | 空間 | 実装の複雑さ | スタックオーバーフロー |
| --- | --- | --- | --- | --- |
| 純粋再帰 | O(2^n) | O(n) | 簡単 | あり |
| メモ化 | O(n) | O(n) | やや複雑 | あり（大きなn） |
| ボトムアップDP | O(n) | O(n) | やや複雑 | なし |
| 最適化DP | O(n) | O(1) | やや複雑 | なし |

### コード行数の比較

```go
// 純粋再帰: 4行
func Fib(n int) int {
	if n <= 1 { return n }
	return Fib(n-1) + Fib(n-2)
}

// ボトムアップDP: 8行
func FibDP(n int) int {
	if n <= 1 { return n }
	dp := make([]int, n+1)
	dp[0], dp[1] = 0, 1
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

// 最適化DP: 6行
func FibDPOpt(n int) int {
	if n <= 1 { return n }
	prev2, prev1 := 0, 1
	for i := 2; i <= n; i++ {
		prev2, prev1 = prev1, prev2+prev1
	}
	return prev1
}
```

## 最適化テクニック

### 1. 空間計算量の削減

dpテーブル全体ではなく、直前の2つの値のみ保持：

```go
// O(n) 空間
dp := make([]int, n+1)
dp[i] = dp[i-1] + dp[i-2]

// O(1) 空間
prev2, prev1 := 0, 1
current := prev1 + prev2
prev2, prev1 = prev1, current
```

### 2. 行列累乗法（O(log n)）

さらに高速化したい場合：

```go
// 行列[[1,1],[1,0]]のn乗を計算
// 時間計算量: O(log n)
// 空間計算量: O(1)
type Matrix [2][2]int

func matrixMultiply(a, b Matrix) Matrix {
	return Matrix{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
		{a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
	}
}

func matrixPower(mat Matrix, n int) Matrix {
	if n == 1 {
		return mat
	}
	if n%2 == 0 {
		half := matrixPower(mat, n/2)
		return matrixMultiply(half, half)
	}
	return matrixMultiply(mat, matrixPower(mat, n-1))
}

func FibonacciMatrix(n int) int {
	if n <= 1 {
		return n
	}
	base := Matrix{{1, 1}, {1, 0}}
	result := matrixPower(base, n)
	return result[0][1]
}
```

### 3. 並列計算

大量のフィボナッチ数を並列生成：

```go
func FibonacciConcurrent(n int) []int {
	results := make([]int, n+1)
	var wg sync.WaitGroup

	// 各値を並列計算（実際はDPの方が速い）
	for i := 0; i <= n; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = FibonacciDPOptimized(index)
		}(i)
	}

	wg.Wait()
	return results
}
```

## パフォーマンス測定

```go
import "time"

func benchmark(name string, n int, fn func(int) int) {
	start := time.Now()
	result := fn(n)
	duration := time.Since(start)
	fmt.Printf("%s: Fib(%d) = %d, 実行時間: %v\\n", name, n, result, duration)
}

func compareMethods() {
	n := 40

	// 純粋再帰は遅いのでn=40でも数秒かかる
	benchmark("純粋再帰", n, FibonacciRecursive)

	// DPは高速
	benchmark("メモ化", n, FibonacciMemoization)
	benchmark("ボトムアップDP", n, FibonacciDP)
	benchmark("最適化DP", n, FibonacciDPOptimized)

	// 大きなnでも高速
	benchmark("最適化DP(90)", 90, FibonacciDPOptimized)
}
```

## まとめ

### メリット

- 時間計算量O(n)で実用的
- 大きなnでも高速に計算可能
- スタックオーバーフローのリスクなし
- 空間をO(1)に最適化可能

### デメリット

- 純粋再帰より実装が複雑
- dpテーブル版は空間O(n)必要
- 小さなn（< 10）では再帰の方がシンプル

### 使うべき時

- 実用的なフィボナッチ数の計算
- n > 30 の場合
- 動的計画法の学習
- パフォーマンスが重要な場合

### 選択基準

| n の範囲 | 推奨実装 | 理由 |
| --- | --- | --- |
| n < 10 | 純粋再帰 | シンプルで十分速い |
| 10 ≤ n < 30 | メモ化 or DP | どちらでも高速 |
| n ≥ 30 | 最適化DP | 最速かつ省メモリ |
| n ≥ 10^6 | 行列累乗 | O(log n)で超高速 |

動的計画法はフィボナッチ数列の実用的な計算方法であり、DPの基本概念を学ぶ優れた例。ボトムアップとトップダウンの違い、時間と空間のトレードオフを理解するのに最適。
