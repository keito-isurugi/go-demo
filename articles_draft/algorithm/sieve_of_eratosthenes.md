## 概要

- Goでエラトステネスの篩（Sieve of Eratosthenes）を実装
- エラトステネスの篩は、指定された範囲内のすべての素数を効率的に列挙するアルゴリズム。古代ギリシャの数学者エラトステネスによって考案された。

### 特徴

- **高速な素数列挙**: 時間計算量O(n log log n)で効率的
- **シンプルな実装**: アルゴリズムが直感的で理解しやすい
- **大量の素数探索に最適**: 単一の素数判定よりも、範囲内のすべての素数を求める場合に効率的
- **メモリ使用**: 空間計算量O(n)でブール配列を使用
- **応用範囲が広い**: 素数の個数カウント、双子素数探索、範囲指定の素数列挙など

## エラトステネスの篩とは

### 基本的な考え方

エラトステネスの篩は、「合成数を除外していく」という考え方で素数を見つける手法。2から始めて、素数の倍数をすべて合成数としてマークしていく。

**基本アルゴリズム**:
```
1. 2〜nまでのすべての数を素数候補とする
2. 2から順に処理:
   - 現在の数が素数候補なら、その倍数をすべて合成数としてマーク
   - √nまで繰り返す
3. 残った素数候補が素数
```

### 視覚的イメージ

```
n=30の例:

初期状態（2〜30をすべて素数候補とする）:
2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30

Step 1: 2は素数 → 2の倍数（4,6,8...）を除外
2  3  ✗  5  ✗  7  ✗  9  ✗  11 ✗  13 ✗  15 ✗  17 ✗  19 ✗  21 ✗  23 ✗  25 ✗  27 ✗  29 ✗

Step 2: 3は素数 → 3の倍数（9,15,21...）を除外（6,12は既に除外済み）
2  3  ✗  5  ✗  7  ✗  ✗  ✗  11 ✗  13 ✗  ✗  ✗  17 ✗  19 ✗  ✗  ✗  23 ✗  25 ✗  ✗  ✗  29 ✗

Step 3: 5は素数 → 5の倍数（25）を除外（10,15,20は既に除外済み）
2  3  ✗  5  ✗  7  ✗  ✗  ✗  11 ✗  13 ✗  ✗  ✗  17 ✗  19 ✗  ✗  ✗  23 ✗  ✗  ✗  ✗  ✗  29 ✗

結果: 2, 3, 5, 7, 11, 13, 17, 19, 23, 29（10個の素数）
```

### なぜ効率的か？

**試し割り法（素朴な方法）**:
- 各数について、2〜√nで割り切れるかチェック
- 時間計算量: O(n√n)

```go
// 試し割り法: 各数について√nまで試し割り
for i := 2; i <= n; i++ {
    isPrime := true
    for j := 2; j*j <= i; j++ {
        if i % j == 0 {
            isPrime = false
            break
        }
    }
    if isPrime {
        // iは素数
    }
}
```

**エラトステネスの篩**:
- 素数の倍数をまとめて除外
- 時間計算量: O(n log log n)

```go
// エラトステネスの篩: 倍数をまとめて除外
isPrime := make([]bool, n+1)
for i := 2; i <= n; i++ {
    isPrime[i] = true
}

for i := 2; i*i <= n; i++ {
    if isPrime[i] {
        // iの倍数をすべて除外
        for j := i * i; j <= n; j += i {
            isPrime[j] = false
        }
    }
}
```

**性能差**:
```
n=10,000の素数を求める場合:
試し割り法:           約10,000 × √10,000 = 1,000,000回の計算
エラトステネスの篩:   約10,000 × log(log 10,000) ≈ 25,000回の計算
→ 約40倍高速！
```

## サンプルコード

### 1. 基本的なエラトステネスの篩

n以下のすべての素数を求める。

```go
package main

import "fmt"

// SieveOfEratosthenes はエラトステネスの篩を使ってn以下の素数をすべて求める
func SieveOfEratosthenes(n int) []int {
	if n < 2 {
		return []int{}
	}

	// すべての数を素数候補とする（true = 素数）
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	// 2からsqrt(n)まで処理
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			// iの倍数をすべて合成数としてマーク
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	// 素数を収集
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}
```

**使用例**:
```go
primes := SieveOfEratosthenes(30)
// 結果: [2, 3, 5, 7, 11, 13, 17, 19, 23, 29]
```

### 2. 素数の個数を数える

n以下の素数が何個あるかを返す。

```go
// CountPrimes はn以下の素数の個数を返す
func CountPrimes(n int) int {
	if n < 2 {
		return 0
	}

	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	count := 0
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			count++
		}
	}

	return count
}
```

**使用例**:
```go
count := CountPrimes(100)
// 結果: 25（100以下の素数は25個）

CountPrimes(10)    // 4個: 2, 3, 5, 7
CountPrimes(1000)  // 168個
CountPrimes(10000) // 1229個
```

### 3. 範囲指定での素数探索

範囲[low, high]内の素数を求める（セグメント篩）。

```go
// SieveOfEratosthenesRange は範囲[low, high]の素数を求める
func SieveOfEratosthenesRange(low, high int) []int {
	if high < 2 {
		return []int{}
	}

	// まず√highまでの素数を求める
	limit := int(float64(high) * 0.5)
	if limit < 2 {
		limit = 2
	}

	smallPrimes := SieveOfEratosthenes(limit)

	// 範囲[low, high]の素数候補
	isPrime := make([]bool, high-low+1)
	for i := range isPrime {
		isPrime[i] = true
	}

	// 小さな素数の倍数を除外
	for _, p := range smallPrimes {
		// pの倍数の開始位置を計算
		start := ((low + p - 1) / p) * p
		if start < p*p {
			start = p * p
		}

		for j := start; j <= high; j += p {
			isPrime[j-low] = false
		}
	}

	// 素数を収集
	primes := []int{}
	start := low
	if start < 2 {
		start = 2
	}

	for i := start; i <= high; i++ {
		if isPrime[i-low] {
			primes = append(primes, i)
		}
	}

	return primes
}
```

**使用例**:
```go
primes := SieveOfEratosthenesRange(50, 100)
// 結果: [53, 59, 61, 67, 71, 73, 79, 83, 89, 97]
```

### 4. 試し割り法（比較用）

単一の数が素数かどうかを判定する（エラトステネスの篩との比較用）。

```go
// IsPrimeSimple は単純な素数判定（試し割り法）
func IsPrimeSimple(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}
```

**使用例**:
```go
IsPrimeSimple(17)  // true
IsPrimeSimple(18)  // false
IsPrimeSimple(97)  // true
IsPrimeSimple(100) // false
```

### 5. ステップごとの可視化

計算過程を表示するバージョン。

```go
// SieveOfEratosthenesWithSteps はステップごとに計算過程を表示
func SieveOfEratosthenesWithSteps(n int) []int {
	fmt.Printf("\nn以下の素数を求める (n=%d)\n", n)
	fmt.Println("=================================")

	if n < 2 {
		fmt.Println("結果: 素数なし")
		return []int{}
	}

	// すべての数を素数候補とする
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	fmt.Printf("\n初期状態: 2〜%dをすべて素数候補とする\n", n)

	fmt.Println("\nエラトステネスの篩を開始:")
	fmt.Println("=================================")

	// 2からsqrt(n)まで処理
	step := 1
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			fmt.Printf("\nStep %d: %dは素数\n", step, i)
			fmt.Printf("        %dの倍数を合成数としてマーク: ", i)

			// iの倍数をすべて合成数としてマーク
			markedNums := []int{}
			for j := i * i; j <= n; j += i {
				if isPrime[j] {
					isPrime[j] = false
					markedNums = append(markedNums, j)
				}
			}

			if len(markedNums) > 0 {
				fmt.Printf("%v\n", markedNums)
			} else {
				fmt.Println("なし")
			}

			step++
		}
	}

	// 素数を収集
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	fmt.Printf("\n結果: %d個の素数が見つかりました\n", len(primes))
	fmt.Printf("素数: %v\n", primes)

	return primes
}
```

**出力例**:
```
n以下の素数を求める (n=30)
=================================

初期状態: 2〜30をすべて素数候補とする

エラトステネスの篩を開始:
=================================

Step 1: 2は素数
        2の倍数を合成数としてマーク: [4 6 8 10 12 14 16 18 20 22 24 26 28 30]

Step 2: 3は素数
        3の倍数を合成数としてマーク: [9 15 21 27]

Step 3: 5は素数
        5の倍数を合成数としてマーク: [25]

結果: 10個の素数が見つかりました
素数: [2 3 5 7 11 13 17 19 23 29]
```

## 動作原理

### なぜi*iから始めるのか？

エラトステネスの篩では、内側のループを`j := i * i`から始める理由：

```
iの倍数を除外する場合:
i×2, i×3, i×4, ..., i×(i-1), i×i, i×(i+1), ...

例: i=5の場合
5×2 = 10 → 既に2の処理で除外済み
5×3 = 15 → 既に3の処理で除外済み
5×4 = 20 → 既に2の処理で除外済み
5×5 = 25 → ここから新しく除外（★）

つまり、i×k (k < i)の倍数は、kの処理で既に除外されている。
したがって、i×iから始めれば十分！
```

**効率化の効果**:
```go
// 非効率な実装
for j := 2*i; j <= n; j += i {
    isPrime[j] = false
}

// 効率的な実装
for j := i*i; j <= n; j += i {
    isPrime[j] = false
}

n=100, i=7の場合:
非効率: 14, 21, 28, 35, 42, 49, ... (14回)
効率的: 49, 56, 63, 70, ...       (8回)
→ 約半分の処理で済む！
```

### なぜ√nまでチェックすればよいのか？

```
合成数nは必ず√n以下の約数を持つ。

証明:
nが合成数なら、n = a × b (a, b > 1)と表せる。
もしa > √n かつ b > √nなら、
a × b > √n × √n = n（矛盾）

したがって、a ≤ √n または b ≤ √nが成立。

例: n=100の場合
√100 = 10なので、10までチェックすればよい。
11以上の合成数は、すでに10以下の素数で除外済み。
```

### アルゴリズムの流れ

n=30の詳細な動作：

```
初期状態: 2〜30をすべて素数候補とする
[2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30]

Step 1: i=2（2は素数）
        2の倍数（4, 6, 8, ...）を除外
        j = 2×2=4から開始: 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30を除外
[2, 3, ✗, 5, ✗, 7, ✗, 9, ✗, 11, ✗, 13, ✗, 15, ✗, 17, ✗, 19, ✗, 21, ✗, 23, ✗, 25, ✗, 27, ✗, 29, ✗]

Step 2: i=3（3は素数）
        3の倍数を除外
        j = 3×3=9から開始: 9, 15, 21, 27を除外（6, 12, 18, 24, 30は既に除外済み）
[2, 3, ✗, 5, ✗, 7, ✗, ✗, ✗, 11, ✗, 13, ✗, ✗, ✗, 17, ✗, 19, ✗, ✗, ✗, 23, ✗, 25, ✗, ✗, ✗, 29, ✗]

Step 3: i=5（5は素数）
        5の倍数を除外
        j = 5×5=25から開始: 25を除外（10, 15, 20は既に除外済み）
[2, 3, ✗, 5, ✗, 7, ✗, ✗, ✗, 11, ✗, 13, ✗, ✗, ✗, 17, ✗, 19, ✗, ✗, ✗, 23, ✗, ✗, ✗, ✗, ✗, 29, ✗]

終了: i=6以降はi×i > 30なので処理不要

結果: 2, 3, 5, 7, 11, 13, 17, 19, 23, 29
```

## 計算量

### 時間計算量

| 手法 | 時間計算量 | 説明 |
| --- | --- | --- |
| 試し割り法（1つの数） | O(√n) | √nまで割り切れるかチェック |
| 試し割り法（n個の数） | O(n√n) | 各数についてO(√n) |
| エラトステネスの篩 | O(n log log n) | 調和級数の和 |

### なぜO(n log log n)なのか？

エラトステネスの篩の時間計算量の導出：

```
各素数pについて、n/p回の除外処理を行う。
素数は2, 3, 5, 7, 11, ...なので、

総処理回数 = n/2 + n/3 + n/5 + n/7 + n/11 + ...
          = n × (1/2 + 1/3 + 1/5 + 1/7 + 1/11 + ...)
          = n × (素数の逆数の和)

素数の逆数の和は、調和級数の和に近似でき、
Σ(1/p) ≈ log(log n)

したがって、時間計算量は O(n log log n)
```

**実測値の比較**:
```
n=10,000の場合:
試し割り法:           約1,000,000回の計算
エラトステネスの篩:   約25,000回の計算

n=100,000の場合:
試し割り法:           約31,600,000回の計算
エラトステネスの篩:   約300,000回の計算

→ nが大きいほど、エラトステネスの篩の優位性が顕著！
```

### 空間計算量

| 手法 | 空間計算量 | 説明 |
| --- | --- | --- |
| 試し割り法 | O(1) | 追加メモリほぼ不要 |
| エラトステネスの篩 | O(n) | n個のブール配列 |
| セグメント篩 | O(√n + (high-low)) | 範囲指定版 |

## 使いどころ

### 向いている場面

1. **大量の素数が必要な場合**
   - n以下のすべての素数を求める
   - 素数表の作成

2. **素数の個数を数える**
   - n以下の素数が何個あるか
   - 素数密度の計算

3. **範囲内の素数探索**
   - [low, high]の素数をすべて求める
   - セグメント篩を使用

4. **素数関連の前処理**
   - 素因数分解の前準備
   - 約数列挙の前処理

### 向いていない場面

1. **単一の大きな数の素数判定**
   - 試し割り法やミラーラビン法が効率的
   - エラトステネスの篩はオーバーキル

2. **メモリ制限が厳しい場合**
   - O(n)のメモリが必要
   - 試し割り法ならO(1)

3. **非常に大きなnの場合**
   - n=10^9以上はメモリ不足の可能性
   - セグメント篩やビット配列で最適化が必要

### 実例1: 双子素数の探索

差が2の素数ペアを見つける。

```go
// FindTwinPrimes は双子素数（差が2の素数ペア）を見つける
func FindTwinPrimes(n int) [][2]int {
	primes := SieveOfEratosthenes(n)
	twinPrimes := [][2]int{}

	for i := 0; i < len(primes)-1; i++ {
		if primes[i+1]-primes[i] == 2 {
			twinPrimes = append(twinPrimes, [2]int{primes[i], primes[i+1]})
		}
	}

	return twinPrimes
}

// 使用例
twins := FindTwinPrimes(100)
// 結果: [(3,5), (5,7), (11,13), (17,19), (29,31), (41,43), (59,61), (71,73)]
```

### 実例2: 素数の分布

範囲ごとの素数の個数を調べる。

```go
// PrimeDistribution は範囲ごとの素数の個数を返す
func PrimeDistribution(n, rangeSize int) []int {
	primes := SieveOfEratosthenes(n)
	distribution := make([]int, (n/rangeSize)+1)

	for _, p := range primes {
		index := p / rangeSize
		distribution[index]++
	}

	return distribution
}

// 使用例
dist := PrimeDistribution(100, 10)
// 結果: [0, 4, 4, 2, 2, 3, 2, 4, 2, 2, 0]
// [1-10]: 4個, [11-20]: 4個, [21-30]: 2個, ...
```

### 実例3: 素因数分解

エラトステネスの篩を使った素因数分解。

```go
// PrimeFactorization は素因数分解を行う
func PrimeFactorization(n int) map[int]int {
	primes := SieveOfEratosthenes(int(float64(n)*0.5) + 1)
	factors := make(map[int]int)

	for _, p := range primes {
		for n%p == 0 {
			factors[p]++
			n /= p
		}
		if n == 1 {
			break
		}
	}

	if n > 1 {
		factors[n]++
	}

	return factors
}

// 使用例
factors := PrimeFactorization(60)
// 結果: {2: 2, 3: 1, 5: 1}（60 = 2² × 3 × 5）

factors2 := PrimeFactorization(100)
// 結果: {2: 2, 5: 2}（100 = 2² × 5²）
```

## 最適化テクニック

### 1. 偶数のスキップ

2以外の偶数は素数ではないので、奇数のみ処理する。

```go
// SieveOfEratosthenesOddOnly は奇数のみを処理（最適化版）
func SieveOfEratosthenesOddOnly(n int) []int {
	if n < 2 {
		return []int{}
	}
	if n == 2 {
		return []int{2}
	}

	// 奇数のみを扱う（index i → 実際の数 2*i+3）
	size := (n - 1) / 2
	isPrime := make([]bool, size)
	for i := range isPrime {
		isPrime[i] = true
	}

	// 3から処理（2ずつ増加）
	for i := 0; i*i < size; i++ {
		if isPrime[i] {
			p := 2*i + 3
			// pの奇数倍を除外
			for j := (p*p - 3) / 2; j < size; j += p {
				isPrime[j] = false
			}
		}
	}

	// 結果を収集（2を先頭に追加）
	primes := []int{2}
	for i := 0; i < size; i++ {
		if isPrime[i] {
			primes = append(primes, 2*i+3)
		}
	}

	return primes
}
```

**効果**: メモリ使用量が約半分、処理速度も向上。

### 2. ビット配列の使用

ブール配列の代わりにビット配列を使用してメモリを削減。

```go
// ビット配列版（実装例は省略）
// 1バイトで8個の数を表現
// メモリ使用量: O(n) → O(n/8)
```

### 3. セグメント篩

巨大なnに対して、範囲を分割して処理。

```go
// 既に実装済みのSieveOfEratosthenesRangeがセグメント篩
// メモリ使用量: O(n) → O(√n + (high-low))
```

## デバッグとテスト

### テストケース

```go
func TestSieveOfEratosthenes(t *testing.T) {
	// 基本ケース
	primes := SieveOfEratosthenes(30)
	expected := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	if !reflect.DeepEqual(primes, expected) {
		t.Errorf("Expected %v, got %v", expected, primes)
	}

	// 小さいn
	if len(SieveOfEratosthenes(1)) != 0 {
		t.Error("n=1 should return empty")
	}

	// n=2
	primes2 := SieveOfEratosthenes(2)
	if len(primes2) != 1 || primes2[0] != 2 {
		t.Error("n=2 should return [2]")
	}

	// 素数の個数チェック
	count := CountPrimes(100)
	if count != 25 {
		t.Errorf("Expected 25 primes up to 100, got %d", count)
	}
}
```

### デバッグのポイント

1. **配列の範囲チェック**
   ```go
   // n+1のサイズで配列を作成（0〜nまで）
   isPrime := make([]bool, n+1)
   ```

2. **ループの範囲**
   ```go
   // i*i <= n（等号を忘れない）
   for i := 2; i*i <= n; i++ {
   ```

3. **倍数の開始位置**
   ```go
   // i*iから始める（2*iではない）
   for j := i * i; j <= n; j += i {
   ```

## まとめ

### メリット

- 時間計算量O(n log log n)で高速
- 試し割り法O(n√n)より圧倒的に効率的
- 大量の素数を一度に求められる
- アルゴリズムがシンプルで理解しやすい
- 素数関連の問題に幅広く応用可能

### デメリット

- 空間計算量O(n)でメモリを消費
- 単一の素数判定には不向き
- 非常に大きなnではメモリ不足の可能性
- 事前に範囲を決める必要がある

### 使うべき時

- **n以下のすべての素数が必要**: エラトステネスの篩
- **素数の個数を数える**: エラトステネスの篩
- **範囲内の素数を列挙**: セグメント篩
- **素因数分解の前処理**: エラトステネスの篩
- **大量の素数判定**: エラトステネスの篩で事前計算

### 選択基準

| 目的 | 最適な手法 | 時間計算量 |
| --- | --- | --- |
| 1つの数の素数判定 | 試し割り法 | O(√n) |
| n個の数の素数判定 | エラトステネスの篩 | O(n log log n) |
| 範囲[low, high]の素数 | セグメント篩 | O((high-low) log log high) |
| 巨大な数の素数判定 | ミラーラビン法 | O(k log³ n) |

エラトステネスの篩は、紀元前3世紀に考案された古典的アルゴリズムながら、現代でも最も効率的な素数列挙法の一つ。シンプルな原理ながら、調和級数の性質を利用した巧妙な時間計算量O(n log log n)を実現。素数の一覧が必要な場面では、試し割り法より圧倒的に効率的で、数論やセキュリティ分野で広く使用される重要なアルゴリズム。
