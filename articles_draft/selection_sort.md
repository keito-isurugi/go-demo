# 選択ソート（Selection Sort）メモ

## 概要

選択ソートは基本的なソートアルゴリズムの一つ。配列から最小値（または最大値）を見つけて、適切な位置に配置することを繰り返す。

### 特徴

- **シンプルな実装**: 直感的で理解しやすい
- **インプレース**: 追加メモリほぼ不要（O(1)）
- **不安定ソート**: 同じ値の相対順序が保持されない
- **比較ソート**: 要素同士の比較でソート

## 動作原理

基本的な流れ：
1. 未ソート部分から最小値を探す
2. 最小値を未ソート部分の先頭と交換
3. 未ソート部分を1つ後ろにずらす
4. 未ソート部分がなくなるまで繰り返し

### 具体例

`[64, 34, 25, 12, 22, 11, 90]` のソート過程：

```
初期: [64, 34, 25, 12, 22, 11, 90]

Step1: 最小値11を位置0と交換
[11, 34, 25, 12, 22, 64, 90]

Step2: 残りから最小値12を位置1と交換
[11, 12, 25, 34, 22, 64, 90]

Step3: 残りから最小値22を位置2と交換
[11, 12, 22, 34, 25, 64, 90]

Step4: 残りから最小値25を位置3と交換
[11, 12, 22, 25, 34, 64, 90]

Step5: 34は既に正しい位置
Step6: 64も既に正しい位置

完了: [11, 12, 22, 25, 34, 64, 90]
```

## 実装

### 基本実装

```go
package main

import "fmt"

// SelectionSort は配列を昇順にソート
func SelectionSort(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        // i番目以降の最小値を探す
        minIndex := i
        for j := i + 1; j < n; j++ {
            if arr[j] < arr[minIndex] {
                minIndex = j
            }
        }
        
        // 最小値をi番目と交換
        if minIndex != i {
            arr[i], arr[minIndex] = arr[minIndex], arr[i]
        }
    }
}
```

### 降順ソート

```go
// SelectionSortDescending は配列を降順にソート
func SelectionSortDescending(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        // 最大値を探す
        maxIndex := i
        for j := i + 1; j < n; j++ {
            if arr[j] > arr[maxIndex] {
                maxIndex = j
            }
        }
        
        if maxIndex != i {
            arr[i], arr[maxIndex] = arr[maxIndex], arr[i]
        }
    }
}
```

### デバッグ用（ステップ表示）

```go
// SelectionSortWithSteps はステップごとの状態を表示
func SelectionSortWithSteps(arr []int) {
    n := len(arr)
    fmt.Printf("初期: %v\n", arr)
    
    for i := 0; i < n-1; i++ {
        minIndex := i
        fmt.Printf("\nStep %d: pos[%d]を決定\n", i+1, i)
        
        for j := i + 1; j < n; j++ {
            if arr[j] < arr[minIndex] {
                minIndex = j
            }
        }
        
        fmt.Printf("  最小値%d(idx:%d) ← → pos[%d]\n", 
                  arr[minIndex], minIndex, i)
        
        if minIndex != i {
            arr[i], arr[minIndex] = arr[minIndex], arr[i]
            fmt.Printf("  結果: %v\n", arr)
        } else {
            fmt.Printf("  交換なし\n")
        }
    }
}
```

## 計算量

### 時間計算量
- 最良・平均・最悪すべて **O(n²)**
- 入力の状態に関わらず常に同じ

### 空間計算量
- **O(1)** - インプレース

### 比較と交換
- 比較回数: 常に n(n-1)/2 回
- 交換回数: 最大 n-1 回、最小 0 回

## 使いどころ

### 向いてる場面
- 小規模データ（要素数 < 100）
- メモリ制約が厳しい
- 交換コストが高い（大きな構造体など）

### 実例：構造体のソート

```go
type Student struct {
    Name  string
    Score int
}

// スコアで降順ソート
func SortStudentsByScore(students []Student) {
    n := len(students)
    
    for i := 0; i < n-1; i++ {
        maxIndex := i
        for j := i + 1; j < n; j++ {
            if students[j].Score > students[maxIndex].Score {
                maxIndex = j
            }
        }
        
        if maxIndex != i {
            students[i], students[maxIndex] = students[maxIndex], students[i]
        }
    }
}
```

## 他のアルゴリズムとの比較

| アルゴリズム | 時間（平均） | 空間 | 安定性 | 実装難易度 |
|------------|------------|------|--------|-----------|
| 選択ソート | O(n²) | O(1) | × | 超簡単 |
| バブルソート | O(n²) | O(1) | ○ | 超簡単 |
| 挿入ソート | O(n²) | O(1) | ○ | 簡単 |
| クイックソート | O(n log n) | O(log n) | × | 普通 |
| マージソート | O(n log n) | O(n) | ○ | 普通 |

## 最適化アイデア

### 双方向選択ソート

一度のパスで最小値と最大値を同時に配置：

```go
func BidirectionalSelectionSort(arr []int) {
    left := 0
    right := len(arr) - 1
    
    for left < right {
        minIndex := left
        maxIndex := right
        
        // 最小値と最大値を同時に探す
        for i := left; i <= right; i++ {
            if arr[i] < arr[minIndex] {
                minIndex = i
            }
            if arr[i] > arr[maxIndex] {
                maxIndex = i
            }
        }
        
        // 最小値を左端へ
        if minIndex != left {
            arr[left], arr[minIndex] = arr[minIndex], arr[left]
        }
        
        // maxIndexの調整（左端と交換した場合）
        if maxIndex == left {
            maxIndex = minIndex
        }
        
        // 最大値を右端へ
        if maxIndex != right {
            arr[right], arr[maxIndex] = arr[maxIndex], arr[right]
        }
        
        left++
        right--
    }
}
```

### 早期終了の検討

本質的に早期終了は難しいが、一応の最適化：

```go
func OptimizedSelectionSort(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        minIndex := i
        isSorted := true
        
        for j := i + 1; j < n; j++ {
            if arr[j] < arr[minIndex] {
                minIndex = j
            }
            // 隣接要素が逆順なら未ソート
            if j > i+1 && arr[j] < arr[j-1] {
                isSorted = false
            }
        }
        
        if minIndex != i {
            arr[i], arr[minIndex] = arr[minIndex], arr[i]
        }
        
        // 残りがソート済みなら終了
        if isSorted && minIndex == i {
            break
        }
    }
}
```

## まとめ

### メリット
- 実装が超シンプル
- メモリ効率良い
- 交換回数少ない
- 動作が予測可能

### デメリット
- O(n²)で遅い
- 不安定ソート
- 最適化の余地が少ない

### 使うべき時
- データ少ない（< 100要素）
- メモリ厳しい
- 学習目的
- 交換コスト高い

実用的には限定的だが、アルゴリズムの基礎として重要。大規模データならクイックソートやマージソート使った方がいい。