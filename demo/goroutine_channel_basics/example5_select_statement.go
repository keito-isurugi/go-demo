package main

import (
	"fmt"
	"time"
)

// ============================================
// 5. selectを使った複数channelの処理
// ============================================

func example5_SelectStatement() {
	fmt.Println("\n=== Example 5: select文で複数channelを処理 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "channel 1からのメッセージ"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "channel 2からのメッセージ"
	}()

	// selectで複数channelを待つ
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("受信:", msg1)
		case msg2 := <-ch2:
			fmt.Println("受信:", msg2)
		}
	}

	fmt.Println("Example 5 完了")
}
