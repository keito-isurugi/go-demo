package main

import (
	"fmt"
	"time"
)

func PerformanceDemo(size int) (appendElapsed, indexElapsed time.Duration) {
	var slice []int

	// appendのパフォーマンス測定
	start := time.Now()
	for i := 0; i < size; i++ {
		slice = append(slice, i)
	}
	appendElapsed = time.Since(start)

	// スライスをリセット
	slice = make([]int, size)

	// インデックスを指定して直接要素を追加のパフォーマンス測定
	start = time.Now()
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	indexElapsed = time.Since(start)

	return appendElapsed, indexElapsed
}

func PerformanceHnadler() {
	const size = 10000000
	count := 100
	var appendElapsed, indexElapsed time.Duration
	for i := 0; i < count; i++ {
		ae, ie := PerformanceDemo(size)
		appendElapsed += ae
		indexElapsed += ie
	}
	fmt.Printf("スライスに %v 個の要素を追加。 %v 回繰り返す\n", size, count)
	fmt.Printf("append: %v\n", appendElapsed)
	fmt.Printf("インデックスを指定して直接要素を追加: %v\n", indexElapsed)

	// appendElapsed が 0 でないことを確認して比率を計算
	if appendElapsed <= 0 {
		fmt.Println("appendElapsedが0のため、比率を計算できません")
	}
	speedRatio := float64(appendElapsed) / float64(indexElapsed)
	fmt.Printf("インデックス指定の方が append より %.2f 倍速い\n", speedRatio)
}