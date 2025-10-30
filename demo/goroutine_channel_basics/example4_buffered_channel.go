package main

import (
	"fmt"
	"time"
)

// ============================================
// 4. buffered channelとunbuffered channel
// ============================================

func example4_BufferedChannel() {
	fmt.Println("\n=== Example 4: buffered channelとunbuffered channel ===")

	// unbuffered channel（バッファなし）
	fmt.Println("Unbuffered channel:")
	unbuffered := make(chan int)

	go func() {
		fmt.Println("  goroutine: 値を送信します")
		unbuffered <- 1
		fmt.Println("  goroutine: 値を送信しました（受信されるまでブロックされていた）")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("  main: 値を受信します")
	<-unbuffered

	// buffered channel（バッファあり）
	fmt.Println("\nBuffered channel:")
	buffered := make(chan int, 2) // バッファサイズ2

	buffered <- 1
	buffered <- 2
	fmt.Println("  main: 2つの値を送信しました（ブロックされない）")

	fmt.Printf("  main: 受信した値: %d\n", <-buffered)
	fmt.Printf("  main: 受信した値: %d\n", <-buffered)

	fmt.Println("Example 4 完了")
}
