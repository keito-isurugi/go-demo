package rsademo

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

const capacity = 2048

// 素数生成関数（簡略化された試行錯誤法）
func generatePrime(bits int) *big.Int {
	for {
		// 指定されたビット長のランダムな数を生成
		primeCandidate := generateRandomNumber(bits)

		// 素数かどうかをチェック（この例では簡略化された方法）
		if isPrime(primeCandidate, 20) {
			return primeCandidate
		}
	}
}

// ランダムなビット長の数を生成
func generateRandomNumber(bits int) *big.Int {
	max := new(big.Int).Lsh(big.NewInt(1), uint(bits)) // 2^bits
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("ランダム数の生成に失敗しました: %v", err)
	}
	return num
}

// Miller-Rabin法による素数判定
func isPrime(n *big.Int, k int) bool {
	// nが2以下の場合
	if n.Cmp(big.NewInt(2)) <= 0 {
		return false
	}
	// 偶数を除外
	if new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// n-1 = 2^s * d という形に分解
	d := new(big.Int).Sub(n, big.NewInt(1))
	s := 0
	for new(big.Int).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		d.Div(d, big.NewInt(2))
		s++
	}

	// k回の繰り返しで確率的に素数か合成数かを判断
	for i := 0; i < k; i++ {
		// ランダムな数aを選ぶ
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
		if err != nil {
			log.Fatalf("ランダム数の生成に失敗しました: %v", err)
		}
		a.Add(a, big.NewInt(2))

		// a^d mod n を計算
		x := new(big.Int).Exp(a, d, n)

		// xが1またはn-1であれば、次の試行へ
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
			continue
		}

		// x^2 mod n の計算
		for r := 0; r < s-1; r++ {
			x = new(big.Int).Exp(x, big.NewInt(2), n)
			// xがn-1になった場合は合成数ではない
			if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
				break
			}
		}

		// 最後までn-1に到達しない場合は合成数
		if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) != 0 {
			return false
		}
	}
	return true
}

// RSA鍵ペアを生成
func generateRSAKeyPair(bits int) (*big.Int, *big.Int, *big.Int) {
	// 2つの素数pとqを生成
	p := generatePrime(bits / 2)
	q := generatePrime(bits / 2)

	// 公開鍵のmodulus n = p * q
	n := new(big.Int).Mul(p, q)

	// オイラーのトーシェント関数φ(n) = (p-1)*(q-1)
	phi := new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))

	// 公開鍵e（一般的に65537を使用）
	e := big.NewInt(65537)

	// 秘密鍵dを計算：d = e^(-1) mod φ(n)
	d := new(big.Int).ModInverse(e, phi)

	return n, e, d
}

func ExecRSA() {
	n, e, b := generateRSAKeyPair(capacity)
	fmt.Println("n: ", n)
	fmt.Println("e: ", e)
	fmt.Println("b: ", b)
}