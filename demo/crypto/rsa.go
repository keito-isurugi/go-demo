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

func gcd(e, l int) int {
	for l != 0 {
		e, l = l, e%l
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

func execRSA() {
	n := generateN()
	fmt.Println("N:", n)

	l := generateL()
	fmt.Println("L:", l)
}
