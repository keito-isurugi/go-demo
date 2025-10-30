package main

import (
	"fmt"
	"time"
)

// ============================================
// 6. タイムアウト処理
// ============================================

func example6_Timeout() {
	fmt.Println("\n=== Example 6: タイムアウト処理 ===")

	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "処理完了"
	}()

	select {
	case result := <-ch:
		fmt.Println("成功:", result)
	case <-time.After(1 * time.Second):
		fmt.Println("タイムアウト: 1秒以内に完了しませんでした")
	}

	fmt.Println("Example 6 完了")
}
