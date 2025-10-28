## 概要

- Goで階乗計算を再帰的に実装
- 階乗（n!）は1からnまでの自然数の積。再帰的定義がそのままコードに落とし込みやすく、再帰アルゴリズムの基本例として最適。

### 特徴

- **シンプルな実装**: 数学的定義をそのままコード化
- **再帰の基本パターン**: ベースケースと再帰ケースの明確な分離
- **メモ化による最適化**: 計算結果のキャッシュで高速化可能
- **大きな数への対応**: big.Intパッケージで巨大な階乗値も計算可能
- **末尾再帰の非対応**: Goは末尾再帰最適化をサポートしないため、スタックオーバーフローに注意

## 動作原理

基本的な流れ：

1. ベースケース：0! = 1, 1! = 1
2. 再帰ケース：n! = n × (n-1)!
3. 各関数呼び出しが1つ小さい問題に分割
4. ベースケースに到達後、結果を順次掛け合わせて返す

### 具体例

`5!` の計算過程：

```
factorial(5)
├── 5 × factorial(4)
    ├── 4 × factorial(3)
        ├── 3 × factorial(2)
            ├── 2 × factorial(1)
                └── 1 (ベースケース)
            └── 2 × 1 = 2
        └── 3 × 2 = 6
    └── 4 × 6 = 24
└── 5 × 24 = 120
```

## サンプルコード

```go
package main

import (
	"fmt"
	"math/big"
)

// 基本的な再帰による階乗計算
func factorial(n int) uint64 {
	// エラーチェック
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	// ベースケース
	if n == 0 || n == 1 {
		return 1
	}
	// 再帰ケース
	return uint64(n) * factorial(n-1)
}

// メモ化を使った階乗計算
func factorialWithMemo(n int, memo map[int]uint64) uint64 {
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	// ベースケース
	if n == 0 || n == 1 {
		return 1
	}
	// メモにあればそれを返す
	if val, exists := memo[n]; exists {
		return val
	}
	// 再帰計算してメモに保存
	result := uint64(n) * factorialWithMemo(n-1, memo)
	memo[n] = result
	return result
}

// 大きな数に対応した階乗計算（big.Int使用）
func factorialBigInt(n int) *big.Int {
	if n < 0 {
		panic("factorial: negative number not allowed")
	}
	// ベースケース
	if n == 0 || n == 1 {
		return big.NewInt(1)
	}
	// 再帰計算
	result := big.NewInt(int64(n))
	return result.Mul(result, factorialBigInt(n-1))
}

// 計算過程を可視化する階乗計算
func visualizeFactorial(n int, depth int) uint64 {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	// ベースケース
	if n == 0 || n == 1 {
		fmt.Printf("%sfactorial(%d) = 1\n", indent, n)
		return 1
	}
	
	// 再帰呼び出しを表示
	fmt.Printf("%sfactorial(%d) = %d * factorial(%d)\n", indent, n, n, n-1)
	result := uint64(n) * visualizeFactorial(n-1, depth+1)
	fmt.Printf("%sfactorial(%d) = %d\n", indent, n, result)
	return result
}

func main() {
	// 基本的な階乗計算
	fmt.Println("基本的な再帰による階乗計算:")
	for i := 0; i <= 10; i++ {
		result := factorial(i)
		fmt.Printf("%d! = %d\n", i, result)
	}
	
	// メモ化を使った階乗計算
	fmt.Println("\nメモ化を使った階乗計算:")
	memo := make(map[int]uint64)
	for i := 15; i <= 20; i++ {
		result := factorialWithMemo(i, memo)
		fmt.Printf("%d! = %d\n", i, result)
	}
	
	// 大きな数の階乗計算
	fmt.Println("\n大きな数の階乗計算 (big.Int使用):")
	for _, n := range []int{25, 50, 100} {
		result := factorialBigInt(n)
		fmt.Printf("%d! = %s\n", n, result.String())
	}
	
	// 計算過程の可視化
	fmt.Println("\n計算過程の可視化 (5!):")
	visualizeFactorial(5, 0)
}
```

## 時間計算量と空間計算量

### 時間計算量
- **基本実装**: O(n) - n回の再帰呼び出し
- **メモ化版**: O(n) - 各値を1回だけ計算

### 空間計算量
- **基本実装**: O(n) - 再帰スタックの深さ
- **メモ化版**: O(n) - スタック深さ + メモのストレージ

## 実装上の注意点

### 1. オーバーフロー対策
階乗は急速に大きくなるため、uint64でも20!までしか扱えない：
- 20! = 2,432,902,008,176,640,000 (uint64の範囲内)
- 21! = 51,090,942,171,709,440,000 (uint64を超える)

大きな値を扱う場合は`math/big`パッケージを使用：

```go
func factorialBigInt(n int) *big.Int {
    if n == 0 || n == 1 {
        return big.NewInt(1)
    }
    result := big.NewInt(int64(n))
    return result.Mul(result, factorialBigInt(n-1))
}
```

### 2. スタックオーバーフロー対策
Goは末尾再帰最適化をサポートしないため、大きなnでスタックオーバーフローの可能性がある。
必要に応じて反復実装を検討：

```go
func factorialIterative(n int) uint64 {
    result := uint64(1)
    for i := 2; i <= n; i++ {
        result *= uint64(i)
    }
    return result
}
```

### 3. エラーハンドリング
負の数に対する適切なエラー処理：

```go
func factorial(n int) (uint64, error) {
    if n < 0 {
        return 0, fmt.Errorf("factorial of negative number is undefined")
    }
    // ...
}
```

## 応用例

### 組み合わせ計算（nCr）
階乗を使った組み合わせの計算：

```go
func combination(n, r int) *big.Int {
    if r > n || r < 0 {
        return big.NewInt(0)
    }
    // nCr = n! / (r! * (n-r)!)
    nFact := factorialBigInt(n)
    rFact := factorialBigInt(r)
    nrFact := factorialBigInt(n - r)
    
    denominator := new(big.Int).Mul(rFact, nrFact)
    return new(big.Int).Div(nFact, denominator)
}
```

### 順列計算（nPr）
階乗を使った順列の計算：

```go
func permutation(n, r int) *big.Int {
    if r > n || r < 0 {
        return big.NewInt(0)
    }
    // nPr = n! / (n-r)!
    nFact := factorialBigInt(n)
    nrFact := factorialBigInt(n - r)
    return new(big.Int).Div(nFact, nrFact)
}
```

## まとめ

再帰的な階乗計算は：
- **再帰の基本を学ぶ最適な例題**
- **数学的定義との対応が明確**
- **メモ化による最適化の効果を実感しやすい**
- **大きな数の扱いを学ぶ良い機会**

実用的には反復実装の方が効率的だが、再帰の概念理解と応用（組み合わせ・順列計算など）において重要な基礎となる。