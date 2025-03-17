package main

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

// - 1  < E < L
// - gcd(E, L) = 1
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

// 最大公約数
func gcd(e, l int) int {
	for l != 0 {
		e, l = l, e % l
	}
	return e
}

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

// m^exp mod n
func modExp(m, exp, n int) int {
	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * m) % n
		}
		m = (m * m) % n
		exp /= 2
	}
	return result
}

func execRSA() {
	n := generateN()
	fmt.Println("N:", n)

	l := generateL()
	fmt.Println("L:", l)

	e := generateE(l)
	fmt.Println("E:", e)

	d := generateD(e, l)
	fmt.Println("D:", d)

	target := 8
	fmt.Println("Target:", target)
	c := modExp(target, e, n)
	fmt.Println("C:", c)
	d = modExp(c, d, n)
	fmt.Println("D:", d)
}
