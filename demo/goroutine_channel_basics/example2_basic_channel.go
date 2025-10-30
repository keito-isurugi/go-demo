package main

import (
	"fmt"
	"time"
)

// ============================================
// 2. 基本的なchannelの使い方
// ============================================

func example2_BasicChannel() {
	fmt.Println("\n=== Example 2: 基本的なchannel ===")

	// channelの作成
	ch := make(chan string)

	// goroutineでchannelに値を送信
	go func() {
		fmt.Println("goroutine: メッセージを送信します...")
		time.Sleep(1 * time.Second)
		ch <- "こんにちは、channelから!" // channelに送信
		fmt.Println("goroutine: メッセージを送信しました")
	}()

	// channelから値を受信（ブロックして待つ）
	fmt.Println("main: メッセージを待っています...")
	message := <-ch // channelから受信
	fmt.Printf("main: 受信したメッセージ: %s\n", message)
	fmt.Println("Example 2 完了")
}
