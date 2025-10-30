package main

import "fmt"

// ============================================
// 8. close()を使ったchanelの終了通知
// ============================================

func example8_CloseChannel() {
	fmt.Println("\n=== Example 8: close()でchannelを閉じる ===")

	ch := make(chan int, 5)

	// 送信側
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // channelを閉じる
	}()

	// 受信側（rangeを使うとchannelが閉じられるまでループ）
	for num := range ch {
		fmt.Printf("受信: %d\n", num)
	}

	// channelが閉じられているか確認
	val, ok := <-ch
	if !ok {
		fmt.Println("channelは閉じられています")
	} else {
		fmt.Printf("受信: %d\n", val)
	}

	fmt.Println("Example 8 完了")
}
