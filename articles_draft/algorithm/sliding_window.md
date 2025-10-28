## 概要

- Goでスライディングウィンドウ（Sliding Window）を実装
- スライディングウィンドウは配列や文字列の連続する部分列を効率的に処理する技法。固定長または可変長のウィンドウを「スライド」させながら最適解を求める。

### 特徴

- **効率的処理**: 時間計算量O(n)で線形探索可能
- **重複計算の削減**: 前回の計算結果を再利用
- **2つのアプローチ**: 固定長ウィンドウと可変長ウィンドウ
- **幅広い応用**: 配列・文字列の部分列問題に適用可能
- **Two Pointersとの関連**: 左右2つのポインタで範囲を管理

## スライディングウィンドウとは

### 基本的な考え方

配列やリストの連続する要素を「窓（ウィンドウ）」として扱い、その窓を少しずつ移動（スライド）させながら処理する手法。

**イメージ図**:
```
配列: [2, 1, 5, 1, 3, 2]
ウィンドウサイズ k=3

Step 1: [2, 1, 5] 1, 3, 2  → 合計 = 8
Step 2:  2,[1, 5, 1] 3, 2  → 合計 = 7
Step 3:  2, 1,[5, 1, 3] 2  → 合計 = 9 ← 最大
Step 4:  2, 1, 5,[1, 3, 2] → 合計 = 6
```

### なぜ効率的か？

**総当たり（Brute Force）**:
- 各ウィンドウごとにゼロから合計を計算
- 時間計算量: O(n × k)

```go
// 総当たり: 毎回k個の要素を足す
for i := 0; i <= n-k; i++ {
    sum := 0
    for j := i; j < i+k; j++ {
        sum += arr[j]  // k回の加算
    }
}
```

**スライディングウィンドウ**:
- 前回の合計から1つ引いて1つ足すだけ
- 時間計算量: O(n)

```go
// スライディングウィンドウ: 1回の引き算と1回の足し算
windowSum = windowSum + arr[i] - arr[i-k]
```

**性能差**:
```
配列サイズ n=100,000, ウィンドウサイズ k=1,000
総当たり:        約1億回の計算 (100,000 × 1,000)
スライディング:  約10万回の計算 (100,000)
→ 1000倍高速！
```

## サンプルコード

### 1. 固定長ウィンドウ - 最大合計

最も基本的なパターン。サイズkの連続する部分配列で最大の合計を求める。

```go
package main

import "fmt"

// MaxSumSubarray は固定長kの部分配列の最大合計を求める
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
```

### 2. 可変長ウィンドウ - 最小長

合計が目標値以上になる最小の長さを求める。ウィンドウサイズが動的に変化。

```go
// MinSubarrayLen は合計がtarget以上になる最小長の部分配列を求める
func MinSubarrayLen(arr []int, target int) int {
	n := len(arr)
	minLen := n + 1
	windowSum := 0
	left := 0

	for right := 0; right < n; right++ {
		// 右端を拡大
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
```

### 3. 文字列 - 最長部分文字列

最大k種類の異なる文字を含む最長部分文字列の長さを求める。

```go
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
```

### 4. スライディングウィンドウの最大値

各ウィンドウ位置での最大値を求める（基本実装）。

```go
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
```

## 動作原理

### 固定長ウィンドウの動き

配列 `[2, 1, 5, 1, 3, 2]`、ウィンドウサイズ `k=3` の例：

```
初期ウィンドウ [0:3]:
arr[0] = 2, 累積和 = 2
arr[1] = 1, 累積和 = 3
arr[2] = 5, 累積和 = 8
→ 初期合計 = 8

ウィンドウをスライド:
[1:4] 削除: arr[0]=2, 追加: arr[3]=1 → 合計 = 8 - 2 + 1 = 7
[2:5] 削除: arr[1]=1, 追加: arr[4]=3 → 合計 = 7 - 1 + 3 = 9 ← 最大値更新!
[3:6] 削除: arr[2]=5, 追加: arr[5]=2 → 合計 = 9 - 5 + 2 = 6

結果: 最大合計 = 9 (位置 [2:5])
```

### 可変長ウィンドウの動き

配列 `[2, 3, 1, 2, 4, 3]`、目標値 `target=7` の例：

```
右端を拡大しながら条件を満たしたら左端を縮める:

right=0: [2]           合計=2
right=1: [2,3]         合計=5
right=2: [2,3,1]       合計=6
right=3: [2,3,1,2]     合計=8 ✓ 目標達成
  → 左端縮小: [3,1,2] 合計=6
right=4: [3,1,2,4]     合計=10 ✓
  → 左端縮小: [1,2,4] 合計=7 ✓ 最小長=3
  → 左端縮小: [2,4]   合計=6
right=5: [2,4,3]       合計=9 ✓
  → 左端縮小: [4,3]   合計=7 ✓
  → 左端縮小: [3]     合計=3

結果: 最小長 = 2 (部分配列 [4,3])
```

## 2つのアプローチ

### 固定長ウィンドウ

ウィンドウサイズが固定されている場合。

**特徴**:
- ウィンドウサイズ k は一定
- 右端が1つ進むと同時に左端も1つ進む
- 実装がシンプル

**典型的な問題**:
- サイズkの部分配列の最大/最小合計
- サイズkの部分配列の平均値
- 各ウィンドウ位置での最大/最小値

**実装パターン**:
```go
// 1. 初期ウィンドウを計算
windowSum := 0
for i := 0; i < k; i++ {
    windowSum += arr[i]
}

// 2. スライドさせながら更新
for i := k; i < n; i++ {
    windowSum = windowSum + arr[i] - arr[i-k]
    // 最適値を更新
}
```

### 可変長ウィンドウ

条件に応じてウィンドウサイズが変化する場合。

**特徴**:
- 左端(left)と右端(right)を独立に管理
- Two Pointers技法を使用
- 条件に応じてウィンドウを拡大・縮小

**典型的な問題**:
- 合計が条件を満たす最小/最大長の部分配列
- k種類の文字を含む最長部分文字列
- 重複のない最長部分文字列

**実装パターン**:
```go
left := 0
for right := 0; right < n; right++ {
    // 右端を拡大
    windowに arr[right] を追加

    // 条件を満たさない間、左端を縮める
    for 条件を満たさない {
        windowから arr[left] を削除
        left++
    }

    // 最適値を更新
}
```

## 計算量

### 時間計算量

| 実装方法 | 時間計算量 | 説明 |
| --- | --- | --- |
| 総当たり（固定長） | O(n × k) | 各位置でk個の要素を計算 |
| スライディングウィンドウ（固定長） | O(n) | 各要素を1回だけ処理 |
| 総当たり（可変長） | O(n²) | すべての部分配列を調べる |
| スライディングウィンドウ（可変長） | O(n) | 各要素を最大2回処理 |

### 空間計算量

| 実装方法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 基本実装 | O(1) | 数値の合計のみ保持 |
| 文字カウント | O(k) | k種類の文字を保持 |
| デック使用 | O(k) | ウィンドウ内の要素を保持 |

### なぜO(n)なのか？

可変長ウィンドウで`while`ループがネストしているのにO(n)になる理由：

```go
for right := 0; right < n; right++ {      // n回
    windowSum += arr[right]

    for windowSum >= target {             // 最悪n回？
        windowSum -= arr[left]
        left++
    }
}
```

**重要なポイント**:
- `left`は常に増加のみ（0 → n-1）
- 内側のループの合計実行回数は最大n回
- したがって全体でO(n + n) = O(n)

```
各要素について:
- right で1回追加される
- left で最大1回削除される
→ 各要素は最大2回処理されるだけ
→ 時間計算量はO(2n) = O(n)
```

## 使いどころ

### 向いている場面

1. **連続する部分配列/部分文字列の問題**
   - k個の連続する要素の最大/最小合計
   - 条件を満たす最小/最大長の部分配列

2. **ストリーミングデータの処理**
   - 直近k個のデータの統計情報
   - 移動平均の計算

3. **文字列の部分文字列問題**
   - 特定の条件を満たす最長/最短部分文字列
   - 重複のない最長部分文字列

### 実例1: 株価の移動平均

直近k日間の株価の平均を求める。

```go
// MovingAverage は株価の移動平均を計算
func MovingAverage(prices []int, k int) []float64 {
	n := len(prices)
	if n < k {
		return []float64{}
	}

	averages := make([]float64, 0, n-k+1)

	// 初期ウィンドウの合計
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += prices[i]
	}
	averages = append(averages, float64(windowSum)/float64(k))

	// スライド
	for i := k; i < n; i++ {
		windowSum = windowSum + prices[i] - prices[i-k]
		averages = append(averages, float64(windowSum)/float64(k))
	}

	return averages
}

// 使用例
prices := []int{100, 102, 101, 105, 110, 108, 107}
ma := MovingAverage(prices, 3)
// [101, 102.67, 105.33, 107.67, 108.33]
```

### 実例2: ログファイルの異常検知

直近k秒間のエラー数が閾値を超えたら警告。

```go
// DetectAnomaly はエラーログの異常を検知
func DetectAnomaly(errorCounts []int, k int, threshold int) []int {
	anomalyIndices := []int{}
	n := len(errorCounts)

	if n < k {
		return anomalyIndices
	}

	// 初期ウィンドウ
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += errorCounts[i]
	}
	if windowSum > threshold {
		anomalyIndices = append(anomalyIndices, 0)
	}

	// スライド
	for i := k; i < n; i++ {
		windowSum = windowSum + errorCounts[i] - errorCounts[i-k]
		if windowSum > threshold {
			anomalyIndices = append(anomalyIndices, i-k+1)
		}
	}

	return anomalyIndices
}
```

### 実例3: DNAシーケンスの解析

特定のパターンを含む最短のDNA配列を見つける。

```go
// ShortestDNASequence は4種類の塩基(A,C,G,T)をすべて含む最短部分列の長さ
func ShortestDNASequence(dna string) int {
	required := map[byte]bool{'A': true, 'C': true, 'G': true, 'T': true}
	charCount := make(map[byte]int)
	minLen := len(dna) + 1
	left := 0
	found := 0

	for right := 0; right < len(dna); right++ {
		char := dna[right]
		if required[char] {
			if charCount[char] == 0 {
				found++
			}
			charCount[char]++
		}

		// 4種類すべて含む間、左端を縮める
		for found == 4 {
			currentLen := right - left + 1
			if currentLen < minLen {
				minLen = currentLen
			}

			leftChar := dna[left]
			if required[leftChar] {
				charCount[leftChar]--
				if charCount[leftChar] == 0 {
					found--
				}
			}
			left++
		}
	}

	if minLen == len(dna)+1 {
		return 0
	}
	return minLen
}
```

## 関連する問題パターン

### 1. 重複のない最長部分文字列

```go
// LengthOfLongestSubstring は重複のない最長部分文字列の長さ
func LengthOfLongestSubstring(s string) int {
	charIndex := make(map[byte]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		// 重複が見つかったら左端を移動
		if index, exists := charIndex[s[right]]; exists && index >= left {
			left = index + 1
		}

		charIndex[s[right]] = right
		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
		}
	}

	return maxLen
}

// 例: "abcabcbb" → 3 ("abc")
// 例: "bbbbb" → 1 ("b")
// 例: "pwwkew" → 3 ("wke")
```

### 2. アナグラムの検出

文字列s1のアナグラムがs2に含まれるかを判定。

```go
// FindAnagrams は文字列sに含まれるpのアナグラムの開始位置をすべて見つける
func FindAnagrams(s string, p string) []int {
	result := []int{}
	if len(s) < len(p) {
		return result
	}

	// pの文字カウント
	pCount := make(map[byte]int)
	for i := 0; i < len(p); i++ {
		pCount[p[i]]++
	}

	// ウィンドウの文字カウント
	windowCount := make(map[byte]int)
	k := len(p)

	// 初期ウィンドウ
	for i := 0; i < k; i++ {
		windowCount[s[i]]++
	}
	if mapsEqual(pCount, windowCount) {
		result = append(result, 0)
	}

	// スライド
	for i := k; i < len(s); i++ {
		windowCount[s[i]]++
		windowCount[s[i-k]]--
		if windowCount[s[i-k]] == 0 {
			delete(windowCount, s[i-k])
		}

		if mapsEqual(pCount, windowCount) {
			result = append(result, i-k+1)
		}
	}

	return result
}

func mapsEqual(m1, m2 map[byte]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}
```

### 3. 最大連続1の個数（k回のフリップ許可）

0をk回まで1に変更できる時、最長の連続1の長さを求める。

```go
// LongestOnes は最大k回のフリップで作れる最長の連続1の長さ
func LongestOnes(nums []int, k int) int {
	maxLen := 0
	left := 0
	zeroCount := 0

	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			zeroCount++
		}

		// 0の個数がkを超えたら左端を縮める
		for zeroCount > k {
			if nums[left] == 0 {
				zeroCount--
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

// 例: nums=[1,1,1,0,0,0,1,1,1,1,0], k=2
// → 6 (0を2個フリップして [0,0,0,1,1,1,1,0] → [1,1,1,1,1,1])
```

## 最適化テクニック

### 1. デックを使った最大値追跡（O(n)で各ウィンドウの最大値）

通常のスライディング最大値はO(n×k)だが、デックを使うとO(n)に。

```go
// MaxSlidingWindowOptimized はデックを使った最適化版
func MaxSlidingWindowOptimized(nums []int, k int) []int {
	n := len(nums)
	if n == 0 || k == 0 {
		return []int{}
	}

	result := make([]int, 0, n-k+1)
	deque := []int{} // インデックスを格納（降順に値が並ぶ）

	for i := 0; i < n; i++ {
		// ウィンドウ外の要素を削除
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}

		// 現在の要素より小さい要素を後ろから削除
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}

		deque = append(deque, i)

		// ウィンドウサイズに達したら結果に追加
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}
```

### 2. ハッシュマップのサイズを制限

文字カウントでアルファベット26文字限定なら配列を使う。

```go
// LengthOfLongestSubstringAlpha は小文字アルファベット限定版
func LengthOfLongestSubstringAlpha(s string) int {
	charIndex := [26]int{} // a-z のインデックス
	for i := range charIndex {
		charIndex[i] = -1
	}

	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		c := s[right] - 'a'
		if charIndex[c] >= left {
			left = charIndex[c] + 1
		}
		charIndex[c] = right

		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
		}
	}

	return maxLen
}
```

### 3. 早期終了

条件を満たしたら即座に返す。

```go
// ContainsSubarraySum は合計がtargetの部分配列が存在するか
func ContainsSubarraySum(arr []int, target int) bool {
	windowSum := 0
	left := 0

	for right := 0; right < len(arr); right++ {
		windowSum += arr[right]

		for windowSum > target {
			windowSum -= arr[left]
			left++
		}

		if windowSum == target {
			return true // 早期終了
		}
	}

	return false
}
```

## デバッグとテスト

### テストケース

```go
func TestSlidingWindow() {
	// 固定長ウィンドウ
	arr1 := []int{2, 1, 5, 1, 3, 2}
	if MaxSumSubarray(arr1, 3) != 9 {
		t.Error("MaxSumSubarray failed")
	}

	// 可変長ウィンドウ
	arr2 := []int{2, 3, 1, 2, 4, 3}
	if MinSubarrayLen(arr2, 7) != 2 {
		t.Error("MinSubarrayLen failed")
	}

	// 文字列
	s := "eceba"
	if LongestSubstringKDistinct(s, 2) != 3 {
		t.Error("LongestSubstringKDistinct failed")
	}

	// エッジケース
	if MaxSumSubarray([]int{}, 3) != 0 {
		t.Error("Empty array should return 0")
	}
	if MinSubarrayLen([]int{1, 2, 3}, 100) != 0 {
		t.Error("Impossible target should return 0")
	}
}
```

## まとめ

### メリット

- 時間計算量O(n)で効率的
- 総当たりO(n²)やO(n×k)を大幅に改善
- 実装がシンプルで理解しやすい
- 幅広い問題に応用可能
- ストリーミングデータにも適用可能

### デメリット

- 連続する部分列にしか適用できない
- 非連続の部分集合問題には使えない
- 条件によっては複雑な管理が必要

### 使うべき時

- **連続する部分配列/文字列の問題**
- **固定長の統計情報**: 移動平均、直近k個の合計
- **可変長の最適化**: 条件を満たす最小/最大長
- **パフォーマンス重視**: O(n²)をO(n)に改善したい

### 選択基準

| 問題の種類 | アプローチ | 時間計算量 |
| --- | --- | --- |
| 固定長の最大/最小合計 | 固定長ウィンドウ | O(n) |
| 条件を満たす最小/最大長 | 可変長ウィンドウ | O(n) |
| 各ウィンドウの最大/最小値 | デック最適化 | O(n) |
| 非連続の部分集合 | 動的計画法 or 全探索 | O(2^n) or O(n×target) |

スライディングウィンドウは配列や文字列の連続する部分列を効率的に処理する強力な技法。Two Pointersの応用として、左右のポインタを巧みに操作することで、総当たりでは不可能な性能を実現する。移動平均、異常検知、文字列マッチングなど、実用的な応用も多い重要アルゴリズム。
