package main

import (
	"fmt"
	"time"
)

// ============================================
// 7. done channelパターン（goroutineの完了通知）
// ============================================

func example7_DoneChannel() {
	fmt.Println("\n=== Example 7: done channelパターン ===")

	done := make(chan bool)

	go func() {
		fmt.Println("goroutine: 作業開始...")
		time.Sleep(2 * time.Second)
		fmt.Println("goroutine: 作業完了")
		done <- true // 完了を通知
	}()

	fmt.Println("main: goroutineの完了を待っています...")
	<-done
	fmt.Println("main: goroutineが完了しました")
	fmt.Println("Example 7 完了")
}
