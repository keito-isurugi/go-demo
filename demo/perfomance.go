package demo

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
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

func PerformanceProfDemo() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
		return
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile: ", err)
		return
	}
	defer pprof.StopCPUProfile()

	// スライスとサイズの定義
	const size = 10000000  // サイズを増やす
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

	// 要素の代入のパフォーマンス測定
	start = time.Now()
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	elapsed = time.Since(start)
	fmt.Printf("slice[i] = i: %v\n", elapsed)
}

func PerformanceTraceDemo() {
	// トレースファイルを作成
	f, err := os.Create("trace.out")
	if err != nil {
		fmt.Println("could not create trace file: ", err)
		return
	}
	defer f.Close()

	// トレースの開始
	if err := trace.Start(f); err != nil {
		fmt.Println("could not start trace: ", err)
		return
	}
	defer trace.Stop()

	// スライスとサイズの定義
	const size = 10000000
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

	// 要素の代入のパフォーマンス測定
	start = time.Now()
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	elapsed = time.Since(start)
	fmt.Printf("slice[i] = i: %v\n", elapsed)
}
