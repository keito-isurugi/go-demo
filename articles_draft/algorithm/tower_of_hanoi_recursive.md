## 概要

- Goでハノイの塔を再帰的に行うサンプルコードを書いてみた
- ハノイの塔は3本の棒と複数の円盤を使う数学パズル。全ての円盤を別の棒に移動させる問題で、再帰アルゴリズムの代表例。

### 特徴

- **分割統治法**: 問題を小さな部分問題に分割
- **再帰の美しさ**: 複雑な問題をシンプルなコードで解決
- **指数的移動回数**: 2^n - 1回の最小移動回数
- **数学的証明可能**: 最適解が数学的に証明済み
- **視覚的理解**: 動きをステップごとに追跡可能

## 動作原理

基本的な流れ：

1. n枚のディスクを元の棒から目的の棒へ移動
2. 補助の棒を使って一時的にディスクを置く
3. 大きいディスクの上に小さいディスクしか置けない
4. 一度に動かせるのは1枚のみ

### 再帰的解法

n枚のディスクを棒Aから棒Cへ移動する手順：

1. 上のn-1枚を棒Aから棒Bへ（棒C経由）
2. 最大のディスクを棒Aから棒Cへ
3. n-1枚を棒Bから棒Cへ（棒A経由）

### 具体例

3枚のディスクの移動過程：

```
初期状態:        移動1後:         移動2後:
  1                              
  2               2              3
  3               3              
-----           -----           -----
A B C           A B C           A B C
                    1            2 1

移動3後:         移動4後:         移動5後:
                                 1
  3               2              2
                  3              3
-----           -----           -----
A B C           A B C           A B C
  1 2               1              

移動6後:         移動7後（完了）:
                    1
                    2
1                   3
-----           -----
A B C           A B C
    2 3

総移動回数: 7回 = 2^3 - 1
```

## サンプルコード

### 基本実装

```go
package recursion

import "fmt"

// HanoiMove は1回の移動を表す構造体
type HanoiMove struct {
    From string
    To   string
    Disk int
}

// TowerOfHanoi はハノイの塔の移動手順を返す
func TowerOfHanoi(n int, from, to, aux string) []HanoiMove {
    moves := []HanoiMove{}
    hanoiHelper(n, from, to, aux, &moves)
    return moves
}

func hanoiHelper(n int, from, to, aux string, moves *[]HanoiMove) {
    // ベースケース: 1枚の場合は直接移動
    if n == 1 {
        *moves = append(*moves, HanoiMove{From: from, To: to, Disk: n})
        return
    }
    
    // ステップ1: n-1枚を補助の棒へ
    hanoiHelper(n-1, from, aux, to, moves)
    
    // ステップ2: 最大のディスクを目的地へ
    *moves = append(*moves, HanoiMove{From: from, To: to, Disk: n})
    
    // ステップ3: n-1枚を目的地へ
    hanoiHelper(n-1, aux, to, from, moves)
}

```

### ステップ表示版

```go
// TowerOfHanoiWithSteps は再帰呼び出しの過程を表示
func TowerOfHanoiWithSteps(n int, from, to, aux string, depth int) {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    
    fmt.Printf("%sHanoi(%d, %s, %s, %s) を呼び出し\\n", indent, n, from, to, aux)
    
    if n == 1 {
        fmt.Printf("%s  ディスク1を%sから%sへ移動\\n", indent, from, to)
        return
    }
    
    fmt.Printf("%s  ステップ1: %d枚のディスクを%sから%sへ（%s経由）\\n", 
               indent, n-1, from, aux, to)
    TowerOfHanoiWithSteps(n-1, from, aux, to, depth+1)
    
    fmt.Printf("%s  ステップ2: ディスク%dを%sから%sへ移動\\n", 
               indent, n, from, to)
    
    fmt.Printf("%s  ステップ3: %d枚のディスクを%sから%sへ（%s経由）\\n", 
               indent, n-1, aux, to, from)
    TowerOfHanoiWithSteps(n-1, aux, to, from, depth+1)
}

```

### シミュレーション版

```go
// SimulateTowerOfHanoi は移動過程を視覚的に表示
func SimulateTowerOfHanoi(n int) {
    // 3本の棒を初期化
    towers := map[string][]int{
        "A": make([]int, 0),
        "B": make([]int, 0),
        "C": make([]int, 0),
    }
    
    // 初期状態: 全てのディスクを棒Aに配置
    for i := n; i >= 1; i-- {
        towers["A"] = append(towers["A"], i)
    }
    
    fmt.Printf("初期状態（%d枚のディスク）:\\n", n)
    PrintHanoiState(towers)
    
    // 移動手順を取得
    moves := TowerOfHanoi(n, "A", "C", "B")
    
    // 各移動を実行
    for i, move := range moves {
        fmt.Printf("移動 %d: ディスク%dを%sから%sへ\\n", 
                   i+1, move.Disk, move.From, move.To)
        
        // ディスクを移動
        fromTower := towers[move.From]
        disk := fromTower[len(fromTower)-1]
        towers[move.From] = fromTower[:len(fromTower)-1]
        towers[move.To] = append(towers[move.To], disk)
        
        PrintHanoiState(towers)
    }
    
    fmt.Printf("\\n完了！ 総移動回数: %d\\n", len(moves))
}

// PrintHanoiState は塔の状態を表示
func PrintHanoiState(towers map[string][]int) {
    maxHeight := 0
    for _, tower := range towers {
        if len(tower) > maxHeight {
            maxHeight = len(tower)
        }
    }
    
    // 上から順に表示
    for level := maxHeight - 1; level >= 0; level-- {
        for _, peg := range []string{"A", "B", "C"} {
            if level < len(towers[peg]) {
                fmt.Printf("  %d  ", towers[peg][level])
            } else {
                fmt.Print("  |  ")
            }
        }
        fmt.Println()
    }
    fmt.Println("-----+-----+-----")
    fmt.Println("  A     B     C  ")
    fmt.Println()
}

```

### 移動回数の計算

```go
// CountHanoiMoves はn枚のディスクの最小移動回数を計算
func CountHanoiMoves(n int) int {
    return (1 << n) - 1  // 2^n - 1
}

// 使用例
func main() {
    for n := 1; n <= 10; n++ {
        moves := CountHanoiMoves(n)
        fmt.Printf("%2d枚: %d回\\n", n, moves)
    }
    // 出力:
    //  1枚: 1回
    //  2枚: 3回
    //  3枚: 7回
    //  4枚: 15回
    //  5枚: 31回
}

```

### 反復版（比較用）

```go
// TowerOfHanoiIterative は反復による実装
func TowerOfHanoiIterative(n int) []HanoiMove {
    moves := []HanoiMove{}
    
    // 偶数と奇数で移動パターンが異なる
    if n%2 == 0 {
        moves = iterativeHanoi(n, "A", "B", "C")
    } else {
        moves = iterativeHanoi(n, "A", "C", "B")
    }
    
    return moves
}

func iterativeHanoi(n int, source, auxiliary, destination string) []HanoiMove {
    moves := []HanoiMove{}
    totalMoves := (1 << n) - 1
    
    // 各棒の状態を管理
    towers := map[string][]int{
        source:      make([]int, 0),
        auxiliary:   make([]int, 0),
        destination: make([]int, 0),
    }
    
    // 初期配置
    for i := n; i >= 1; i-- {
        towers[source] = append(towers[source], i)
    }
    
    pegs := []string{source, auxiliary, destination}
    
    // ビット演算を使った移動パターン
    for i := 1; i <= totalMoves; i++ {
        from := pegs[(i&(i-1))%3]
        to := pegs[((i|(i-1))+1)%3]
        
        // 小さいディスクを大きいディスクの上に移動
        if canMove(towers, from, to) {
            moveDisc(towers, from, to, &moves)
        } else {
            moveDisc(towers, to, from, &moves)
        }
    }
    
    return moves
}

```

## 計算量

### 時間計算量

- **移動回数**: O(2^n) - 正確には2^n - 1回
- **再帰呼び出し**: O(2^n) - 各ノードで2つの再帰呼び出し

### 空間計算量

- **再帰版**: O(n) - 再帰スタックの深さ
- **反復版**: O(n) - ディスクの状態を保持

### ディスク枚数と移動回数

| ディスク枚数 | 移動回数 | 計算式 |
| --- | --- | --- |
| 1 | 1 | 2^1 - 1 |
| 2 | 3 | 2^2 - 1 |
| 3 | 7 | 2^3 - 1 |
| 4 | 15 | 2^4 - 1 |
| 5 | 31 | 2^5 - 1 |
| 10 | 1,023 | 2^10 - 1 |
| 20 | 1,048,575 | 2^20 - 1 |
| 64 | 約1,844京 | 2^64 - 1 |

## 使いどころ

### 向いてる場面

- 再帰アルゴリズムの教育
- 分割統治法の理解
- 問題解決思考の訓練
- アルゴリズム設計の練習

### 実例：パズルゲーム実装

```go
// HanoiGame はハノイの塔ゲームの実装
type HanoiGame struct {
    Towers    map[string][]int
    NumDisks  int
    MoveCount int
    History   []HanoiMove
}

func NewHanoiGame(n int) *HanoiGame {
    game := &HanoiGame{
        Towers:    make(map[string][]int),
        NumDisks:  n,
        MoveCount: 0,
        History:   []HanoiMove{},
    }
    
    // 初期配置
    game.Towers["A"] = make([]int, 0)
    game.Towers["B"] = make([]int, 0)
    game.Towers["C"] = make([]int, 0)
    
    for i := n; i >= 1; i-- {
        game.Towers["A"] = append(game.Towers["A"], i)
    }
    
    return game
}

// Move はディスクを移動
func (g *HanoiGame) Move(from, to string) error {
    fromTower := g.Towers[from]
    toTower := g.Towers[to]
    
    // 移動可能チェック
    if len(fromTower) == 0 {
        return fmt.Errorf("棒%sにディスクがありません", from)
    }
    
    disk := fromTower[len(fromTower)-1]
    if len(toTower) > 0 && toTower[len(toTower)-1] < disk {
        return fmt.Errorf("大きいディスクを小さいディスクの上に置けません")
    }
    
    // 移動実行
    g.Towers[from] = fromTower[:len(fromTower)-1]
    g.Towers[to] = append(toTower, disk)
    g.MoveCount++
    g.History = append(g.History, HanoiMove{From: from, To: to, Disk: disk})
    
    return nil
}

// IsComplete はパズルが完成したか判定
func (g *HanoiGame) IsComplete() bool {
    return len(g.Towers["C"]) == g.NumDisks
}

```

## 他の実装との比較

| 実装方法 | 時間計算量 | 空間計算量 | 可読性 | 実装難易度 |
| --- | --- | --- | --- | --- |
| 再帰 | O(2^n) | O(n) | ◎ | 易 |
| 反復（ビット演算） | O(2^n) | O(n) | △ | 難 |
| 反復（スタック） | O(2^n) | O(n) | ○ | 中 |
| 並列処理 | O(2^n) | O(n) | △ | 難 |

## 最適化アイデア

### パフォーマンス測定

```go
import "time"

func BenchmarkHanoi(n int) {
    fmt.Printf("n = %d での比較:\\n", n)
    
    // 再帰版
    start := time.Now()
    recursiveMoves := TowerOfHanoi(n, "A", "C", "B")
    recursiveTime := time.Since(start)
    
    // 反復版
    start = time.Now()
    iterativeMoves := TowerOfHanoiIterative(n)
    iterativeTime := time.Since(start)
    
    fmt.Printf("再帰版: %v (%d moves)\\n", recursiveTime, len(recursiveMoves))
    fmt.Printf("反復版: %v (%d moves)\\n", iterativeTime, len(iterativeMoves))
}

```

### 並列処理版

```go
// ParallelHanoi は並列処理を使った実装（概念的）
func ParallelHanoi(n int) []HanoiMove {
    // 左右の部分問題を並列実行
    ch1 := make(chan []HanoiMove)
    ch2 := make(chan []HanoiMove)
    
    go func() {
        // n-1枚を補助棒へ
        ch1 <- TowerOfHanoi(n-1, "A", "B", "C")
    }()
    
    go func() {
        // n-1枚を目的棒へ
        ch2 <- TowerOfHanoi(n-1, "B", "C", "A")
    }()
    
    // 結果を結合
    moves1 := <-ch1
    middleMove := HanoiMove{From: "A", To: "C", Disk: n}
    moves2 := <-ch2
    
    result := append(moves1, middleMove)
    result = append(result, moves2...)
    
    return result
}

```

### メモリ効率版

```go
// HanoiGenerator は必要に応じて移動を生成
func HanoiGenerator(n int) <-chan HanoiMove {
    ch := make(chan HanoiMove)
    
    go func() {
        defer close(ch)
        generateMoves(n, "A", "C", "B", ch)
    }()
    
    return ch
}

func generateMoves(n int, from, to, aux string, ch chan<- HanoiMove) {
    if n == 1 {
        ch <- HanoiMove{From: from, To: to, Disk: n}
        return
    }
    
    generateMoves(n-1, from, aux, to, ch)
    ch <- HanoiMove{From: from, To: to, Disk: n}
    generateMoves(n-1, aux, to, from, ch)
}

// 使用例
func main() {
    for move := range HanoiGenerator(3) {
        fmt.Printf("移動: ディスク%dを%sから%sへ\\n", 
                   move.Disk, move.From, move.To)
    }
}

```

## まとめ

### メリット

- エレガントな再帰解法
- 数学的に最適解が保証
- 分割統治法の良い例
- 実装がシンプル

### デメリット

- 指数的な実行時間
- 大きなnでは実用不可
- 再帰スタックの制限
- 並列化が困難

### 使うべき時

- アルゴリズム教育
- 再帰の理解
- パズルゲーム実装
- 問題解決思考の訓練

### 学習ポイント

- **再帰的思考**: 大きな問題を小さな問題に分割
- **数学的帰納法**: ベースケースと再帰ケース
- **計算量の理解**: 指数的増加の実感
- **最適性の証明**: 数学的に最小手数が証明可能

ハノイの塔は再帰アルゴリズムの美しさと威力を示す最良の例。シンプルなルールから複雑な動作が生まれ、数学的な美しさとプログラミングの実用性を兼ね備えた問題。