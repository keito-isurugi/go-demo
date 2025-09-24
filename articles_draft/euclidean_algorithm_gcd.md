## 概要

- Goでユークリッドの互除法を実装
- ユークリッドの互除法は2つの自然数の最大公約数（Greatest Common Divisor: GCD）を効率的に求めるアルゴリズム。2300年以上前から使われている最古のアルゴリズムの一つ。

### 特徴

- **効率的**: 最悪でもO(log min(a,b))の計算量
- **シンプル**: 除算と剰余の繰り返しだけ
- **汎用的**: 最小公倍数の計算にも応用可能
- **拡張可能**: 拡張版で一次不定方程式も解ける
- **実装容易**: 再帰でも反復でも簡潔に書ける

## 動作原理

基本原理：
「2つの自然数a, bについて、aをbで割った余りをrとすると、GCD(a, b) = GCD(b, r)が成り立つ」

### アルゴリズムの流れ

1. 大きい数を小さい数で割る
2. 余りを求める
3. 小さい数と余りで同じ処理を繰り返す
4. 余りが0になったら、その時の除数が最大公約数

### 具体例

`GCD(48, 18)` の計算過程：

```
Step 1: 48 = 18 × 2 + 12  (余り12)
Step 2: 18 = 12 × 1 + 6   (余り6)
Step 3: 12 = 6 × 2 + 0    (余り0)

最大公約数: 6
```

視覚的に理解：
```
48と18の最大公約数を求める
48 |||||||||||||||||||||||||||||||||||||||||||||||||
18 ||||||||||||||||||

48を18で割る → 商2、余り12
18 ||||||||||||||||||
12 ||||||||||||

18を12で割る → 商1、余り6
12 ||||||||||||
6  ||||||

12を6で割る → 商2、余り0
6  ||||||
0

GCD = 6
```

## サンプルコード

### 基本実装（再帰版）

```go
package main

import "fmt"

// GCD はユークリッドの互除法で最大公約数を求める
func GCD(a, b int) int {
    if b == 0 {
        return a
    }
    return GCD(b, a%b)
}
```

### 反復版

```go
// GCDIterative は反復版のユークリッドの互除法
func GCDIterative(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}
```

### デバッグ用（ステップ表示）

```go
// GCDWithSteps はステップごとの計算過程を表示
func GCDWithSteps(a, b int) int {
    fmt.Printf("GCD(%d, %d)の計算:\\n", a, b)
    step := 1

    for b != 0 {
        remainder := a % b
        fmt.Printf("Step %d: %d = %d × %d + %d\\n", step, a, b, a/b, remainder)
        a, b = b, remainder
        step++
    }

    fmt.Printf("\\n最大公約数: %d\\n", a)
    return a
}
```

### 拡張ユークリッドの互除法

一次不定方程式 `ax + by = gcd(a,b)` の整数解x, yも求める：

```go
// ExtendedGCD は拡張ユークリッドの互除法
// ax + by = gcd(a,b) となるx, yも返す
func ExtendedGCD(a, b int) (gcd, x, y int) {
    if b == 0 {
        return a, 1, 0
    }

    gcd, x1, y1 := ExtendedGCD(b, a%b)
    x = y1
    y = x1 - (a/b)*y1

    return gcd, x, y
}
```

### 最小公倍数（LCM）の計算

GCDを利用してLCMを求める：

```go
// LCM は最小公倍数を求める
func LCM(a, b int) int {
    return (a * b) / GCD(a, b)
}
```

### 複数の数の最大公約数

```go
// GCDMultiple は複数の数の最大公約数を求める
func GCDMultiple(nums ...int) int {
    if len(nums) == 0 {
        return 0
    }

    result := nums[0]
    for i := 1; i < len(nums); i++ {
        result = GCD(result, nums[i])
        if result == 1 {
            break // 1になったらそれ以上小さくならない
        }
    }

    return result
}
```

### 互いに素の判定

```go
// IsCoprime は2つの数が互いに素かどうかを判定
func IsCoprime(a, b int) bool {
    return GCD(a, b) == 1
}
```

### 分数の約分

```go
// SimplifyFraction は分数を最簡分数に約分
func SimplifyFraction(numerator, denominator int) (int, int) {
    if denominator == 0 {
        return numerator, denominator
    }

    gcd := GCD(abs(numerator), abs(denominator))
    return numerator / gcd, denominator / gcd
}

// abs は絶対値を返す
func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
```

## 計算量

### 時間計算量

- **最良ケース**: O(1) - bが0の場合
- **平均・最悪ケース**: O(log min(a,b))
- フィボナッチ数列の隣接項が最悪ケース

### 空間計算量

- **再帰版**: O(log min(a,b)) - 再帰スタック
- **反復版**: O(1) - 定数空間

### 計算回数の上界

ラメの定理により、除算回数は高々 `5 × (10進数での桁数)` 回

## 使いどころ

### 向いてる場面

- 分数の約分
- 最小公倍数の計算
- 暗号理論（RSA暗号など）
- 座標系での格子点問題
- 音楽理論での周波数比の計算

### 実例：分数計算機

```go
type Fraction struct {
    Numerator   int
    Denominator int
}

// NewFraction は新しい分数を作成（自動約分）
func NewFraction(num, den int) Fraction {
    if den == 0 {
        panic("分母が0です")
    }

    // 負の数の処理
    if den < 0 {
        num = -num
        den = -den
    }

    // 約分
    gcd := GCD(abs(num), abs(den))
    return Fraction{
        Numerator:   num / gcd,
        Denominator: den / gcd,
    }
}

// Add は分数の加算
func (f Fraction) Add(other Fraction) Fraction {
    num := f.Numerator*other.Denominator + other.Numerator*f.Denominator
    den := f.Denominator * other.Denominator
    return NewFraction(num, den)
}

// String は分数を文字列で表現
func (f Fraction) String() string {
    if f.Denominator == 1 {
        return fmt.Sprintf("%d", f.Numerator)
    }
    return fmt.Sprintf("%d/%d", f.Numerator, f.Denominator)
}
```

### 実例：最適な画面アスペクト比

```go
// GetAspectRatio は画面解像度からアスペクト比を求める
func GetAspectRatio(width, height int) string {
    gcd := GCD(width, height)
    w := width / gcd
    h := height / gcd
    return fmt.Sprintf("%d:%d", w, h)
}

// 使用例
// GetAspectRatio(1920, 1080) → "16:9"
// GetAspectRatio(1280, 1024) → "5:4"
// GetAspectRatio(2560, 1440) → "16:9"
```

## 数学的性質

### 重要な性質

1. **可換性**: GCD(a, b) = GCD(b, a)
2. **結合性**: GCD(a, GCD(b, c)) = GCD(GCD(a, b), c)
3. **分配性**: GCD(ma, mb) = m × GCD(a, b)
4. **ベズーの等式**: 任意のa, bに対して、ax + by = GCD(a, b)となる整数x, yが存在

### 証明の概要

なぜ `GCD(a, b) = GCD(b, a mod b)` が成り立つか：

```
a = bq + r とすると（qは商、rは余り）

aとbの公約数をgとすると：
- a = g × m
- b = g × n
（m, nは整数）

r = a - bq = gm - gnq = g(m - nq)
よってgはrの約数でもある

逆も同様に証明でき、
「aとbの公約数」=「bとrの公約数」
となるため、最大公約数も等しい
```

## 他のアルゴリズムとの比較

| 手法 | 計算量 | メモリ | 特徴 |
| --- | --- | --- | --- |
| ユークリッド互除法 | O(log n) | O(1) | 高速・シンプル |
| 素因数分解 | O(√n) | O(log n) | 小さい数向け |
| 二進GCD法 | O(log²n) | O(1) | ビット演算使用 |
| 総当たり | O(n) | O(1) | 非効率 |

## 最適化と発展

### 二進GCD法（Stein's Algorithm）

除算を使わずビット演算で高速化：

```go
func BinaryGCD(a, b int) int {
    if a == 0 {
        return b
    }
    if b == 0 {
        return a
    }

    // 2の共通因子を取り出す
    shift := 0
    for ((a | b) & 1) == 0 {
        a >>= 1
        b >>= 1
        shift++
    }

    // aを奇数にする
    for (a & 1) == 0 {
        a >>= 1
    }

    for b != 0 {
        // bを奇数にする
        for (b & 1) == 0 {
            b >>= 1
        }

        // a > bになるよう調整
        if a > b {
            a, b = b, a
        }

        b = b - a
    }

    return a << shift
}
```

### 連分数展開への応用

ユークリッドの互除法の商の列は連分数展開と対応：

```go
// ContinuedFraction は連分数展開を求める
func ContinuedFraction(num, den int) []int {
    result := []int{}

    for den != 0 {
        q := num / den
        result = append(result, q)
        num, den = den, num-q*den
    }

    return result
}

// 例: 355/113 (円周率の近似値)
// → [3, 7, 15, 1, 292, ...]
```

## まとめ

### メリット

- 極めて効率的（対数時間）
- 実装が簡単
- 数学的に美しい
- 応用範囲が広い
- 2300年の歴史が証明する信頼性

### デメリット

- 大きな数では除算がボトルネック
- 並列化が困難
- 浮動小数点数には適用不可

### 使うべき時

- 分数を扱う
- 最小公倍数が必要
- 暗号処理
- 数論的問題を解く
- 座標・図形処理

ユークリッドの互除法は、プログラミングで最も基本的かつ重要なアルゴリズムの一つ。シンプルながら強力で、数学とコンピュータサイエンスの橋渡しとなる美しいアルゴリズム。