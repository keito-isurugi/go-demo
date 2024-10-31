package performance

import (
	"fmt"
	"time"
)

func PerformanceDemo() {
	// スライスとサイズの定義
	const size = 100000
	var slice []int

	// appendのパフォーマンス測定
	start := time.Now()
	for i := 0; i < size; i++ {
		slice = append(slice, i)
	}
	elapsed := time.Since(start)
	fmt.Printf("append: %v\n", elapsed)

	// スライスをリセット
	slice = make([]int, size)

	// hoge[i] = fooのパフォーマンス測定
	start = time.Now()
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	elapsed = time.Since(start)
	fmt.Printf("hoge[i] = foo: %v\n", elapsed)
}
