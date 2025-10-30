package main

import (
	"fmt"
	"time"
)

// ============================================
// 1. 基本的なgoroutineの使い方
// ============================================

func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Hello, %s! (count: %d)\n", name, i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func example1_BasicGoroutine() {
	fmt.Println("\n=== Example 1: 基本的なgoroutine ===")

	// 通常の関数呼び出し（同期）
	fmt.Println("同期実行:")
	sayHello("Alice")

	// goroutineを使った並行実行
	fmt.Println("\n並行実行:")
	go sayHello("Bob")   // goroutineで実行
	go sayHello("Carol") // goroutineで実行

	// main関数が終了するとgoroutineも終了してしまうので待つ
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Example 1 完了")
}
