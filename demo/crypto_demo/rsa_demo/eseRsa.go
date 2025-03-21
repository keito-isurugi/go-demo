package rsademo

import "fmt"

const (
	p, q = 3, 5
)

func generateN() int {
	n := p * q
	return n
}

func generateL() int {
	l := (p - 1) * (q - 1)
	return l
}

// 1 < E < L
// gcd(E, L) = 1
func generateE(l int) int {
	e := 2
	for e < l {
		if gcd(e, l) == 1 {
			break
		}
		e++
	}
	return e
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
	return d
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

func ExecEseRSA() {
	n := generateN()
	fmt.Println("N:", n)

	l := generateL()
	fmt.Println("L:", l)

	e := generateE(l)
	fmt.Println("E:", e)

	d := generateD(e, l)
	fmt.Println("D:", d)

	// N未満の数値を指定
	target := 8 
	encrypted := modExp(target, e, n)
	decrypted := modExp(encrypted, d, n)
	
	fmt.Println("平文:", target)
	fmt.Println("暗号化:", encrypted)
	fmt.Println("復号化:", decrypted)
}

// 出力
// 平文: 8
// 暗号化: 2
// 復号化: 8
