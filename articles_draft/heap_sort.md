## 概要

- Goでヒープソートを実装
- ヒープソートはヒープデータ構造を使うソートアルゴリズム。配列を最大ヒープに変換し、最大値を順次取り出すことでソートする。

### 特徴

- **ヒープ利用**: 完全二分木の性質を持つヒープを使用
- **インプレース**: 追加メモリほぼ不要（O(1)）
- **不安定ソート**: 同じ値の相対順序が保持されない
- **予測可能な性能**: 常にO(n log n)の時間計算量
- **比較ソート**: 要素同士の比較でソート

## 動作原理

基本的な流れ：

1. 配列を最大ヒープに変換（ヒープ構築）
2. ヒープのルート（最大値）を配列の末尾と交換
3. ヒープサイズを1つ減らす
4. ルートから再度ヒープ化
5. ヒープが空になるまで2-4を繰り返す

### 具体例

`[12, 11, 13, 5, 6, 7]` のソート過程：

```
初期: [12, 11, 13, 5, 6, 7]

ヒープ構築フェーズ:
[12, 11, 13, 5, 6, 7] → 最大ヒープ化
[13, 11, 12, 5, 6, 7] (最大ヒープ完成)

ソートフェーズ:
Step1: 最大値13を末尾に移動
[7, 11, 12, 5, 6] | [13]
再ヒープ化: [12, 11, 7, 5, 6] | [13]

Step2: 最大値12を末尾に移動
[6, 11, 7, 5] | [12, 13]
再ヒープ化: [11, 6, 7, 5] | [12, 13]

Step3: 最大値11を末尾に移動
[5, 6, 7] | [11, 12, 13]
再ヒープ化: [7, 6, 5] | [11, 12, 13]

... 続く

完了: [5, 6, 7, 11, 12, 13]

```

## サンプルコード

### 基本実装

```go
package main

import "fmt"

// HeapSort は配列を昇順にソート
func HeapSort(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    heapSortHelper(result)
    return result
}

// heapSortHelper は実際のヒープソート処理
func heapSortHelper(arr []int) {
    n := len(arr)
    
    // ヒープを構築（最大ヒープ）
    for i := n/2 - 1; i >= 0; i-- {
        heapify(arr, n, i)
    }
    
    // 一つずつ要素をヒープから取り出す
    for i := n - 1; i > 0; i-- {
        // ルート（最大値）を末尾に移動
        arr[0], arr[i] = arr[i], arr[0]
        
        // 縮小されたヒープで再ヒープ化
        heapify(arr, i, 0)
    }
}

// heapify は部分木をヒープ化
func heapify(arr []int, heapSize, rootIndex int) {
    largest := rootIndex
    leftChild := 2*rootIndex + 1
    rightChild := 2*rootIndex + 2
    
    // 左の子が存在し、ルートより大きい場合
    if leftChild < heapSize && arr[leftChild] > arr[largest] {
        largest = leftChild
    }
    
    // 右の子が存在し、現在の最大値より大きい場合
    if rightChild < heapSize && arr[rightChild] > arr[largest] {
        largest = rightChild
    }
    
    // ルートが最大値でない場合
    if largest != rootIndex {
        arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
        
        // 影響を受けた部分木を再帰的にヒープ化
        heapify(arr, heapSize, largest)
    }
}

```

### ヒープ操作関数

```go
// BuildMaxHeap は配列を最大ヒープに変換
func BuildMaxHeap(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    n := len(result)
    for i := n/2 - 1; i >= 0; i-- {
        heapify(result, n, i)
    }
    
    return result
}

// IsMaxHeap は配列が最大ヒープかどうかを確認
func IsMaxHeap(arr []int) bool {
    n := len(arr)
    
    for i := 0; i < n/2; i++ {
        leftChild := 2*i + 1
        rightChild := 2*i + 2
        
        if leftChild < n && arr[i] < arr[leftChild] {
            return false
        }
        
        if rightChild < n && arr[i] < arr[rightChild] {
            return false
        }
    }
    
    return true
}

```

### デバッグ用（ステップ表示）

```go
// HeapSortWithSteps はステップごとの状態を表示
func HeapSortWithSteps(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    fmt.Printf("初期: %v\\n", result)
    
    n := len(result)
    
    // ヒープ構築フェーズ
    fmt.Println("\\n--- ヒープ構築フェーズ ---")
    for i := n/2 - 1; i >= 0; i-- {
        fmt.Printf("ヒープ化: インデックス %d から\\n", i)
        heapifyWithSteps(result, n, i)
        fmt.Printf("結果: %v\\n\\n", result)
    }
    
    fmt.Println("最大ヒープ完成:", result)
    
    // ソートフェーズ
    fmt.Println("\\n--- ソートフェーズ ---")
    for i := n - 1; i > 0; i-- {
        fmt.Printf("ステップ %d: 最大値 %d を位置 %d に移動\\n", n-i, result[0], i)
        
        result[0], result[i] = result[i], result[0]
        fmt.Printf("交換後: %v\\n", result)
        
        heapifyWithSteps(result, i, 0)
        fmt.Printf("ヒープ化後: %v\\n\\n", result)
    }
    
    return result
}

```

## 計算量

### 時間計算量

- **ヒープ構築**: O(n)
- **ヒープソート全体**: O(n log n)
- 最良・平均・最悪すべて **O(n log n)**

### 空間計算量

- **O(1)** - インプレース（再帰スタック除く）
- **O(log n)** - 再帰呼び出しのスタック

### 比較と交換

- 比較回数: 約 2n log n 回
- 交換回数: 約 n log n 回

## 使いどころ

### 向いてる場面

- メモリ制約が厳しい環境
- 最悪ケースの性能保証が必要
- 優先度付きキューの実装
- 部分ソート（k個の最大値が欲しい場合）

### 実例：優先度付きキュー

```go
type PriorityQueue struct {
    heap []int
}

// Insert は要素を優先度付きキューに追加
func (pq *PriorityQueue) Insert(value int) {
    pq.heap = append(pq.heap, value)
    
    // ボトムアップでヒープ性質を回復
    index := len(pq.heap) - 1
    for index > 0 {
        parent := (index - 1) / 2
        
        if pq.heap[parent] >= pq.heap[index] {
            break
        }
        
        pq.heap[parent], pq.heap[index] = pq.heap[index], pq.heap[parent]
        index = parent
    }
}

// ExtractMax は最大値を取り出し
func (pq *PriorityQueue) ExtractMax() int {
    if len(pq.heap) == 0 {
        return 0
    }
    
    max := pq.heap[0]
    
    // 最後の要素をルートに移動
    pq.heap[0] = pq.heap[len(pq.heap)-1]
    pq.heap = pq.heap[:len(pq.heap)-1]
    
    // ヒープ性質を回復
    if len(pq.heap) > 0 {
        heapify(pq.heap, len(pq.heap), 0)
    }
    
    return max
}

```

## 他のアルゴリズムとの比較

| アルゴリズム | 時間（平均） | 空間 | 安定性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| ヒープソート | O(n log n) | O(1) | × | やや難 |
| クイックソート | O(n log n) | O(log n) | × | 普通 |
| マージソート | O(n log n) | O(n) | ○ | 普通 |
| 挿入ソート | O(n²) | O(1) | ○ | 簡単 |
| 選択ソート | O(n²) | O(1) | × | 超簡単 |

## 最適化アイデア

### イテレーティブなヒープ化

再帰を避けてスタック使用量を削減：

```go
func heapifyIterative(arr []int, heapSize, rootIndex int) {
    for {
        largest := rootIndex
        leftChild := 2*rootIndex + 1
        rightChild := 2*rootIndex + 2
        
        if leftChild < heapSize && arr[leftChild] > arr[largest] {
            largest = leftChild
        }
        
        if rightChild < heapSize && arr[rightChild] > arr[largest] {
            largest = rightChild
        }
        
        if largest == rootIndex {
            break
        }
        
        arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
        rootIndex = largest
    }
}

```

### k個の最大値を取得

全体をソートせずに最大のk個を取得：

```go
func GetTopK(arr []int, k int) []int {
    if k >= len(arr) {
        return HeapSort(arr)
    }
    
    heap := BuildMaxHeap(arr)
    result := make([]int, k)
    
    for i := 0; i < k; i++ {
        result[i] = heap[0]
        
        // 最大値を取り出し
        heap[0] = heap[len(heap)-1]
        heap = heap[:len(heap)-1]
        
        if len(heap) > 0 {
            heapify(heap, len(heap), 0)
        }
    }
    
    return result
}

```

### ボトムアップヒープ構築

効率的なヒープ構築：

```go
func BuildMaxHeapBottomUp(arr []int) {
    n := len(arr)
    
    // 最後の非葉ノードから開始
    for i := n/2 - 1; i >= 0; i-- {
        // 各ノードを適切な位置に下ろす
        siftDown(arr, i, n-1)
    }
}

func siftDown(arr []int, start, end int) {
    root := start
    
    for root*2+1 <= end {
        child := root*2 + 1
        
        // 右の子が存在し、左の子より大きい
        if child+1 <= end && arr[child] < arr[child+1] {
            child++
        }
        
        if arr[root] < arr[child] {
            arr[root], arr[child] = arr[child], arr[root]
            root = child
        } else {
            break
        }
    }
}

```

## まとめ

### メリット

- メモリ効率が良い（O(1)）
- 予測可能な性能（常にO(n log n)）
- 優先度付きキューに応用可能
- 部分ソートが効率的

### デメリット

- 不安定ソート
- キャッシュ効率が悪い
- 実装がやや複雑
- 実際の実行時間は他より遅い傾向

### 使うべき時

- メモリ制約厳しい
- 最悪ケース回避したい
- 優先度付きキュー必要
- 部分ソートで十分

実用的には組み込みシステムなどメモリ制約が厳しい環境や、優先度付きキューの実装で使われる。多くの言語の標準ライブラリにはヒープ（priority queue）として提供されている重要なデータ構造。