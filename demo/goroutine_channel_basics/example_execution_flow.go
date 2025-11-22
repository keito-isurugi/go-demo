package main

import (
	"fmt"
	"time"
)

// ============================================
// チャネル送信の実行フローを詳しく見る
// ============================================

func exampleExecutionFlow() {
	fmt.Println("\n=== チャネル送信の実行フロー ===")

	ch := make(chan string)

	// ゴルーチンを起動
	go func() {
		fmt.Println("[goroutine] ステップ1: 関数開始")
		fmt.Println("[goroutine] ステップ2: チャネル送信を試みます...")

		// ここで実行が「開始」されるが「完了」しない
		fmt.Println("[goroutine] ステップ3: ch <- を実行中... (次の行に進まない)")
		ch <- "メッセージ"  // ← この行は実行されるが、完了までブロックされる

		// 上の行が完了するまで、ここには到達しない
		fmt.Println("[goroutine] ステップ4: 送信完了！次の行に進めた")
		fmt.Println("[goroutine] ステップ5: 関数終了")
	}()

	// main は意図的に待つ
	fmt.Println("[main] ステップ1: 3秒待ちます（ゴルーチンは送信中...）")
	time.Sleep(3 * time.Second)

	fmt.Println("[main] ステップ2: 今から受信します")
	msg := <-ch  // ← ここで初めて送信が完了できる
	fmt.Printf("[main] ステップ3: 受信完了: %s\n", msg)

	// ゴルーチンの残りの処理が実行される
	time.Sleep(100 * time.Millisecond)  // ゴルーチンの出力を待つ
}

// 実行と完了の違いをもっと明確に
func exampleExecutionVsCompletion() {
	fmt.Println("\n=== 実行 vs 完了 ===")

	ch := make(chan string)

	go func() {
		fmt.Println("[goroutine] 行1: fmt.Println実行 → すぐ完了")
		fmt.Println("[goroutine] 行2: time.Sleep実行開始...")
		time.Sleep(1 * time.Second)
		fmt.Println("[goroutine] 行2: time.Sleep完了")

		fmt.Println("[goroutine] 行3: ch <- 実行開始... (受信者待ち)")
		startTime := time.Now()
		ch <- "データ"  // ← 実行されるが完了しない
		elapsed := time.Since(startTime)

		fmt.Printf("[goroutine] 行3: ch <- 完了 (待機時間: %v)\n", elapsed)
		fmt.Println("[goroutine] 行4: すべて完了")
	}()

	fmt.Println("[main] 5秒待ってから受信します...")
	time.Sleep(5 * time.Second)

	fmt.Println("[main] 受信開始")
	<-ch
	fmt.Println("[main] 受信完了")

	time.Sleep(100 * time.Millisecond)
}

// CPUを使わずに待つことを確認
func exampleCPUUsage() {
	fmt.Println("\n=== CPU使用率の確認 ===")

	ch := make(chan string)

	// このゴルーチンは待機中にCPUを使わない
	go func() {
		fmt.Println("[goroutine] 送信します（ブロックされます）")
		ch <- "メッセージ"
		fmt.Println("[goroutine] 送信完了")
	}()

	fmt.Println("[main] 10秒待ちます（この間、ゴルーチンはCPUを使いません）")
	fmt.Println("[main] top コマンドなどでCPU使用率を見てみてください")

	// カウントダウン表示
	for i := 10; i > 0; i-- {
		fmt.Printf("[main] 残り %d秒...\n", i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("[main] 受信します")
	<-ch
	fmt.Println("[main] 完了")

	time.Sleep(100 * time.Millisecond)
}

// 対比: busyループ（CPUを使い続ける）
func exampleBusyLoop() {
	fmt.Println("\n=== 対比: Busyループ（悪い例）===")

	done := false

	// こちらはCPUを使い続ける（悪い実装）
	go func() {
		fmt.Println("[goroutine] busyループ開始（CPUを使い続けます）")
		count := 0
		for !done {  // ← 何度も何度もチェックし続ける（CPUを消費）
			count++
			if count % 100000000 == 0 {
				fmt.Println("[goroutine] まだ待っています...")
			}
		}
		fmt.Println("[goroutine] ループ終了")
	}()

	fmt.Println("[main] 3秒待ちます（この間、ゴルーチンはCPUを使い続けます）")
	time.Sleep(3 * time.Second)

	done = true
	fmt.Println("[main] done=true に設定しました")

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n【比較】")
	fmt.Println("- チャネル: 待機中はCPUを使わない（効率的）")
	fmt.Println("- busyループ: 待機中もCPUを使い続ける（非効率）")
}

func runExecutionFlowDemo() {
	exampleExecutionFlow()
	exampleExecutionVsCompletion()
	// exampleCPUUsage()  // 時間がかかるのでコメントアウト
	exampleBusyLoop()
}
