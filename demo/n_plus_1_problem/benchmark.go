package main

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BenchmarkResult ベンチマーク結果
type BenchmarkResult struct {
	Name        string
	Duration    time.Duration
	QueryCount  int
	ResultCount int
}

// RunBenchmark ベンチマークを実行
func RunBenchmark(name string, fn func() (int, error)) BenchmarkResult {
	start := time.Now()
	count, err := fn()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error in %s: %v\n", name, err)
	}

	return BenchmarkResult{
		Name:        name,
		Duration:    duration,
		ResultCount: count,
	}
}

// PrintBenchmarkResults ベンチマーク結果を表示
func PrintBenchmarkResults(results []BenchmarkResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ベンチマーク結果")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-50s %15s %10s\n", "テスト名", "実行時間", "件数")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		fmt.Printf("%-50s %15s %10d\n",
			result.Name,
			result.Duration.Round(time.Microsecond),
			result.ResultCount)
	}
	fmt.Println(strings.Repeat("=", 80))
}

// CompareBenchmarks 2つのベンチマーク結果を比較
func CompareBenchmarks(bad, good BenchmarkResult) {
	improvement := float64(bad.Duration) / float64(good.Duration)
	timeSaved := bad.Duration - good.Duration

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("改善結果: %s vs %s\n", bad.Name, good.Name)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("N+1問題あり: %v\n", bad.Duration.Round(time.Microsecond))
	fmt.Printf("最適化後:     %v\n", good.Duration.Round(time.Microsecond))
	fmt.Printf("改善率:       %.2f倍速くなりました\n", improvement)
	fmt.Printf("短縮時間:     %v\n", timeSaved.Round(time.Microsecond))
	fmt.Println(strings.Repeat("=", 80))
}

// QueryCounter クエリ数をカウントするためのコールバック
type QueryCounter struct {
	Count int
}

// NewQueryCounter クエリカウンターを作成
func NewQueryCounter(db *gorm.DB) *QueryCounter {
	counter := &QueryCounter{}

	// クエリ実行前のコールバックを登録
	db.Callback().Query().Before("gorm:query").Register("query_counter", func(db *gorm.DB) {
		counter.Count++
	})

	// CREATE/UPDATE/DELETEのコールバックも登録
	db.Callback().Create().Before("gorm:create").Register("create_counter", func(db *gorm.DB) {
		counter.Count++
	})
	db.Callback().Update().Before("gorm:update").Register("update_counter", func(db *gorm.DB) {
		counter.Count++
	})
	db.Callback().Delete().Before("gorm:delete").Register("delete_counter", func(db *gorm.DB) {
		counter.Count++
	})

	return counter
}

// Reset カウンターをリセット
func (c *QueryCounter) Reset() {
	c.Count = 0
}

// GetCount 現在のカウントを取得
func (c *QueryCounter) GetCount() int {
	return c.Count
}
