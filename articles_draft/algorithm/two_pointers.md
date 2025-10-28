## 概要

- GoでTwo Pointers（二つのポインタ）を実装
- Two Pointersは配列や文字列を2つのポインタで走査する技法。左右のポインタを巧みに操作することで、効率的に問題を解く。

### 特徴

- **効率的処理**: 時間計算量O(n)で線形探索可能
- **空間計算量O(1)**: 追加のメモリをほとんど使わない
- **多様なパターン**: 対向ポインタ、同方向ポインタ、固定+可変など
- **幅広い応用**: ソート済み配列、文字列、部分配列問題に適用可能
- **Sliding Windowとの関連**: Sliding Windowは可変長Two Pointersの応用

## Two Pointersとは

### 基本的な考え方

配列や文字列に対して2つのインデックス（ポインタ）を使い、それぞれを適切に移動させながら問題を解く手法。

**基本パターン**:
```
1. 対向ポインタ（Opposite Direction）:
   配列: [1, 2, 3, 4, 5]
          ↑           ↑
        left       right
   → 左右から中央に向かって移動

2. 同方向ポインタ（Same Direction）:
   配列: [1, 2, 3, 4, 5]
          ↑     ↑
        slow  fast
   → 同じ方向に異なる速度で移動

3. 固定+可変ポインタ:
   配列: [1, 2, 3, 4, 5]
          ↑     ↑
          i   left/right
   → 一方を固定、もう一方を可変
```

### なぜ効率的か？

**総当たり（Brute Force）**:
- すべての組み合わせを調べる
- 時間計算量: O(n²)

```go
// 総当たり: 2重ループ
for i := 0; i < n; i++ {
    for j := i+1; j < n; j++ {
        if arr[i] + arr[j] == target {
            return []int{i, j}
        }
    }
}
```

**Two Pointers**:
- ポインタを賢く移動させて探索範囲を絞る
- 時間計算量: O(n)

```go
// Two Pointers: 1回のループ
left, right := 0, n-1
for left < right {
    sum := arr[left] + arr[right]
    if sum == target {
        return []int{left, right}
    } else if sum < target {
        left++  // 合計を増やす
    } else {
        right-- // 合計を減らす
    }
}
```

**性能差**:
```
配列サイズ n=10,000
総当たり:     約5000万回の計算 (10,000 × 10,000 / 2)
Two Pointers: 約1万回の計算 (10,000)
→ 5000倍高速！
```

## サンプルコード

### 1. Two Sum（ソート済み配列）

最も基本的なパターン。昇順配列から合計が目標値になる2要素を見つける。

```go
package main

import "fmt"

// TwoSumSorted は昇順配列から合計がtargetになる2要素のインデックスを返す
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
```

**使用例**:
```go
nums := []int{2, 7, 11, 15}
target := 9
result := TwoSumSorted(nums, target)
// 結果: [0, 1] (nums[0]=2, nums[1]=7, 2+7=9)
```

### 2. 重複の削除（ソート済み配列）

ソート済み配列から重複を削除し、ユニーク要素の数を返す。

```go
// RemoveDuplicates はソート済み配列から重複を削除
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
```

**動作イメージ**:
```
配列: [1, 1, 2, 2, 3]
       ↑
     slow,fast

fast=1: nums[1]=1 = nums[0]=1 → スキップ
fast=2: nums[2]=2 ≠ nums[0]=1 → slow++, nums[1]=2
配列: [1, 2, 2, 2, 3]
          ↑  ↑
        slow fast

fast=3: nums[3]=2 = nums[1]=2 → スキップ
fast=4: nums[4]=3 ≠ nums[1]=2 → slow++, nums[2]=3
配列: [1, 2, 3, 2, 3]
             ↑     ↑
           slow   fast

結果: ユニーク要素数 = 3, 配列 = [1, 2, 3]
```

### 3. 文字列の反転

文字列を反転する（Two Pointers）。

```go
// ReverseString は文字列を反転する
func ReverseString(s []byte) {
	left := 0
	right := len(s) - 1

	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}
```

**使用例**:
```go
s := []byte("hello")
ReverseString(s)
// 結果: "olleh"
```

### 4. 回文判定

文字列が回文（前から読んでも後ろから読んでも同じ）かどうかを判定。

```go
// IsPalindrome は文字列が回文かどうかを判定
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
```

**使用例**:
```go
IsPalindrome("racecar") // true
IsPalindrome("hello")   // false
```

### 5. Three Sum（3要素の合計が0）

配列から合計が0になる3要素の組み合わせをすべて返す。

```go
// ThreeSum は配列から合計が0になる3要素の組み合わせをすべて返す
func ThreeSum(nums []int) [][]int {
	result := [][]int{}
	n := len(nums)

	// まずソート
	sort.Ints(nums)

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
```

**使用例**:
```go
nums := []int{-1, 0, 1, 2, -1, -4}
result := ThreeSum(nums)
// 結果: [[-1, -1, 2], [-1, 0, 1]]
```

**アルゴリズムの流れ**:
```
ソート後: [-4, -1, -1, 0, 1, 2]

i=0: nums[0]=-4 を固定、残り2要素の合計目標値=4
     Two Pointersで探索 → 見つからない

i=1: nums[1]=-1 を固定、残り2要素の合計目標値=1
     left=2, right=5
     -1 + 2 = 1 ✓ → [-1, -1, 2]
     left=3, right=4
     0 + 1 = 1 ✓ → [-1, 0, 1]

i=2: nums[2]=-1 は重複のためスキップ
```

### 6. Container With Most Water

最大の水を貯められる2つの垂直線を見つける。

```go
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

**なぜ低い方を移動するのか**:
```
高さ配列: [1, 8, 6, 2, 5, 4, 8, 3, 7]
           ↑                       ↑
         left(1)              right(7)

現在の面積 = 8 × min(1, 7) = 8

もし right を移動すると:
  → 幅が減る（確実にマイナス）
  → 高さは max(1, next) = 1 （変わらないか悪化）
  → 面積は必ず悪化

もし left を移動すると:
  → 幅が減る（マイナス）
  → 高さは max(next, 7) （改善の可能性あり）
  → 面積が改善する可能性あり

結論: 低い方を移動して改善を狙う
```

### 7. Move Zeroes（0を末尾に移動）

配列内の0をすべて末尾に移動する（相対順序は保持）。

```go
// MoveZeroes は配列内の0をすべて末尾に移動する
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
```

**使用例**:
```go
nums := []int{0, 1, 0, 3, 12}
MoveZeroes(nums)
// 結果: [1, 3, 12, 0, 0]
```

## 動作原理

### 対向ポインタの動き

配列 `[2, 7, 11, 15]`、目標値 `target=9` の例：

```
初期状態:
配列: [2, 7, 11, 15]
       ↑           ↑
     left(0)    right(3)

Step 1: sum = 2 + 15 = 17 > 9
        → 合計が大きい → right--

配列: [2, 7, 11, 15]
       ↑       ↑
     left(0) right(2)

Step 2: sum = 2 + 11 = 13 > 9
        → 合計が大きい → right--

配列: [2, 7, 11, 15]
       ↑   ↑
     left(0) right(1)

Step 3: sum = 2 + 7 = 9 ✓
        → 見つかりました!

結果: インデックス [0, 1]
```

### 同方向ポインタの動き

配列 `[1, 1, 2, 2, 3]` の重複削除：

```
初期状態:
配列: [1, 1, 2, 2, 3]
       ↑
     slow=0, fast=0

fast=1: nums[1]=1 = nums[0]=1
        → 重複のためスキップ

fast=2: nums[2]=2 ≠ nums[0]=1
        → slow++, nums[1]=2
配列: [1, 2, 2, 2, 3]
          ↑  ↑
        slow fast

fast=3: nums[3]=2 = nums[1]=2
        → 重複のためスキップ

fast=4: nums[4]=3 ≠ nums[1]=2
        → slow++, nums[2]=3
配列: [1, 2, 3, 2, 3]
             ↑     ↑
           slow   fast

結果: ユニーク要素数 = 3
```

## 3つのパターン

### パターン1: 対向ポインタ（Opposite Direction）

左右から中央に向かって移動。

**特徴**:
- 左端(left)と右端(right)からスタート
- 条件に応じてどちらかを移動
- ソート済み配列に有効

**典型的な問題**:
- Two Sum（ソート済み配列）
- 回文判定
- Container With Most Water

**実装パターン**:
```go
left := 0
right := len(arr) - 1

for left < right {
    if 条件を満たす {
        // 処理
        break または return
    } else if 値が小さい {
        left++
    } else {
        right--
    }
}
```

### パターン2: 同方向ポインタ（Same Direction）

同じ方向に異なる速度で移動。

**特徴**:
- slowとfastポインタを使用
- fastが先行、slowが追従
- インプレース操作に有効

**典型的な問題**:
- 重複の削除
- 特定要素の削除
- Move Zeroes

**実装パターン**:
```go
slow := 0

for fast := 0; fast < len(arr); fast++ {
    if 条件を満たす {
        arr[slow] = arr[fast]
        slow++
    }
}
```

### パターン3: 固定+可変ポインタ

一方を固定、もう一方を可変にして探索。

**特徴**:
- 外側ループで一方を固定
- 内側でTwo Pointersを使用
- Three Sumなどの複雑な問題に有効

**典型的な問題**:
- Three Sum
- Four Sum
- Closest Sum

**実装パターン**:
```go
for i := 0; i < n; i++ {
    // i を固定
    left := i + 1
    right := n - 1

    for left < right {
        // Two Pointers探索
        if 条件を満たす {
            // 処理
        } else if 値が小さい {
            left++
        } else {
            right--
        }
    }
}
```

## 計算量

### 時間計算量

| 実装方法 | 時間計算量 | 説明 |
| --- | --- | --- |
| 総当たり（2要素） | O(n²) | すべての組み合わせを調べる |
| Two Pointers（2要素） | O(n) | 各要素を最大1回処理 |
| 総当たり（3要素） | O(n³) | 3重ループ |
| Three Sum (固定+Two Pointers) | O(n²) | 外側O(n)×内側O(n) |

### 空間計算量

| 実装方法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 基本Two Pointers | O(1) | ポインタ2つのみ |
| インプレース操作 | O(1) | 追加メモリ不要 |
| 結果を保存 | O(k) | k個の結果を保存 |

### なぜO(n)なのか？

対向ポインタの場合：

```go
left := 0
right := n - 1

for left < right {  // 最大n回
    // 処理
    left++ または right--
}
```

**重要なポイント**:
- `left`は0からn-1まで最大n回増加
- `right`はn-1から0まで最大n回減少
- 合計でも最大n回のループ
- したがって時間計算量はO(n)

```
各ステップでleftまたはrightのどちらかが1つ近づく:
初期: left=0, right=n-1, 差=n-1
...
終了: left=right, 差=0

→ 最大n回のステップで終了
```

## 使いどころ

### 向いている場面

1. **ソート済み配列の探索**
   - Two Sum、Three Sum
   - 特定の合計値を見つける

2. **インプレース操作**
   - 重複削除
   - 要素の移動・整理

3. **回文や対称性の判定**
   - 回文判定
   - 対称配列のチェック

4. **最適化問題**
   - Container With Most Water
   - 最小・最大値の探索

### 実例1: Two Sum（未ソート配列）

ハッシュマップとTwo Pointersを組み合わせる。

```go
// TwoSum は未ソート配列から合計がtargetになる2要素のインデックスを返す
func TwoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if index, ok := numMap[complement]; ok {
			return []int{index, i}
		}
		numMap[num] = i
	}

	return []int{}
}

// 使用例
nums := []int{2, 7, 11, 15}
target := 9
result := TwoSum(nums, target)
// 結果: [0, 1]
```

### 実例2: 配列の分割

配列を特定の条件で2つに分割する（Dutch National Flag問題）。

```go
// PartitionArray は配列をpivot未満とpivot以上に分割
func PartitionArray(nums []int, pivot int) int {
	left := 0
	right := len(nums) - 1

	for left <= right {
		if nums[left] < pivot {
			left++
		} else {
			nums[left], nums[right] = nums[right], nums[left]
			right--
		}
	}

	return left // pivot以上の要素の開始位置
}

// 使用例
nums := []int{3, 5, 2, 1, 4}
pivot := 3
pos := PartitionArray(nums, pivot)
// 結果: nums=[2, 1, 3, 5, 4], pos=2
//       [2, 1] < 3, [3, 5, 4] >= 3
```

### 実例3: 最短距離のペア

ソート済み配列から差が最小の2要素を見つける。

```go
// SmallestDifferencePair はソート済み配列から差が最小の2要素を返す
func SmallestDifferencePair(arr1, arr2 []int) []int {
	i, j := 0, 0
	minDiff := int(^uint(0) >> 1) // 最大int値
	result := []int{0, 0}

	for i < len(arr1) && j < len(arr2) {
		diff := abs(arr1[i] - arr2[j])

		if diff < minDiff {
			minDiff = diff
			result[0] = arr1[i]
			result[1] = arr2[j]
		}

		if arr1[i] < arr2[j] {
			i++
		} else {
			j++
		}
	}

	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
```

## 関連する問題パターン

### 1. 重複のない最長部分文字列

Two Pointersで重複のない最長部分文字列の長さを求める。

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
// 例: "pwwkew" → 3 ("wke")
```

### 2. Sort Colors（3色ソート）

0, 1, 2の3種類の要素を1回のパスでソートする。

```go
// SortColors は配列を1回のパスで[0,1,2]にソート
func SortColors(nums []int) {
	low, mid, high := 0, 0, len(nums)-1

	for mid <= high {
		switch nums[mid] {
		case 0:
			nums[low], nums[mid] = nums[mid], nums[low]
			low++
			mid++
		case 1:
			mid++
		case 2:
			nums[mid], nums[high] = nums[high], nums[mid]
			high--
		}
	}
}

// 使用例
nums := []int{2, 0, 2, 1, 1, 0}
SortColors(nums)
// 結果: [0, 0, 1, 1, 2, 2]
```

### 3. Trapping Rain Water（雨水トラップ）

高さ配列から雨水がどれだけ溜まるかを計算。

```go
// TrapRainWater は雨水が溜まる量を計算
func TrapRainWater(height []int) int {
	if len(height) == 0 {
		return 0
	}

	left, right := 0, len(height)-1
	leftMax, rightMax := 0, 0
	water := 0

	for left < right {
		if height[left] < height[right] {
			if height[left] >= leftMax {
				leftMax = height[left]
			} else {
				water += leftMax - height[left]
			}
			left++
		} else {
			if height[right] >= rightMax {
				rightMax = height[right]
			} else {
				water += rightMax - height[right]
			}
			right--
		}
	}

	return water
}

// 例: [0,1,0,2,1,0,1,3,2,1,2,1]
// 結果: 6
```

## Sliding Windowとの関係

Sliding WindowはTwo Pointersの応用パターンの一つ。

### 違い

| 特徴 | Two Pointers | Sliding Window |
| --- | --- | --- |
| 目的 | 2要素の組み合わせ探索 | 連続部分配列の処理 |
| ポインタの動き | 独立に移動 | ウィンドウとして移動 |
| 典型的な用途 | ソート済み配列探索 | 部分配列の最大/最小 |
| 実装 | left/rightを条件で移動 | 固定長または可変長ウィンドウ |

### Sliding Windowの例

```go
// MaxSumSubarray は固定長kの部分配列の最大合計（Sliding Window）
func MaxSumSubarray(arr []int, k int) int {
	if len(arr) < k {
		return 0
	}

	// 最初のウィンドウ
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}

	maxSum := windowSum

	// ウィンドウをスライド
	for i := k; i < len(arr); i++ {
		windowSum = windowSum + arr[i] - arr[i-k]
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum
}
```

## 最適化テクニック

### 1. 早期終了

条件を満たしたら即座に返す。

```go
// HasPairWithSum は合計がtargetのペアが存在するか
func HasPairWithSum(nums []int, target int) bool {
	left, right := 0, len(nums)-1

	for left < right {
		sum := nums[left] + nums[right]
		if sum == target {
			return true // 早期終了
		} else if sum < target {
			left++
		} else {
			right--
		}
	}

	return false
}
```

### 2. 重複のスキップ

重複する要素をスキップして計算量を削減。

```go
// Three Sumで重複をスキップ
for i := 0; i < n-2; i++ {
	// 同じ値をスキップ
	if i > 0 && nums[i] == nums[i-1] {
		continue
	}
	// ...
}
```

### 3. ソート済み配列の活用

配列がソート済みであることを活用して効率化。

```go
// ソート済み配列で目標値を超えたら終了
for left < right {
	sum := nums[left] + nums[right]
	if sum > target {
		right--
	} else if sum < target {
		left++
	} else {
		return []int{left, right}
	}
}
```

## デバッグとテスト

### テストケース

```go
func TestTwoPointers() {
	// Two Sum
	nums1 := []int{2, 7, 11, 15}
	if !reflect.DeepEqual(TwoSumSorted(nums1, 9), []int{0, 1}) {
		t.Error("TwoSumSorted failed")
	}

	// 重複削除
	nums2 := []int{1, 1, 2}
	if RemoveDuplicates(nums2) != 2 {
		t.Error("RemoveDuplicates failed")
	}

	// 回文判定
	if !IsPalindrome("racecar") {
		t.Error("IsPalindrome failed")
	}

	// エッジケース
	if len(TwoSumSorted([]int{}, 0)) != 0 {
		t.Error("Empty array should return empty")
	}

	if IsPalindrome("") {
		t.Error("Empty string should not be palindrome")
	}
}
```

### デバッグのポイント

1. **ポインタの範囲チェック**
   ```go
   // 範囲外アクセスに注意
   for left < right {  // <= ではなく <
       // ...
   }
   ```

2. **無限ループの回避**
   ```go
   // 必ずポインタが進むことを確認
   if condition {
       left++  // または right--
   } else {
       right-- // または left++
   }
   ```

3. **エッジケースのテスト**
   - 空配列
   - 要素1つ
   - すべて同じ要素
   - 解が存在しない

## まとめ

### メリット

- 時間計算量O(n)で効率的
- 総当たりO(n²)を大幅に改善
- 空間計算量O(1)で省メモリ
- 実装がシンプルで理解しやすい
- 幅広い問題に応用可能

### デメリット

- ソート済み配列に限定される場合が多い
- 未ソート配列では事前ソートが必要（O(n log n)）
- 3要素以上の組み合わせは複雑化

### 使うべき時

- **ソート済み配列の探索**: Two Sum、Three Sum
- **インプレース操作**: 重複削除、要素移動
- **対称性の判定**: 回文、対称配列
- **最適化問題**: 最大面積、最小差
- **パフォーマンス重視**: O(n²)をO(n)に改善したい

### 選択基準

| 問題の種類 | アプローチ | 時間計算量 |
| --- | --- | --- |
| ソート済み配列の2要素探索 | 対向Two Pointers | O(n) |
| 未ソート配列の2要素探索 | ハッシュマップ | O(n) |
| 重複削除・要素移動 | 同方向Two Pointers | O(n) |
| 3要素以上の組み合わせ | 固定+Two Pointers | O(n²) |
| 連続部分配列の処理 | Sliding Window | O(n) |

Two Pointersは配列や文字列を効率的に処理する強力な技法。2つのポインタを巧みに操作することで、総当たりでは不可能な性能を実現する。ソート済み配列の探索、インプレース操作、最適化問題など、実用的な応用も多い重要アルゴリズム。Sliding Windowの基礎としても重要。
