package main

import (
	"fmt"
)

const (
	p, q = 3, 5
)

// n = p * q 	
func generateN() int {
	n := p * q
	return n // 15 = 3 * 5
}

// φ(n) = (p-1)*(q-1) 
func generateL() int {
	l := (p - 1) * (q - 1)
	return l // 8 = (3-1)*(5-1)
}

// 1 < E < L
// gcd(E, L)
func generateE(l int) int {
	e := 2
	for e < l {
		if gcd(e, l) == 1 {
			break
		}
		e++
	}
	return e // 3
}

// gcd(E, L) = 1を満たしているかをチェック
func gcd(e, l int) int {
	for l != 0 {
		e, l = l, e % l
	}
	return e
}

// 1 < D < L
// E * D mod L = 1
func generateD(e, l int) int {
	d := 1
	for {
		if (e * d) % l == 1 {
			break
		}
		d++
	}
	return d // 3
}

// m^e mod n
func modExp(m, e, n int) int {
	result := 1
	for e > 0 {
		if e % 2 == 1 {
			result = (result * m) % n
		}
		m = (m * m) % n
		e /= 2
	}
	return result
}

func rsaSignature() {
	n := generateN()
	fmt.Println("N:", n)

	l := generateL()
	fmt.Println("L:", l)

	e := generateE(l)
	fmt.Println("E:", e)

	d := generateD(e, l)
	fmt.Println("D:", d)

	// N未満の数値を指定
	message := 8

	// 署名生成: s = m^d mod n
	signature := modExp(message, d, n) // 2 = 8^3 mod 15
	fmt.Println("メッセージ:", message)
	fmt.Println("署名:", signature)

	// 検証: v = s^e mod n
	verified := modExp(signature, e, n) // 8 = 2^3 mod 15
	fmt.Println("検証結果:", verified)

	if verified == message {
		fmt.Println("検証成功: メッセージが正しいです。")
	} else {
		fmt.Println("検証失敗: メッセージが改ざんされています。")
	}
}
