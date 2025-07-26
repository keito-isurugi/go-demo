## 概要

- Goで挿入ソートを実装
- 挿入ソートは基本的なソートアルゴリズムの一つ。トランプの手札を整理する際の動作に似ており、未ソート部分から要素を一つずつ取り出し、ソート済み部分の適切な位置に挿入する。

### 特徴

- **シンプルな実装**: 直感的で理解しやすい
- **インプレース**: 追加メモリほぼ不要（O(1)）
- **安定ソート**: 同じ値の相対順序が保持される
- **適応的**: 部分的にソート済みのデータに効率的
- **比較ソート**: 要素同士の比較でソート

## 動作原理

基本的な流れ：

1. 配列の2番目の要素から開始（1番目は既にソート済みとみなす）
2. 現在の要素を取り出し、ソート済み部分の適切な位置を探す
3. ソート済み部分の要素を右にシフトしながら挿入位置を見つける
4. 適切な位置に要素を挿入
5. すべての要素が処理されるまで繰り返す

### 具体例

`[64, 34, 25, 12, 22, 11, 90]` のソート過程：

```
初期: [64, 34, 25, 12, 22, 11, 90]

Step1: 34を先頭部分[64]の適切な位置に挿入
[34, 64, 25, 12, 22, 11, 90]

Step2: 25をソート済み部分[34, 64]の適切な位置に挿入
[25, 34, 64, 12, 22, 11, 90]

Step3: 12をソート済み部分[25, 34, 64]の適切な位置に挿入
[12, 25, 34, 64, 22, 11, 90]

Step4: 22をソート済み部分[12, 25, 34, 64]の適切な位置に挿入
[12, 22, 25, 34, 64, 11, 90]

Step5: 11をソート済み部分[12, 22, 25, 34, 64]の適切な位置に挿入
[11, 12, 22, 25, 34, 64, 90]

Step6: 90をソート済み部分[11, 12, 22, 25, 34, 64]の適切な位置に挿入
[11, 12, 22, 25, 34, 64, 90]

```

## サンプルコード

### 基本実装

```go
package main

import "fmt"

// InsertionSort は配列を昇順にソート
func InsertionSort(arr []int) {
    n := len(arr)

    for i := 1; i < n; i++ {
        // 現在の要素を保存
        key := arr[i]
        j := i - 1

        // keyより大きい要素を右にシフト
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }

        // 適切な位置にkeyを挿入
        arr[j+1] = key
    }
}

```

### 降順ソート

```go
// InsertionSortDescending は配列を降順にソート
func InsertionSortDescending(arr []int) {
    n := len(arr)

    for i := 1; i < n; i++ {
        key := arr[i]
        j := i - 1

        // keyより小さい要素を右にシフト
        for j >= 0 && arr[j] < key {
            arr[j+1] = arr[j]
            j--
        }

        arr[j+1] = key
    }
}

```

### デバッグ用（ステップ表示）

```go
// InsertionSortWithSteps はステップごとの状態を表示
func InsertionSortWithSteps(arr []int) {
    n := len(arr)
    fmt.Printf("初期: %v\\n", arr)

    for i := 1; i < n; i++ {
        key := arr[i]
        j := i - 1

        fmt.Printf("\\nStep %d: 要素%d を適切な位置に挿入\\n", i, key)

        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }

        arr[j+1] = key
        fmt.Printf("  結果: %v\\n", arr)
    }
}

```

## 計算量

### 時間計算量

- 最良：**O(n)** - 既にソート済みの配列
- 平均・最悪：**O(n²)** - ランダムな配列や逆順の配列

### 空間計算量

- **O(1)** - インプレース

### 比較と交換

- 比較回数: 最大 n(n-1)/2 回、最小 n-1 回
- 交換回数: 最大 n(n-1)/2 回、最小 0 回

## 使いどころ

### 向いてる場面

- 小規模データ（要素数 < 100）
- ほぼソート済みのデータ
- オンラインアルゴリズム（データが逐次到着）
- 安定性が必要

### 実例：構造体のソート

```go
type Student struct {
    Name  string
    Score int
}

// スコアで昇順ソート
func SortStudentsByScore(students []Student) {
    n := len(students)

    for i := 1; i < n; i++ {
        key := students[i]
        j := i - 1

        for j >= 0 && students[j].Score > key.Score {
            students[j+1] = students[j]
            j--
        }

        students[j+1] = key
    }
}

```

## 他のアルゴリズムとの比較

| アルゴリズム | 時間（平均） | 空間 | 安定性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| 挿入ソート | O(n²) | O(1) | ○ | 簡単 |
| 選択ソート | O(n²) | O(1) | × | 超簡単 |
| バブルソート | O(n²) | O(1) | ○ | 超簡単 |
| クイックソート | O(n log n) | O(log n) | × | 普通 |
| マージソート | O(n log n) | O(n) | ○ | 普通 |

## 最適化アイデア

### バイナリサーチ版

挿入位置を二分探索で見つけることで比較回数を削減：

```go
func BinaryInsertionSort(arr []int) {
    n := len(arr)

    for i := 1; i < n; i++ {
        key := arr[i]
        
        // 挿入位置を二分探索で見つける
        left, right := 0, i
        for left < right {
            mid := (left + right) / 2
            if arr[mid] <= key {
                left = mid + 1
            } else {
                right = mid
            }
        }

        // 要素をシフト
        for j := i; j > left; j-- {
            arr[j] = arr[j-1]
        }

        arr[left] = key
    }
}

```

### 早期終了の最適化

```go
func OptimizedInsertionSort(arr []int) {
    n := len(arr)

    for i := 1; i < n; i++ {
        key := arr[i]
        
        // 既に正しい位置にある場合は何もしない
        if arr[i-1] <= key {
            continue
        }

        j := i - 1
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }

        arr[j+1] = key
    }
}

```

## まとめ

### メリット

- 実装がシンプル
- メモリ効率良い
- 安定ソート
- 適応的（部分的ソート済みで高速）
- オンライン対応

### デメリット

- 最悪時間計算量 O(n²)
- 大規模データには不向き

### 使うべき時

- データ少ない（< 100要素）
- ほぼソート済み
- 安定性が必要
- オンライン処理

実用的には小規模データで有効。大規模データならクイックソートやマージソート使った方がいい。