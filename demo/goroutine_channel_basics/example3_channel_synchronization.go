package main

import (
	"fmt"
	"time"
)

// ============================================
// 3. channelを使った複数goroutineの同期
// ============================================

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d: ジョブ %d を処理中...\n", id, job)
		time.Sleep(500 * time.Millisecond)
		results <- job * 2 // 結果を送信
	}
}

func example3_ChannelSynchronization() {
	fmt.Println("\n=== Example 3: channelで複数goroutineを同期 ===")

	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// 3つのworkerを起動
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 5つのジョブを送信
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // これ以上ジョブを送らないことを示す

	// 5つの結果を受信
	for a := 1; a <= 5; a++ {
		result := <-results
		fmt.Printf("結果: %d\n", result)
	}

	fmt.Println("Example 3 完了")
}
