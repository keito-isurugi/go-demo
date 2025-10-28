## 概要

- Goで配列を再帰的に逆順にするサンプルコードを書いてみた
- 配列の逆順は、配列の要素を前後から入れ替える基本的なアルゴリズム。再帰を使うことで、ループを使わずにエレガントに実装できる。

### 特徴

- **分割統治法**: 配列を左右から中央に向かって処理
- **インプレース処理**: 追加メモリを最小限に抑える実装が可能
- **再帰の理解**: 基本的な再帰パターンの良い例
- **複数の実装方法**: インプレース版とスライス操作版
- **シンプルな実装**: 少ないコード行数で実現可能

## 動作原理

基本的な流れ：

1. 配列の両端の要素を交換
2. 範囲を内側に狭める
3. 中央に到達したら終了
4. 再帰的に処理を繰り返す

### 再帰的解法

配列を逆順にする手順：

1. 配列の先頭と末尾の要素を交換
2. 残りの部分配列（先頭+1から末尾-1）を再帰的に逆順
3. 中央に到達（start >= end）したら終了

### 具体例

5要素の配列[1, 2, 3, 4, 5]の逆順過程：

```
初期状態:     [1, 2, 3, 4, 5]
             ↑           ↑
           start        end

ステップ1:    [5, 2, 3, 4, 1]  (1と5を交換)
                ↑     ↑
              start  end

ステップ2:    [5, 4, 3, 2, 1]  (2と4を交換)
                   ↑↑
               start,end

ステップ3:    終了（start >= end）

最終結果:     [5, 4, 3, 2, 1]
```

## サンプルコード

### 基本実装（インプレース版）

```go
package recursion

import "fmt"

// ReverseArray は配列を再帰的に逆順にする
func ReverseArray(arr []int, start, end int) []int {
    // ベースケース: 配列の中央に到達したら終了
    if start >= end {
        return arr
    }
    
    // 要素を交換
    arr[start], arr[end] = arr[end], arr[start]
    
    // 再帰呼び出し: 範囲を狭めて継続
    return ReverseArray(arr, start+1, end-1)
}
```

### スライス操作版

```go
// ReverseArrayWithSlice はスライス操作で配列を逆順にする（再帰）
func ReverseArrayWithSlice(arr []int) []int {
    // ベースケース: 要素が1つ以下なら逆順にする必要なし
    if len(arr) <= 1 {
        return arr
    }
    
    // 最初の要素を取り出し、残りを再帰的に逆順にしてから結合
    return append(ReverseArrayWithSlice(arr[1:]), arr[0])
}
```

### ステップ表示版

```go
// ReverseArrayWithSteps は再帰呼び出しの過程を表示
func ReverseArrayWithSteps(arr []int, start, end, depth int) []int {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    
    fmt.Printf("%sReverseArray([%v], start=%d, end=%d)\n", 
               indent, arr, start, end)
    
    if start >= end {
        fmt.Printf("%s  ベースケース到達: 終了\n", indent)
        return arr
    }
    
    fmt.Printf("%s  交換: arr[%d]=%d <-> arr[%d]=%d\n", 
        indent, start, arr[start], end, arr[end])
    arr[start], arr[end] = arr[end], arr[start]
    fmt.Printf("%s  交換後: %v\n", indent, arr)
    
    return ReverseArrayWithSteps(arr, start+1, end-1, depth+1)
}
```

### 反復版（比較用）

```go
// ReverseArrayWithFor は反復処理で配列を逆順にする
func ReverseArrayWithFor(arr []int) []int {
    n := len(arr)
    result := make([]int, n)
    copy(result, arr)
    
    for i := 0; i < n/2; i++ {
        result[i], result[n-1-i] = result[n-1-i], result[i]
    }
    
    return result
}
```

### デモ実装

```go
// ReverseArrayDemo は配列逆順のデモを実行
func ReverseArrayDemo() {
    fmt.Println("■ 配列の逆順（再帰）デモ")
    fmt.Println("==================================================")
    
    // 基本的な例
    fmt.Println("\n1. 基本的な配列の逆順:")
    arr1 := []int{1, 2, 3, 4, 5}
    fmt.Printf("元の配列: %v\n", arr1)
    
    // コピーを作成して逆順にする
    reversed1 := make([]int, len(arr1))
    copy(reversed1, arr1)
    reversed1 = ReverseArray(reversed1, 0, len(reversed1)-1)
    fmt.Printf("逆順後:   %v\n", reversed1)
    
    // スライス操作版
    fmt.Println("\n2. スライス操作での逆順:")
    arr2 := []int{10, 20, 30, 40, 50, 60}
    fmt.Printf("元の配列: %v\n", arr2)
    reversed2 := ReverseArrayWithSlice(arr2)
    fmt.Printf("逆順後:   %v\n", reversed2)
}
```

### 実行例

```go
// main.goでの使用例
package main

import (
    "fmt"
    "algorithm/recursion"
)

func main() {
    fmt.Println("\n--- 配列の逆順（再帰）のデモ ---")
    recursion.ReverseArrayDemo()
}
```

実行結果：

```
--- 配列の逆順（再帰）のデモ ---
■ 配列の逆順（再帰）デモ
==================================================

1. 基本的な配列の逆順:
元の配列: [1 2 3 4 5]
逆順後:   [5 4 3 2 1]

2. スライス操作での逆順:
元の配列: [10 20 30 40 50 60]
逆順後:   [60 50 40 30 20 10]

3. 再帰呼び出しの詳細:
初期配列: [1 2 3 4 5]
再帰呼び出しの過程:
ReverseArray([[1 2 3 4 5]], start=0, end=4)
  交換: arr[0]=1 <-> arr[4]=5
  交換後: [5 2 3 4 1]
  ReverseArray([[5 2 3 4 1]], start=1, end=3)
    交換: arr[1]=2 <-> arr[3]=4
    交換後: [5 4 3 2 1]
    ReverseArray([[5 4 3 2 1]], start=2, end=2)
      ベースケース到達: 終了
最終結果: [5 4 3 2 1]
```

## 計算量

### 時間計算量

- **インプレース版**: O(n/2) = O(n) - 各要素を1回だけ処理
- **スライス操作版**: O(n²) - 各ステップでスライスのコピーが発生
- **反復版**: O(n) - 各要素を1回だけ処理

### 空間計算量

- **インプレース版**: O(n) - 再帰スタックの深さ（最大n/2）
- **スライス操作版**: O(n²) - 各再帰で新しいスライスを作成
- **反復版**: O(1) - 追加メモリは定数

### 配列サイズと処理性能

| 配列サイズ | 再帰呼び出し回数 | 交換回数 |
| --- | --- | --- |
| 5 | 3 | 2 |
| 10 | 6 | 5 |
| 100 | 51 | 50 |
| 1000 | 501 | 500 |
| n | n/2 + 1 | n/2 |

## 使いどころ

### 向いてる場面

- 再帰アルゴリズムの学習
- 分割統治法の理解
- 小規模な配列の処理
- コードの可読性を重視する場合

### 実例：文字列の逆順

```go
// ReverseString は文字列を逆順にする
func ReverseString(s string) string {
    runes := []rune(s)
    reverseRunes(runes, 0, len(runes)-1)
    return string(runes)
}

func reverseRunes(runes []rune, start, end int) {
    if start >= end {
        return
    }
    runes[start], runes[end] = runes[end], runes[start]
    reverseRunes(runes, start+1, end-1)
}

// 使用例
func main() {
    str := "Hello, 世界!"
    reversed := ReverseString(str)
    fmt.Printf("元の文字列: %s\n", str)
    fmt.Printf("逆順:       %s\n", reversed)
    // 出力:
    // 元の文字列: Hello, 世界!
    // 逆順:       !界世 ,olleH
}
```

### 実例：回文判定

```go
// IsPalindrome は配列が回文かどうかを判定
func IsPalindrome(arr []int) bool {
    return isPalindromeHelper(arr, 0, len(arr)-1)
}

func isPalindromeHelper(arr []int, start, end int) bool {
    // ベースケース: 中央に到達
    if start >= end {
        return true
    }
    
    // 両端が異なれば回文ではない
    if arr[start] != arr[end] {
        return false
    }
    
    // 内側を再帰的にチェック
    return isPalindromeHelper(arr, start+1, end-1)
}

// 使用例
func main() {
    palindrome := []int{1, 2, 3, 2, 1}
    notPalindrome := []int{1, 2, 3, 4, 5}
    
    fmt.Printf("%v は回文: %v\n", palindrome, IsPalindrome(palindrome))
    fmt.Printf("%v は回文: %v\n", notPalindrome, IsPalindrome(notPalindrome))
    // 出力:
    // [1 2 3 2 1] は回文: true
    // [1 2 3 4 5] は回文: false
}
```

## 他の実装との比較

| 実装方法 | 時間計算量 | 空間計算量 | 可読性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| 再帰（インプレース） | O(n) | O(n) | ◎ | 易 |
| 再帰（スライス） | O(n²) | O(n²) | ○ | 易 |
| 反復（Two Pointers） | O(n) | O(1) | ◎ | 易 |
| 組み込み関数 | O(n) | O(n) | ◎ | 最易 |

## 最適化アイデア

### 尾再帰最適化

```go
// TailRecursiveReverse は尾再帰で配列を逆順にする
// （Goは尾再帰最適化をサポートしていないが、概念として）
func TailRecursiveReverse(arr []int, start, end int) {
    if start >= end {
        return
    }
    
    arr[start], arr[end] = arr[end], arr[start]
    
    // 尾位置での再帰呼び出し
    TailRecursiveReverse(arr, start+1, end-1)
}
```

### ジェネリクス版

```go
// ReverseGeneric はジェネリクスを使った汎用逆順関数
func ReverseGeneric[T any](slice []T) []T {
    result := make([]T, len(slice))
    copy(result, slice)
    reverseHelper(result, 0, len(result)-1)
    return result
}

func reverseHelper[T any](slice []T, start, end int) {
    if start >= end {
        return
    }
    slice[start], slice[end] = slice[end], slice[start]
    reverseHelper(slice, start+1, end-1)
}

// 使用例
func main() {
    // 整数配列
    ints := []int{1, 2, 3, 4, 5}
    reversedInts := ReverseGeneric(ints)
    fmt.Println(reversedInts)  // [5 4 3 2 1]
    
    // 文字列配列
    strs := []string{"a", "b", "c", "d"}
    reversedStrs := ReverseGeneric(strs)
    fmt.Println(reversedStrs)  // [d c b a]
}
```

### ベンチマーク

```go
import "testing"

func BenchmarkReverseRecursive(b *testing.B) {
    arr := make([]int, 1000)
    for i := range arr {
        arr[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        temp := make([]int, len(arr))
        copy(temp, arr)
        ReverseArray(temp, 0, len(temp)-1)
    }
}

func BenchmarkReverseIterative(b *testing.B) {
    arr := make([]int, 1000)
    for i := range arr {
        arr[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ReverseArrayWithFor(arr)
    }
}
```

## まとめ

### メリット

- エレガントで読みやすいコード
- 分割統治法の良い学習例
- 少ないコード行数で実装可能
- 再帰的思考の訓練に最適

### デメリット

- スタックオーバーフローの可能性（大規模配列）
- 反復版より若干遅い（関数呼び出しのオーバーヘッド）
- Goでは尾再帰最適化が効かない
- スライス操作版は効率が悪い

### 使うべき時

- 教育目的や学習用途
- 小〜中規模の配列処理
- コードの可読性を重視する場合
- 再帰的な問題の一部として使用

### 学習ポイント

- **再帰の基本**: ベースケースと再帰ケースの設計
- **インデックス操作**: 配列の両端から中央への処理
- **分割統治**: 問題を小さく分割して解決
- **トレードオフ**: 可読性と性能のバランス

配列の逆順は再帰アルゴリズムの入門として最適な問題。シンプルな実装で再帰の基本概念を理解でき、より複雑な再帰問題への足がかりとなる。