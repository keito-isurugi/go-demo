## 概要

- Goでフィボナッチ数列を再帰実装
- フィボナッチ数列は前の2つの数を足した値が次の数になる数列。再帰は数学的定義をそのままコードに表現できるが、計算量の問題がある典型例。

### 特徴

- **直感的実装**: 数学的定義をそのままコード化
- **関数型プログラミング**: 副作用のない純粋関数
- **指数的計算量**: 同じ計算を何度も繰り返すため非効率
- **メモ化で改善可能**: 計算結果を保存することで高速化
- **再帰の典型例**: 再帰アルゴリズムの学習に最適

## 動作原理

基本的な流れ：

1. ベースケース：Fib(0) = 0, Fib(1) = 1
2. 再帰ケース：Fib(n) = Fib(n-1) + Fib(n-2)
3. 各関数呼び出しが2つの子問題に分割
4. 結果を合算して返す

### 具体例

`Fib(5)` の計算過程：

```
Fib(5)
├── Fib(4)
│   ├── Fib(3)
│   │   ├── Fib(2)
│   │   │   ├── Fib(1) = 1
│   │   │   └── Fib(0) = 0
│   │   │   = 1 + 0 = 1
│   │   └── Fib(1) = 1
│   │   = 1 + 1 = 2
│   └── Fib(2)
│       ├── Fib(1) = 1
│       └── Fib(0) = 0
│       = 1 + 0 = 1
│   = 2 + 1 = 3
└── Fib(3)
    ├── Fib(2)
    │   ├── Fib(1) = 1
    │   └── Fib(0) = 0
    │   = 1 + 0 = 1
    └── Fib(1) = 1
    = 1 + 1 = 2
= 3 + 2 = 5

結果: Fib(5) = 5
```

## サンプルコード

### 基本実装

```go
package main

import "fmt"

// FibonacciRecursive は再帰でフィボナッチ数列のn番目を計算
func FibonacciRecursive(n int) int {
    if n <= 1 {
        return n
    }
    return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

```

### ステップ表示版

```go
// FibonacciRecursiveWithSteps は再帰の呼び出し過程を表示
func FibonacciRecursiveWithSteps(n int, depth int) int {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    
    fmt.Printf("%sFib(%d) を計算開始\\n", indent, n)
    
    if n <= 1 {
        fmt.Printf("%sFib(%d) = %d (ベースケース)\\n", indent, n, n)
        return n
    }
    
    fmt.Printf("%sFib(%d) = Fib(%d) + Fib(%d)\\n", indent, n, n-1, n-2)
    
    left := FibonacciRecursiveWithSteps(n-1, depth+1)
    right := FibonacciRecursiveWithSteps(n-2, depth+1)
    
    result := left + right
    fmt.Printf("%sFib(%d) = %d + %d = %d\\n", indent, n, left, right, result)
    
    return result
}

```

### メモ化版（改善）

```go
// FibonacciMemoization はメモ化を使った再帰実装
func FibonacciMemoization(n int) int {
    memo := make(map[int]int)
    return fibMemo(n, memo)
}

func fibMemo(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    
    // メモに結果があれば使用
    if val, exists := memo[n]; exists {
        return val
    }
    
    // 計算してメモに保存
    result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
    memo[n] = result
    return result
}

```

### 反復版（比較用）

```go
// FibonacciIterative は反復による実装
func FibonacciIterative(n int) int {
    if n <= 1 {
        return n
    }
    
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

```

## 計算量

### 時間計算量

- **純粋再帰**: O(2^n) - 指数的増加
- **メモ化版**: O(n) - 各値を1回だけ計算
- **反復版**: O(n) - 線形時間

### 空間計算量

- **純粋再帰**: O(n) - 再帰スタックの深さ
- **メモ化版**: O(n) - メモテーブル + スタック
- **反復版**: O(1) - 定数空間

### 関数呼び出し回数

Fib(n)の計算における関数呼び出し回数：
- Fib(5): 15回
- Fib(10): 177回  
- Fib(20): 21,891回
- Fib(30): 2,692,537回

## 使いどころ

### 向いてる場面

- 小さなn（n < 30程度）
- 数学的概念の理解・教育
- 再帰アルゴリズムの学習
- プロトタイピング

### 実例：数列生成

```go
// FibonacciSequence はn番目までのフィボナッチ数列を生成
func FibonacciSequence(n int) []int {
    if n < 0 {
        return []int{}
    }
    
    sequence := make([]int, n+1)
    for i := 0; i <= n; i++ {
        sequence[i] = FibonacciRecursive(i)
    }
    return sequence
}

// 使用例
func main() {
    sequence := FibonacciSequence(10)
    fmt.Printf("フィボナッチ数列: %v\\n", sequence)
    // 出力: [0 1 1 2 3 5 8 13 21 34 55]
}

```

## 他の実装との比較

| 実装方法 | 時間計算量 | 空間計算量 | 可読性 | 実用性 |
| --- | --- | --- | --- | --- |
| 純粋再帰 | O(2^n) | O(n) | ◎ | × |
| メモ化再帰 | O(n) | O(n) | ○ | ○ |
| 反復 | O(n) | O(1) | ○ | ◎ |
| 行列べき乗 | O(log n) | O(1) | × | ◎ |

## 最適化アイデア

### パフォーマンス測定

```go
import "time"

func measureTime(name string, fn func()) {
    start := time.Now()
    fn()
    duration := time.Since(start)
    fmt.Printf("%s の実行時間: %v\\n", name, duration)
}

func compareImplementations(n int) {
    fmt.Printf("n = %d での比較:\\n", n)
    
    measureTime("純粋再帰", func() {
        FibonacciRecursive(n)
    })
    
    measureTime("メモ化", func() {
        FibonacciMemoization(n)
    })
    
    measureTime("反復", func() {
        FibonacciIterative(n)
    })
}

```

### 末尾再帰版

```go
// FibonacciTailRecursive は末尾再帰による実装
func FibonacciTailRecursive(n int) int {
    return fibTail(n, 0, 1)
}

func fibTail(n, a, b int) int {
    if n == 0 {
        return a
    }
    return fibTail(n-1, b, a+b)
}

```

### チャネルを使った並行実装

```go
func FibonacciChannel(n int) <-chan int {
    ch := make(chan int)
    
    go func() {
        defer close(ch)
        a, b := 0, 1
        for i := 0; i < n; i++ {
            ch <- a
            a, b = b, a+b
        }
    }()
    
    return ch
}

// 使用例
func main() {
    for fib := range FibonacciChannel(10) {
        fmt.Printf("%d ", fib)
    }
    // 出力: 0 1 1 2 3 5 8 13 21 34
}

```

## まとめ

### メリット

- 直感的で理解しやすい
- 数学的定義に忠実
- 再帰の良い学習例
- コードが簡潔

### デメリット

- 指数的計算量で非効率
- 大きなnで実用不可
- スタックオーバーフローの危険
- 同じ計算の重複実行

### 使うべき時

- 小さなn（< 30）
- 教育・学習目的
- プロトタイピング
- 再帰の理解

### 改善案

- **メモ化**: 計算結果をキャッシュ
- **反復版**: ボトムアップで計算
- **行列べき乗**: O(log n)で計算

実用的にはメモ化や反復版を使うべきだが、純粋再帰版は再帰アルゴリズムの典型例として重要。計算量の違いを実感できる良い教材。