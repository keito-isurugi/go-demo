package main

import (
	"crypto/rand"
	mathRand "math/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func partA() {
	partA13()
}

func partA13() {
	a := make([]byte, 20)
	rand.Read(a)
	fmt.Println(hex.EncodeToString(a))

	// 乱数の種を設定
	mathRand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		// 浮動小数点数（float64）の乱数を生成
		fmt.Println(mathRand.Float64())
	}
}
