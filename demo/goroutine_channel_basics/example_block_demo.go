package main

import (
	"fmt"
	"time"
)

// ============================================
// ブロックの仕組みを理解するデモ
// ============================================

func exampleBlockDemo() {
	fmt.Println("\n=== ブロックの仕組み ===")

	// バッファなしチャネル
	ch := make(chan string)

	fmt.Println("[main] チャネルに送信を試みます...")

	// ここで送信しようとすると、受信者がいないため
	// この行で永遠に止まる（デッドロック）
	// ch <- "hello"  // ← これを実行するとプログラム全体がフリーズ

	// だから goroutine で送信する
	go func() {
		fmt.Println("[goroutine] 1秒待ってから送信します...")
		time.Sleep(1 * time.Second)
		fmt.Println("[goroutine] 送信中... (ここでブロックされるかも)")
		ch <- "メッセージ"  // ← 受信者が現れるまでここで待つ
		fmt.Println("[goroutine] 送信完了！(受信者が受け取った)")
	}()

	// main は2秒待つ
	fmt.Println("[main] 2秒待ちます...")
	time.Sleep(2 * time.Second)

	fmt.Println("[main] 受信します...")
	msg := <-ch  // ← goroutine がここで送信を完了できる
	fmt.Printf("[main] 受信: %s\n", msg)
}

// バッファありチャネルの場合
func exampleBufferedChannel() {
	fmt.Println("\n=== バッファありチャネル ===")

	// バッファサイズ2のチャネル
	ch := make(chan string, 2)

	fmt.Println("[main] 2つ送信します（受信者なしでもOK）")
	ch <- "メッセージ1"  // ← ブロックしない（バッファに空きがある）
	fmt.Println("[main] 1つ目送信完了")
	ch <- "メッセージ2"  // ← ブロックしない（バッファに空きがある）
	fmt.Println("[main] 2つ目送信完了")
	// ch <- "メッセージ3"  // ← これはブロックする（バッファが満杯）

	fmt.Println("[main] 受信します...")
	msg1 := <-ch
	msg2 := <-ch
	fmt.Printf("[main] 受信: %s, %s\n", msg1, msg2)
}

// ゴルーチンリークの例
func exampleGoroutineLeakBad() {
	fmt.Println("\n=== ゴルーチンリーク（悪い例）===")

	ch := make(chan string)

	// このゴルーチンは永遠に終わらない
	go func() {
		fmt.Println("[goroutine] 送信します...")
		ch <- "メッセージ1"  // ← 受信されるのでOK
		fmt.Println("[goroutine] 1つ目送信完了")

		ch <- "メッセージ2"  // ← 誰も受信しないので永遠にここで待つ（ハング）
		fmt.Println("[goroutine] 2つ目送信完了（この行には到達しない！）")
	}()

	// 1つしか受信しない
	msg := <-ch
	fmt.Printf("[main] 受信: %s\n", msg)
	fmt.Println("[main] 終了（でもゴルーチンはまだ生きている...）")

	// ゴルーチンは ch <- "メッセージ2" の行で永遠に待ち続ける
	// = ゴルーチンリーク
	time.Sleep(1 * time.Second)
}

// ゴルーチンリークの修正版
func exampleGoroutineLeakGood() {
	fmt.Println("\n=== ゴルーチンリーク（良い例）===")

	ch := make(chan string)

	go func() {
		fmt.Println("[goroutine] 送信します...")
		ch <- "メッセージ1"
		fmt.Println("[goroutine] 1つ目送信完了")

		ch <- "メッセージ2"
		fmt.Println("[goroutine] 2つ目送信完了")
	}()

	// 2つとも受信する
	msg1 := <-ch
	msg2 := <-ch
	fmt.Printf("[main] 受信: %s, %s\n", msg1, msg2)
	fmt.Println("[main] 終了（ゴルーチンも正常終了）")
}

func runBlockDemo() {
	// exampleBlockDemo()
	exampleBufferedChannel()
	exampleGoroutineLeakBad()
	exampleGoroutineLeakGood()
}
