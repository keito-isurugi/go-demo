## 概要

- GoでLCM（最小公倍数）のサンプルコード
- 最小公倍数（Least Common Multiple: LCM）は、2つ以上の自然数の公倍数のうち最小のもの。周期性を扱う問題や分数の通分で頻繁に使用される基本的なアルゴリズム。

### 特徴

- **GCD利用**: 最大公約数を使って効率的に計算
- **簡潔な公式**: LCM(a,b) = (a × b) / GCD(a,b)
- **拡張可能**: 3つ以上の数にも適用可能
- **実用的**: 周期問題、分数計算、スケジューリングに活用
- **オーバーフロー注意**: 積が大きくなりやすい

## 動作原理

基本原理：
「2つの自然数a, bの最小公倍数は、その積をGCD（最大公約数）で割った値に等しい」

### アルゴリズムの流れ

1. 2つの数のGCD（最大公約数）を求める
2. 2つの数の積を計算
3. 積をGCDで割る
4. 結果が最小公倍数

### 具体例

`LCM(12, 18)` の計算過程：

```
Step 1: GCD(12, 18)を求める
        12 = 18 × 0 + 12
        18 = 12 × 1 + 6
        12 = 6 × 2 + 0
        GCD = 6

Step 2: LCMを計算
        LCM = (12 × 18) / 6
        LCM = 216 / 6
        LCM = 36

検証: 36は12の3倍、18の2倍
```

視覚的に理解：
```
12の倍数: 12, 24, [36], 48, 60, 72...
18の倍数: 18, [36], 54, 72, 90, 108...

共通の倍数の最小値 = 36
```

## サンプルコード

### 基本実装（GCD利用）

```go
package main

import "fmt"

// LCM は2つの数の最小公倍数を求める
func LCM(a, b int) int {
    if a == 0 || b == 0 {
        return 0
    }
    // 絶対値で計算
    absA, absB := abs(a), abs(b)
    return (absA * absB) / GCD(absA, absB)
}

// GCD はユークリッドの互除法で最大公約数を求める
func GCD(a, b int) int {
    if b == 0 {
        return a
    }
    return GCD(b, a%b)
}

// abs は絶対値を返す
func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
```

### 複数の数のLCM

```go
// LCMMultiple は複数の数の最小公倍数を求める
func LCMMultiple(numbers ...int) int {
    if len(numbers) == 0 {
        return 0
    }
    if len(numbers) == 1 {
        return abs(numbers[0])
    }

    result := numbers[0]
    for i := 1; i < len(numbers); i++ {
        result = LCM(result, numbers[i])
    }
    return result
}
```

### 素因数分解を使った実装

```go
// LCMByPrimeFactorization は素因数分解を使ってLCMを計算
func LCMByPrimeFactorization(a, b int) int {
    if a == 0 || b == 0 {
        return 0
    }

    absA, absB := abs(a), abs(b)
    factorsA := primeFactors(absA)
    factorsB := primeFactors(absB)

    // 各素因数の最大の指数を取る
    maxFactors := make(map[int]int)
    for prime, count := range factorsA {
        maxFactors[prime] = count
    }
    for prime, count := range factorsB {
        if maxFactors[prime] < count {
            maxFactors[prime] = count
        }
    }

    // 素因数からLCMを計算
    result := 1
    for prime, count := range maxFactors {
        for i := 0; i < count; i++ {
            result *= prime
        }
    }
    return result
}

// primeFactors は素因数分解を行う
func primeFactors(n int) map[int]int {
    factors := make(map[int]int)

    // 2で割れるだけ割る
    for n%2 == 0 {
        factors[2]++
        n = n / 2
    }

    // 奇数の素因数
    for i := 3; i*i <= n; i += 2 {
        for n%i == 0 {
            factors[i]++
            n = n / i
        }
    }

    // nが素数の場合
    if n > 2 {
        factors[n]++
    }

    return factors
}
```

### オーバーフロー対策版

```go
// LCMSafe はオーバーフローを防ぐ実装
func LCMSafe(a, b int) (int, error) {
    if a == 0 || b == 0 {
        return 0, nil
    }

    absA, absB := abs(a), abs(b)
    gcd := GCD(absA, absB)

    // オーバーフローチェック
    // a/gcd * b の形で計算して中間値を小さく保つ
    quotient := absA / gcd

    // quotient * absB がオーバーフローしないかチェック
    if quotient > 0 && absB > math.MaxInt/quotient {
        return 0, fmt.Errorf("LCM overflow for %d and %d", a, b)
    }

    return quotient * absB, nil
}
```

### 分数の通分への応用

```go
// Fraction は分数を表す構造体
type Fraction struct {
    Numerator   int // 分子
    Denominator int // 分母
}

// AddFractions は分数の加算（通分が必要）
func AddFractions(f1, f2 Fraction) Fraction {
    // 分母のLCMを求めて通分
    lcd := LCM(f1.Denominator, f2.Denominator)

    // 各分数を通分
    num1 := f1.Numerator * (lcd / f1.Denominator)
    num2 := f2.Numerator * (lcd / f2.Denominator)

    // 分子を加算
    resultNum := num1 + num2

    // 約分
    gcd := GCD(abs(resultNum), lcd)

    return Fraction{
        Numerator:   resultNum / gcd,
        Denominator: lcd / gcd,
    }
}

// 使用例
// 1/6 + 1/9 = 3/18 + 2/18 = 5/18
```

## 計算量

### 時間計算量

- **基本的なLCM**: O(log min(a,b)) - GCDの計算に依存
- **複数の数のLCM**: O(n × log m) - n個の数、mは最大値
- **素因数分解版**: O(√n) - 素因数分解がボトルネック

### 空間計算量

- **基本版**: O(1) - 定数空間（反復版GCD使用時）
- **素因数分解版**: O(log n) - 素因数の格納

### 数値の大きさ

LCMは積に近い値になるため、オーバーフローに注意：
- LCM(a,b) ≤ a × b
- 互いに素の場合、LCM(a,b) = a × b

## 使いどころ

### 向いてる場面

- 周期的な事象の同期タイミング計算
- 分数の通分（最小公分母の計算）
- 歯車の回転周期
- スケジューリング問題
- 音楽理論での拍子の計算

### 実例：周期的イベントの同期

```go
// EventScheduler は複数の周期的イベントの同期を管理
type EventScheduler struct {
    Events map[string]int // イベント名と周期（分）
}

// NextSyncTime は全イベントが同時に発生する次の時刻を計算
func (s *EventScheduler) NextSyncTime() int {
    periods := []int{}
    for _, period := range s.Events {
        periods = append(periods, period)
    }

    return LCMMultiple(periods...)
}

// GetSchedule は指定時間内のイベント発生を表示
func (s *EventScheduler) GetSchedule(duration int) {
    fmt.Printf("スケジュール（%d分間）:\n", duration)

    for name, period := range s.Events {
        fmt.Printf("%s（%d分ごと）: ", name, period)
        times := []int{}
        for t := period; t <= duration; t += period {
            times = append(times, t)
        }
        fmt.Printf("%v\n", times)
    }

    syncTime := s.NextSyncTime()
    fmt.Printf("\n全イベント同時発生: %d分後\n", syncTime)
}

// 使用例
scheduler := &EventScheduler{
    Events: map[string]int{
        "バックアップ": 12,
        "ログローテート": 18,
        "メトリクス収集": 15,
    },
}
// 同期時刻 = LCM(12, 18, 15) = 180分後
```

### 実例：画面解像度の最適化

```go
// OptimalResolution は複数の画面に最適な共通解像度を計算
func OptimalResolution(widths, heights []int) (int, int) {
    // 幅と高さそれぞれのGCDを求める
    widthGCD := GCDMultiple(widths...)
    heightGCD := GCDMultiple(heights...)

    // 共通の最小単位
    unitWidth := widthGCD
    unitHeight := heightGCD

    // 各画面のグリッド数を計算
    maxGridX := 0
    maxGridY := 0

    for i := range widths {
        gridX := widths[i] / unitWidth
        gridY := heights[i] / unitHeight

        if gridX > maxGridX {
            maxGridX = gridX
        }
        if gridY > maxGridY {
            maxGridY = gridY
        }
    }

    // 最適な共通解像度
    optimalWidth := unitWidth * maxGridX
    optimalHeight := unitHeight * maxGridY

    return optimalWidth, optimalHeight
}
```

## 数学的性質

### 重要な性質

1. **GCDとの関係**: LCM(a,b) × GCD(a,b) = a × b
2. **可換性**: LCM(a,b) = LCM(b,a)
3. **結合性**: LCM(a, LCM(b,c)) = LCM(LCM(a,b), c)
4. **単調性**: a ≤ LCM(a,b) かつ b ≤ LCM(a,b)
5. **互いに素**: GCD(a,b) = 1 ⇒ LCM(a,b) = a × b

### 素因数分解との関係

```
例: 12 = 2² × 3, 18 = 2 × 3²

LCM(12, 18) = 2² × 3² = 36
（各素因数の最大の指数を取る）

GCD(12, 18) = 2 × 3 = 6
（各素因数の最小の指数を取る）
```

## 他のアルゴリズムとの比較

| 手法 | 計算量 | メモリ | 特徴 |
| --- | --- | --- | --- |
| GCD利用 | O(log n) | O(1) | 最も効率的 |
| 素因数分解 | O(√n) | O(log n) | 理解しやすい |
| 倍数列挙 | O(a×b) | O(min(a,b)) | 非効率 |
| ブルートフォース | O(a×b) | O(1) | 単純だが遅い |

## 応用例

### 音楽のポリリズム計算

```go
// PolyrhythmCalculator はポリリズムの周期を計算
type PolyrhythmCalculator struct {
    BeatPatterns map[string]int // 楽器名とビートパターン
}

// CalculateCycle はポリリズムの1周期を計算
func (p *PolyrhythmCalculator) CalculateCycle() {
    patterns := []int{}
    for _, beats := range p.BeatPatterns {
        patterns = append(patterns, beats)
    }

    cycle := LCMMultiple(patterns...)
    fmt.Printf("ポリリズムの1周期: %dビート\n\n", cycle)

    for instrument, beats := range p.BeatPatterns {
        repetitions := cycle / beats
        fmt.Printf("%s: %dビート × %d回\n", instrument, beats, repetitions)
    }
}

// 使用例：3対4のポリリズム
calc := &PolyrhythmCalculator{
    BeatPatterns: map[string]int{
        "右手": 3,
        "左手": 4,
    },
}
// 周期 = LCM(3, 4) = 12ビート
```

### 複数タイマーの同期

```go
// TimerSync は複数のタイマーが同時に鳴るタイミングを計算
func TimerSync(intervals []int) {
    lcm := LCMMultiple(intervals...)

    fmt.Printf("タイマー設定: %v秒\n", intervals)
    fmt.Printf("同時に鳴るタイミング: %d秒後\n\n", lcm)

    // 各タイマーの鳴る回数
    for i, interval := range intervals {
        count := lcm / interval
        fmt.Printf("タイマー%d（%d秒）: %d回鳴る\n", i+1, interval, count)
    }
}
```

## まとめ

### メリット

- GCDを利用した効率的な計算
- 周期性の問題に強い
- 分数計算の基礎
- 実装がシンプル
- 数学的に美しい性質

### デメリット

- 値が急速に大きくなる（オーバーフロー注意）
- 大量の数のLCMは巨大になりやすい
- 素因数分解版は大きな数に弱い

### 使うべき時

- 周期的なイベントの同期
- 分数の通分・計算
- グリッドシステムの設計
- 音楽理論・リズムパターン
- スケジューリング最適化

最小公倍数は、複数の周期を扱う問題で必須のアルゴリズム。GCDとセットで理解することで、より深い数論的な問題にも対応できる基礎となる。