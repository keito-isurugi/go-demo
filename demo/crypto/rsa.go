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

func execRSA() {
	n := generateN()
	fmt.Println("N:", n)

	l := generateL()
	fmt.Println("L:", l)
}
