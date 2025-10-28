## 概要

- Goで部分和問題を動的計画法で実装
- 部分和問題は与えられた数値の集合から、合計が目標値になる部分集合が存在するかを判定する問題。動的計画法の典型的な応用例の一つ。

### 特徴

- **効率的計算**: 時間計算量O(n×target)で解ける
- **メモリトレードオフ**: 空間O(n×target)または最適化でO(target)
- **2つのアプローチ**: 判定問題と個数カウント
- **バックトラック**: 実際の解（部分集合）も復元可能
- **実用性**: ナップサック問題の基礎となる重要アルゴリズム

## 部分和問題とは

### 問題の定義

**入力**: 整数の配列 nums = [n₁, n₂, ..., nₙ] と目標値 target

**出力**: nums の部分集合で合計が target になるものが存在するか

### 具体例

```
配列: [3, 34, 4, 12, 5, 2]
目標値: 9

解: [4, 5] または [3, 4, 2]
→ 4 + 5 = 9 ✓
→ 3 + 4 + 2 = 9 ✓
```

### 応用例

- ナップサック問題
- 予算配分問題
- 硬貨での支払い問題
- パーティション問題（配列を等しい和に分割）
- 組み合わせ最適化

## 動的計画法のアプローチ

### 基本的な考え方

**部分問題の定義**:
```
dp[i][j] = nums[0..i-1] の要素を使って合計 j を作れるか（true/false）
```

**漸化式**:
```
dp[i][j] = dp[i-1][j]  OR  dp[i-1][j - nums[i-1]]
           ↑                ↑
      使わない場合      使う場合
```

**初期条件**:
```
dp[i][0] = true  （合計0は常に作れる：何も選ばない）
dp[0][j] = false （j > 0の場合、要素がないので作れない）
```

### 計算過程の可視化

配列 `[3, 4, 5, 2]` で目標値 `9` の例：

```
DPテーブル (◯: 作れる, ×: 作れない):
       0   1   2   3   4   5   6   7   8   9
初期   ◯   ×   ×   ×   ×   ×   ×   ×   ×   ×
  3    ◯   ×   ×   ◯   ×   ×   ×   ×   ×   ×
  4    ◯   ×   ×   ◯   ◯   ×   ×   ◯   ×   ×
  5    ◯   ×   ×   ◯   ◯   ◯   ×   ◯   ◯   ◯
  2    ◯   ×   ◯   ◯   ◯   ◯   ◯   ◯   ◯   ◯

結果: dp[4][9] = ◯ → 合計9は作れる ✓
```

#### ステップバイステップ

**Step 1: 初期状態**
- 合計0は常に作れる（何も選ばない）

**Step 2: 要素3を追加**
- 3が作れるようになる

**Step 3: 要素4を追加**
- 4が作れる（4単体）
- 7が作れる（3 + 4）

**Step 4: 要素5を追加**
- 5が作れる（5単体）
- 8が作れる（3 + 5）
- 9が作れる（4 + 5） ← 目標達成！

**Step 5: 要素2を追加**
- さらに多くの合計が作れる

## サンプルコード

### 基本実装（判定問題）

```go
package main

import "fmt"

// SubsetSum は動的計画法で部分和問題を解く
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
```

### 実際の要素を見つける（バックトラック）

```go
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
```

### 解の個数を数える

```go
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
```

### 空間最適化版（O(target)）

```go
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
```

**重要**: 1次元配列で実装する場合、**逆順**で更新する必要がある。順方向だと同じ要素を複数回使ってしまう。

```go
// ✗ NG: 順方向 → 同じ要素を複数回使う
for j := num; j <= target; j++ {
    dp[j] = dp[j] || dp[j-num]
}

// ✓ OK: 逆順 → 各要素を1回だけ使う
for j := target; j >= num; j-- {
    dp[j] = dp[j] || dp[j-num]
}
```

## 計算量

### 時間計算量

| 実装方法 | 時間計算量 | 説明 |
| --- | --- | --- |
| 全探索（ビット演算） | O(2^n) | すべての部分集合を調べる |
| 2次元DP | O(n × target) | n個の要素×targetまでの値 |
| 1次元DP（最適化） | O(n × target) | 同上 |

### 空間計算量

| 実装方法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 全探索 | O(n) | 再帰スタック |
| 2次元DP | O(n × target) | dpテーブル |
| 1次元DP（最適化） | O(target) | 1次元配列のみ |

### パフォーマンス比較

```
配列サイズ n=20, target=1000 の場合:

全探索:       約1秒（2^20 = 100万回の計算）
2次元DP:      約0.001秒（20×1000 = 2万回の計算）
最適化DP:     約0.001秒（同上、メモリは1/20）
```

## 動作原理の詳細

### なぜ漸化式が成り立つか

`dp[i][j]` を計算する際、`nums[i-1]` を使うか使わないかの2択：

**1. 使わない場合**
- `nums[0..i-2]` で合計 `j` を作れればOK
- → `dp[i][j] = dp[i-1][j]`

**2. 使う場合**
- `nums[0..i-2]` で合計 `j - nums[i-1]` を作れればOK
- そこに `nums[i-1]` を足せば `j` になる
- → `dp[i][j] = dp[i-1][j - nums[i-1]]`

どちらか一方でも可能なら `dp[i][j] = true`

### バックトラックの仕組み

dpテーブル構築後、逆順に辿って実際の解を復元：

```
配列: [3, 4, 5, 2], 目標: 9

dpテーブルから逆順に辿る:
i=4, j=9: nums[3]=2 を使ったか？
  → dp[3][9]=true なので使わない

i=3, j=9: nums[2]=5 を使ったか？
  → dp[2][9]=false, dp[2][4]=true なので使う！
  → j = 9 - 5 = 4

i=2, j=4: nums[1]=4 を使ったか？
  → dp[1][4]=false, dp[1][0]=true なので使う！
  → j = 4 - 4 = 0

i=1, j=0: 終了

解: [5, 4]
```

## 使いどころ

### 向いてる場面

- **ナップサック問題**: 重さと価値のある品物を選ぶ
- **予算配分**: 限られた予算で目標額を達成
- **組み合わせ最適化**: 複数の選択肢から最適な組み合わせ
- **パーティション問題**: 配列を等しい和に分割できるか

### 実例1: 硬貨での支払い

```go
func canPayExactly(coins []int, price int) bool {
    return SubsetSum(coins, price)
}

// 使用例
coins := []int{1, 5, 10, 25, 50, 100}
price := 87

if canPayExactly(coins, price) {
    subset := FindSubsetSum(coins, price)
    fmt.Printf("87円の支払い: %v円硬貨\n", subset)
    // 例: [50, 25, 10, 1, 1] など
}
```

### 実例2: 配列の分割（Equal Partition）

配列を等しい和の2グループに分けられるか：

```go
// CanPartition は配列を等しい和に分割できるか判定
func CanPartition(nums []int) bool {
    sum := 0
    for _, num := range nums {
        sum += num
    }

    // 合計が奇数なら分割不可能
    if sum%2 != 0 {
        return false
    }

    // 半分の値を作れるか？
    return SubsetSum(nums, sum/2)
}

// 使用例
nums := []int{1, 5, 11, 5}
// 合計 = 22, 半分 = 11
// [1, 5, 5] と [11] に分割可能 → true
```

### 実例3: ターゲットサム問題

配列の各要素に+または-を付けて目標値を作る：

```go
// FindTargetSumWays は目標値を作る方法の数
func FindTargetSumWays(nums []int, target int) int {
    sum := 0
    for _, num := range nums {
        sum += num
    }

    // 不可能な場合
    if sum < target || (sum+target)%2 != 0 {
        return 0
    }

    // P - N = target かつ P + N = sum
    // → P = (sum + target) / 2
    // Pの部分集合を作る方法の数を求める
    return CountSubsetSums(nums, (sum+target)/2)
}
```

## 関連する問題

### 1. 0/1ナップサック問題

部分和問題の発展版（価値の概念を追加）：

```go
// Knapsack は0/1ナップサック問題を解く
func Knapsack(weights []int, values []int, capacity int) int {
    n := len(weights)
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, capacity+1)
    }

    for i := 1; i <= n; i++ {
        for w := 0; w <= capacity; w++ {
            // 入れない場合
            dp[i][w] = dp[i-1][w]

            // 入れる場合
            if w >= weights[i-1] {
                value := dp[i-1][w-weights[i-1]] + values[i-1]
                if value > dp[i][w] {
                    dp[i][w] = value
                }
            }
        }
    }

    return dp[n][capacity]
}
```

### 2. 完全部分和問題（無制限使用可）

各要素を何回でも使える場合：

```go
// SubsetSumUnlimited は要素を無制限に使える部分和問題
func SubsetSumUnlimited(nums []int, target int) bool {
    dp := make([]bool, target+1)
    dp[0] = true

    // 順方向で更新（同じ要素を複数回使えるため）
    for j := 1; j <= target; j++ {
        for _, num := range nums {
            if j >= num && dp[j-num] {
                dp[j] = true
                break
            }
        }
    }

    return dp[target]
}
```

### 3. 最小要素数の部分和

目標値を作る最小の要素数：

```go
// MinSubsetSum は目標値を作る最小要素数
func MinSubsetSum(nums []int, target int) int {
    dp := make([]int, target+1)
    for i := range dp {
        dp[i] = target + 1 // 十分大きい値
    }
    dp[0] = 0

    for _, num := range nums {
        for j := target; j >= num; j-- {
            dp[j] = min(dp[j], dp[j-num]+1)
        }
    }

    if dp[target] > target {
        return -1 // 不可能
    }
    return dp[target]
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

## 最適化テクニック

### 1. 早期終了

目標値に到達したら計算を打ち切る：

```go
func SubsetSumEarlyExit(nums []int, target int) bool {
    dp := make([]bool, target+1)
    dp[0] = true

    for _, num := range nums {
        for j := target; j >= num; j-- {
            if dp[j-num] {
                dp[j] = true
                if j == target {
                    return true // 早期終了
                }
            }
        }
    }

    return dp[target]
}
```

### 2. 事前ソート

配列をソートして枝刈り：

```go
import "sort"

func SubsetSumPruned(nums []int, target int) bool {
    sort.Ints(nums)

    // 合計が目標値未満なら不可能
    sum := 0
    for _, num := range nums {
        sum += num
    }
    if sum < target {
        return false
    }

    // 最大値が目標値を超える要素は除外
    validNums := []int{}
    for _, num := range nums {
        if num <= target {
            validNums = append(validNums, num)
        }
    }

    return SubsetSum(validNums, target)
}
```

### 3. ビットセット最適化

ビット演算で高速化（targetが小さい場合）：

```go
// SubsetSumBitset はビットセットを使った高速化版
func SubsetSumBitset(nums []int, target int) bool {
    // dp を uint64 のビットで表現（target < 64の場合）
    if target >= 64 {
        return SubsetSum(nums, target)
    }

    var dp uint64 = 1 // dp[0] = true

    for _, num := range nums {
        if num < 64 {
            dp |= dp << uint(num)
        }
    }

    return (dp & (1 << uint(target))) != 0
}
```

## デバッグとテスト

### テストケース

```go
func TestSubsetSum() {
    testCases := []struct {
        nums   []int
        target int
        want   bool
    }{
        {[]int{3, 34, 4, 12, 5, 2}, 9, true},     // [4, 5]
        {[]int{3, 34, 4, 12, 5, 2}, 30, false},   // 不可能
        {[]int{1, 2, 3, 7}, 6, true},             // [1, 2, 3]
        {[]int{1, 2, 7, 1, 5}, 10, true},         // [1, 2, 7]
        {[]int{1, 5, 11, 5}, 11, true},           // [1, 5, 5]
        {[]int{2, 4, 6}, 5, false},               // 奇数は作れない
    }

    for _, tc := range testCases {
        got := SubsetSum(tc.nums, tc.target)
        if got != tc.want {
            fmt.Printf("FAIL: nums=%v, target=%d, got=%v, want=%v\n",
                tc.nums, tc.target, got, tc.want)
        }
    }
}
```

### デバッグ用可視化

```go
// PrintDPTable はdpテーブルを見やすく表示
func PrintDPTable(dp [][]bool, nums []int, target int) {
    fmt.Print("     ")
    for j := 0; j <= target; j++ {
        fmt.Printf("%3d", j)
    }
    fmt.Println()

    for i := 0; i < len(dp); i++ {
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
}
```

## まとめ

### メリット

- 時間計算量O(n×target)で効率的
- 解の存在判定だけでなく、実際の解や解の個数も求められる
- 空間を最適化してO(target)に削減可能
- ナップサック問題の基礎となる重要アルゴリズム
- 実用的な応用が多い

### デメリット

- targetが大きいとメモリと時間が問題に
- 浮動小数点数には適用できない
- targetが非常に大きい場合は別のアプローチが必要

### 使うべき時

- **判定問題**: 特定の合計を作れるか
- **組み合わせ最適化**: 複数の選択肢から最適な組み合わせ
- **ナップサック系問題**: 容量制約のある選択問題
- **パーティション問題**: 配列の分割

### targetの大きさによる選択

| targetの範囲 | 推奨実装 | 理由 |
| --- | --- | --- |
| target < 1000 | 1次元DP | 十分高速かつ省メモリ |
| 1000 ≤ target < 10^5 | 1次元DP | メモリ効率重視 |
| target ≥ 10^5 | 別手法検討 | Meet-in-the-middle など |
| n < 20 | ビット全探索も可 | 実装がシンプル |

部分和問題は動的計画法の典型問題であり、ナップサック問題や組み合わせ最適化の基礎となる。dpテーブルの構築、漸化式の立て方、バックトラックによる解の復元など、DPの重要な概念を学ぶのに最適なアルゴリズム。
