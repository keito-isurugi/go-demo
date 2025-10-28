## 概要

- Goでマージソートを実装
- マージソートは分割統治法を使う安定なソートアルゴリズム。配列を再帰的に半分に分割し、それぞれをソートしてからマージ（併合）することで全体をソートする。

### 特徴

- **分割統治法**: 配列を分割して再帰的にソート
- **安定ソート**: 同じ値の相対順序が保持される
- **予測可能な性能**: 常にO(n log n)の時間計算量
- **追加メモリ必要**: マージ時に一時配列が必要（O(n)）
- **比較ソート**: 要素同士の比較でソート

## 動作原理

基本的な流れ：

1. 配列を中央で2つに分割
2. 左右の部分配列をそれぞれ再帰的にソート
3. ソート済みの2つの配列をマージ
4. 要素が1つになるまで分割を繰り返す

### 具体例

`[38, 27, 43, 3, 9, 82, 10]` のソート過程：

```
初期: [38, 27, 43, 3, 9, 82, 10]

分割フェーズ:
[38, 27, 43, 3, 9, 82, 10]
    ↓
[38, 27, 43] [3, 9, 82, 10]
    ↓
[38] [27, 43] [3, 9] [82, 10]
    ↓
[38] [27] [43] [3] [9] [82] [10]

マージフェーズ:
[38] [27] [43] → [27, 38] [43] → [27, 38, 43]
[3] [9] → [3, 9]
[82] [10] → [10, 82]
[3, 9] [10, 82] → [3, 9, 10, 82]
[27, 38, 43] [3, 9, 10, 82] → [3, 9, 10, 27, 38, 43, 82]

完了: [3, 9, 10, 27, 38, 43, 82]

```

## サンプルコード

### 基本実装

```go
package main

import "fmt"

// MergeSort は配列を昇順にソート
func MergeSort(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    mergeSortHelper(result, 0, len(result)-1)
    return result
}

// mergeSortHelper は再帰的にソート
func mergeSortHelper(arr []int, left, right int) {
    if left < right {
        mid := left + (right-left)/2
        
        // 左右を再帰的にソート
        mergeSortHelper(arr, left, mid)
        mergeSortHelper(arr, mid+1, right)
        
        // マージ
        merge(arr, left, mid, right)
    }
}

// merge は2つのソート済み配列をマージ
func merge(arr []int, left, mid, right int) {
    // 一時配列を作成
    n1 := mid - left + 1
    n2 := right - mid
    
    leftArr := make([]int, n1)
    rightArr := make([]int, n2)
    
    // データをコピー
    for i := 0; i < n1; i++ {
        leftArr[i] = arr[left+i]
    }
    for j := 0; j < n2; j++ {
        rightArr[j] = arr[mid+1+j]
    }
    
    // マージ処理
    i, j, k := 0, 0, left
    
    for i < n1 && j < n2 {
        if leftArr[i] <= rightArr[j] {
            arr[k] = leftArr[i]
            i++
        } else {
            arr[k] = rightArr[j]
            j++
        }
        k++
    }
    
    // 残りをコピー
    for i < n1 {
        arr[k] = leftArr[i]
        i++
        k++
    }
    
    for j < n2 {
        arr[k] = rightArr[j]
        j++
        k++
    }
}

```

### ボトムアップ版（非再帰）

```go
// MergeSortBottomUp は非再帰的にソート
func MergeSortBottomUp(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    n := len(result)
    
    // サブ配列のサイズを2倍ずつ増やす
    for size := 1; size < n; size *= 2 {
        for start := 0; start < n-size; start += 2*size {
            mid := start + size - 1
            end := min(start+2*size-1, n-1)
            
            merge(result, start, mid, end)
        }
    }
    
    return result
}

```

### デバッグ用（ステップ表示）

```go
// MergeSortWithSteps はステップごとの状態を表示
func MergeSortWithSteps(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    fmt.Printf("初期: %v\\n", result)
    mergeSortHelperWithSteps(result, 0, len(result)-1, 0)
    
    return result
}

func mergeSortHelperWithSteps(arr []int, left, right, depth int) {
    if left < right {
        indent := strings.Repeat("  ", depth)
        mid := left + (right-left)/2
        
        fmt.Printf("%s分割: [%d:%d] → [%d:%d] + [%d:%d]\\n", 
            indent, left, right, left, mid, mid+1, right)
        
        mergeSortHelperWithSteps(arr, left, mid, depth+1)
        mergeSortHelperWithSteps(arr, mid+1, right, depth+1)
        
        // マージ前後の状態を表示
        fmt.Printf("%sマージ前: ", indent)
        for i := left; i <= right; i++ {
            fmt.Printf("%d ", arr[i])
        }
        fmt.Printf("\\n")
        
        merge(arr, left, mid, right)
        
        fmt.Printf("%sマージ後: ", indent)
        for i := left; i <= right; i++ {
            fmt.Printf("%d ", arr[i])
        }
        fmt.Printf("\\n\\n")
    }
}

```

## 計算量

### 時間計算量

- 最良・平均・最悪すべて **O(n log n)**
- 入力の状態に関わらず常に同じ

### 空間計算量

- **O(n)** - マージ時の一時配列
- **O(log n)** - 再帰呼び出しのスタック（トップダウン版）
- **O(1)** - スタック使用量（ボトムアップ版）

### 比較と交換

- 比較回数: 常に O(n log n) 回
- データ移動回数: 常に O(n log n) 回

## 使いどころ

### 向いてる場面

- 安定性が必要
- 最悪ケースの性能保証が必要
- 大規模データ
- 連結リスト（追加メモリ不要）
- 外部ソート（ファイルソート）

### 実例：構造体のソート

```go
type Person struct {
    Name string
    Age  int
}

// 年齢で昇順、同じ年齢なら名前順（安定性を活用）
func SortPeople(people []Person) {
    mergeSortPeople(people, 0, len(people)-1)
}

func mergeSortPeople(people []Person, left, right int) {
    if left < right {
        mid := left + (right-left)/2
        
        mergeSortPeople(people, left, mid)
        mergeSortPeople(people, mid+1, right)
        
        mergePeople(people, left, mid, right)
    }
}

func mergePeople(people []Person, left, mid, right int) {
    // 一時配列作成
    leftPart := make([]Person, mid-left+1)
    rightPart := make([]Person, right-mid)
    
    copy(leftPart, people[left:mid+1])
    copy(rightPart, people[mid+1:right+1])
    
    i, j, k := 0, 0, left
    
    for i < len(leftPart) && j < len(rightPart) {
        // 年齢で比較、同じなら名前で比較
        if leftPart[i].Age < rightPart[j].Age ||
           (leftPart[i].Age == rightPart[j].Age && 
            leftPart[i].Name <= rightPart[j].Name) {
            people[k] = leftPart[i]
            i++
        } else {
            people[k] = rightPart[j]
            j++
        }
        k++
    }
    
    // 残りをコピー
    copy(people[k:], leftPart[i:])
    copy(people[k:], rightPart[j:])
}

```

## 他のアルゴリズムとの比較

| アルゴリズム | 時間（平均） | 空間 | 安定性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| マージソート | O(n log n) | O(n) | ○ | 普通 |
| クイックソート | O(n log n) | O(log n) | × | 普通 |
| ヒープソート | O(n log n) | O(1) | × | やや難 |
| 挿入ソート | O(n²) | O(1) | ○ | 簡単 |
| 選択ソート | O(n²) | O(1) | × | 超簡単 |

## 最適化アイデア

### インプレースマージ

メモリ使用量を削減（ただし実装が複雑）：

```go
// インプレースマージ（簡易版）
func inPlaceMerge(arr []int, left, mid, right int) {
    start2 := mid + 1
    
    // 既にマージ済み
    if arr[mid] <= arr[start2] {
        return
    }
    
    for left <= mid && start2 <= right {
        if arr[left] <= arr[start2] {
            left++
        } else {
            value := arr[start2]
            index := start2
            
            // 要素をシフト
            for index != left {
                arr[index] = arr[index-1]
                index--
            }
            arr[left] = value
            
            left++
            mid++
            start2++
        }
    }
}

```

### 小規模配列での最適化

挿入ソートとのハイブリッド：

```go
func HybridMergeSort(arr []int, left, right int) {
    if right-left < 10 {
        // 小規模なら挿入ソート
        insertionSortRange(arr, left, right)
        return
    }
    
    if left < right {
        mid := left + (right-left)/2
        HybridMergeSort(arr, left, mid)
        HybridMergeSort(arr, mid+1, right)
        merge(arr, left, mid, right)
    }
}

```

### 自然なマージソート

既存の順序を活用：

```go
func NaturalMergeSort(arr []int) {
    n := len(arr)
    
    for {
        // 自然な順序の部分列を探す
        runs := findRuns(arr)
        
        if len(runs) <= 1 {
            break
        }
        
        // 隣接するrunをマージ
        mergeRuns(arr, runs)
    }
}

```

## まとめ

### メリット

- 安定ソート
- 予測可能な性能（常にO(n log n)）
- 分割統治の典型例
- 外部ソートに最適
- 並列化しやすい

### デメリット

- 追加メモリ必要（O(n)）
- キャッシュ効率悪い
- 小規模データで遅い

### 使うべき時

- 安定性が必要
- 最悪ケース回避したい
- 大規模データ
- 外部ソート
- 並列処理

実用的には安定性と予測可能な性能が必要な場面で使われる。Javaの`Arrays.sort()`（オブジェクト）やPythonの`sorted()`（Timsort）のベースとなってる重要なアルゴリズム。