## 概要

- Goでクイックソートを実装
- クイックソートは分割統治法を使う高速なソートアルゴリズム。ピボット（基準値）を選んで、それより小さい要素を左、大きい要素を右に分割することを再帰的に繰り返す。

### 特徴

- **分割統治法**: 配列を分割して再帰的にソート
- **インプレース**: 追加メモリ最小限（O(log n)）
- **不安定ソート**: 同じ値の相対順序が保持されない
- **平均的に高速**: 平均O(n log n)の時間計算量
- **比較ソート**: 要素同士の比較でソート

## 動作原理

基本的な流れ：

1. ピボット（基準値）を選択
2. ピボットより小さい要素を左側、大きい要素を右側に分割
3. 左右の部分配列に対して同じ処理を再帰的に実行
4. 再帰が終われば全体がソート済み

### 具体例

`[3, 7, 8, 5, 2, 1, 9, 4]` のソート過程：

```
初期: [3, 7, 8, 5, 2, 1, 9, 4]

Step1: ピボット=4で分割
[3, 2, 1] [4] [7, 8, 5, 9]

Step2: 左側[3, 2, 1]を処理（ピボット=1）
[1] [3, 2] → [1] [2] [3]

Step3: 右側[7, 8, 5, 9]を処理（ピボット=9）
[7, 8, 5] [9] → [5] [7, 8] [9] → [5] [7] [8] [9]

完了: [1, 2, 3, 4, 5, 7, 8, 9]

```

## サンプルコード

### 基本実装

```go
package main

import "fmt"

// QuickSort は配列を昇順にソート
func QuickSort(arr []int) {
    quickSortHelper(arr, 0, len(arr)-1)
}

// quickSortHelper は再帰的にソート
func quickSortHelper(arr []int, low, high int) {
    if low < high {
        // パーティション操作でピボット位置を取得
        pi := partition(arr, low, high)
        
        // 左右を再帰的にソート
        quickSortHelper(arr, low, pi-1)
        quickSortHelper(arr, pi+1, high)
    }
}

// partition は配列を分割
func partition(arr []int, low, high int) int {
    // 最後の要素をピボットに
    pivot := arr[high]
    i := low - 1
    
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    // ピボットを正しい位置に
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}
```

### ランダムピボット版

```go
import "math/rand"

// QuickSortRandom はランダムピボットを使用
func QuickSortRandom(arr []int) {
    quickSortRandomHelper(arr, 0, len(arr)-1)
}

func quickSortRandomHelper(arr []int, low, high int) {
    if low < high {
        // ランダムにピボットを選択
        randomIndex := low + rand.Intn(high-low+1)
        arr[randomIndex], arr[high] = arr[high], arr[randomIndex]
        
        pi := partition(arr, low, high)
        quickSortRandomHelper(arr, low, pi-1)
        quickSortRandomHelper(arr, pi+1, high)
    }
}

```

### デバッグ用（ステップ表示）

```go
// QuickSortWithSteps はステップごとの状態を表示
func QuickSortWithSteps(arr []int) {
    fmt.Printf("初期: %v\\n", arr)
    quickSortHelperWithSteps(arr, 0, len(arr)-1, 0)
}

func quickSortHelperWithSteps(arr []int, low, high, depth int) {
    if low < high {
        indent := strings.Repeat("  ", depth)
        fmt.Printf("%s範囲[%d:%d]を処理\\n", indent, low, high)
        
        pi := partition(arr, low, high)
        fmt.Printf("%sピボット=%d, 位置=%d\\n", indent, arr[pi], pi)
        fmt.Printf("%s結果: %v\\n\\n", indent, arr)
        
        quickSortHelperWithSteps(arr, low, pi-1, depth+1)
        quickSortHelperWithSteps(arr, pi+1, high, depth+1)
    }
}

```

## 計算量

### 時間計算量

- 最良：**O(n log n)** - バランスよく分割される場合
- 平均：**O(n log n)** - ランダムなデータ
- 最悪：**O(n²)** - 常に偏った分割（ソート済みなど）

### 空間計算量

- **O(log n)** - 再帰呼び出しのスタック（平均）
- **O(n)** - 最悪の場合のスタック深さ

### 比較と交換

- 比較回数: 平均 O(n log n)回、最悪 O(n²)回
- 交換回数: 平均 O(n log n)回

## 使いどころ

### 向いてる場面

- 大規模データ
- 平均的な性能重視
- メモリ制約あり
- ランダムなデータ

### 実例：構造体のソート

```go
type Item struct {
    ID    int
    Value int
}

// Valueで昇順ソート
func SortItemsByValue(items []Item) {
    quickSortItems(items, 0, len(items)-1)
}

func quickSortItems(items []Item, low, high int) {
    if low < high {
        pi := partitionItems(items, low, high)
        quickSortItems(items, low, pi-1)
        quickSortItems(items, pi+1, high)
    }
}

func partitionItems(items []Item, low, high int) int {
    pivot := items[high].Value
    i := low - 1
    
    for j := low; j < high; j++ {
        if items[j].Value <= pivot {
            i++
            items[i], items[j] = items[j], items[i]
        }
    }
    
    items[i+1], items[high] = items[high], items[i+1]
    return i + 1
}

```

## 他のアルゴリズムとの比較

| アルゴリズム | 時間（平均） | 空間 | 安定性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| クイックソート | O(n log n) | O(log n) | × | 普通 |
| マージソート | O(n log n) | O(n) | ○ | 普通 |
| ヒープソート | O(n log n) | O(1) | × | やや難 |
| 挿入ソート | O(n²) | O(1) | ○ | 簡単 |
| 選択ソート | O(n²) | O(1) | × | 超簡単 |

## 最適化アイデア

### 3-way クイックソート

重複要素が多い場合に有効：

```go
func ThreeWayQuickSort(arr []int, low, high int) {
    if low >= high {
        return
    }
    
    lt, gt := low, high
    pivot := arr[low]
    i := low + 1
    
    for i <= gt {
        if arr[i] < pivot {
            arr[lt], arr[i] = arr[i], arr[lt]
            lt++
            i++
        } else if arr[i] > pivot {
            arr[gt], arr[i] = arr[i], arr[gt]
            gt--
        } else {
            i++
        }
    }
    
    ThreeWayQuickSort(arr, low, lt-1)
    ThreeWayQuickSort(arr, gt+1, high)
}

```

### ハイブリッドソート

小規模配列で挿入ソートに切り替え：

```go
func HybridQuickSort(arr []int, low, high int) {
    if high-low < 10 {
        // 小規模なら挿入ソート
        insertionSortRange(arr, low, high)
        return
    }
    
    if low < high {
        pi := partition(arr, low, high)
        HybridQuickSort(arr, low, pi-1)
        HybridQuickSort(arr, pi+1, high)
    }
}

```

## まとめ

### メリット

- 平均的に高速
- メモリ効率良い
- キャッシュ効率良い
- 実装が比較的シンプル

### デメリット

- 最悪O(n²)
- 不安定ソート
- 再帰によるスタック消費
- ピボット選択が重要

### 使うべき時

- 大規模データ
- 平均性能重視
- メモリ制約あり
- 安定性不要

実用的には最も使われるソートアルゴリズムの一つ。多くの言語の標準ライブラリで採用されてる（Go の sort.Slice も内部でクイックソート使用）。