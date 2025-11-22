package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// チャネルの同期機構としての利点
// ============================================

// 1. 作業完了の待機（done pattern）
func exampleDonePattern() {
	fmt.Println("\n=== Done Pattern ===")

	done := make(chan bool)

	go func() {
		fmt.Println("[worker] 作業開始...")
		time.Sleep(2 * time.Second)
		fmt.Println("[worker] 作業完了！")
		done <- true  // 完了を通知
	}()

	fmt.Println("[main] 作業の完了を待っています...")
	<-done  // ← ここでブロック（完了まで待つ）
	fmt.Println("[main] 作業が完了したので次に進めます")
}

// 2. 複数の作業の完了を待つ
func exampleMultipleDone() {
	fmt.Println("\n=== 複数の作業の完了待機 ===")

	done := make(chan string, 3)  // バッファ付き

	// 3つの作業を並行実行
	go workerWithDuration("A", 3*time.Second, done)
	go workerWithDuration("B", 2*time.Second, done)
	go workerWithDuration("C", 1*time.Second, done)

	// 3つすべての完了を待つ
	for i := 0; i < 3; i++ {
		msg := <-done
		fmt.Printf("[main] 受信: %s\n", msg)
	}

	fmt.Println("[main] すべての作業が完了しました")
}

func workerWithDuration(name string, duration time.Duration, done chan string) {
	fmt.Printf("[worker %s] 開始 (所要時間: %v)\n", name, duration)
	time.Sleep(duration)
	fmt.Printf("[worker %s] 完了\n", name)
	done <- fmt.Sprintf("worker %s 完了", name)
}

// 3. データ競合の例（悪い例）
func exampleDataRaceBad() {
	fmt.Println("\n=== データ競合（悪い例）===")

	var counter int

	// 100個のゴルーチンで同時にカウントアップ
	for i := 0; i < 100; i++ {
		go func() {
			counter++  // ← データ競合！
		}()
	}

	time.Sleep(100 * time.Millisecond)  // 適当に待つ（不確実）
	fmt.Printf("counter = %d (期待値: 100)\n", counter)
	fmt.Println("→ 値が100にならないことがある（データ競合）")
}

// 4. データ競合の解決（チャネル版）
func exampleDataRaceGoodChannel() {
	fmt.Println("\n=== データ競合の解決（チャネル版）===")

	ch := make(chan int, 100)  // バッファ付き

	// 100個のゴルーチンで値を送信
	for i := 0; i < 100; i++ {
		go func() {
			ch <- 1  // チャネルに送信（安全）
		}()
	}

	// すべて受信してカウント
	counter := 0
	for i := 0; i < 100; i++ {
		counter += <-ch
	}

	fmt.Printf("counter = %d (期待値: 100)\n", counter)
	fmt.Println("→ 必ず100になる（安全）")
}

// 5. データ競合の解決（Mutex版）
func exampleDataRaceGoodMutex() {
	fmt.Println("\n=== データ競合の解決（Mutex版）===")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 100個のゴルーチンで同時にカウントアップ
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			counter++  // ← ロックで保護
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("counter = %d (期待値: 100)\n", counter)
	fmt.Println("→ 必ず100になる（Mutexで安全）")
}

// 6. パイプライン（データの流れ）
func examplePipeline() {
	fmt.Println("\n=== パイプライン ===")

	// チャネルで3段階のパイプラインを作る
	numbers := make(chan int)
	squares := make(chan int)
	sums := make(chan int)

	// ステージ1: 数値を生成
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("  [生成] %d\n", i)
			numbers <- i
		}
		close(numbers)  // 終了を通知
	}()

	// ステージ2: 2乗する
	go func() {
		for n := range numbers {  // チャネルが閉じるまで受信
			squared := n * n
			fmt.Printf("  [2乗] %d -> %d\n", n, squared)
			squares <- squared
		}
		close(squares)
	}()

	// ステージ3: 合計する
	go func() {
		sum := 0
		for s := range squares {
			sum += s
			fmt.Printf("  [合計] +%d = %d\n", s, sum)
		}
		sums <- sum
		close(sums)
	}()

	// 最終結果を受信
	result := <-sums
	fmt.Printf("[main] 最終結果: 1²+2²+3²+4²+5² = %d\n", result)
}

// 7. デッドロック検出のデモ
func exampleDeadlockDetection() {
	fmt.Println("\n=== デッドロック検出 ===")
	fmt.Println("デッドロックを発生させると Go がエラーを出してくれます")
	fmt.Println("（このデモは実際には実行しません）")

	// 以下のコードは実行するとデッドロックエラーになる
	/*
	ch := make(chan int)
	ch <- 42  // ← 受信者がいないので永遠に待つ
	          //    main goroutine しかないので全体がブロック
	          //    → Go がデッドロックを検出してパニック
	*/

	fmt.Println(`
コード例:
    ch := make(chan int)
    ch <- 42  // デッドロック！

エラーメッセージ:
    fatal error: all goroutines are asleep - deadlock!
`)
}

func runSyncBenefitsDemo() {
	exampleDonePattern()
	exampleMultipleDone()
	exampleDataRaceBad()
	exampleDataRaceGoodChannel()
	exampleDataRaceGoodMutex()
	examplePipeline()
	exampleDeadlockDetection()
}
