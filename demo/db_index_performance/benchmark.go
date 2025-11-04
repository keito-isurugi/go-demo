package main

import (
	"fmt"
	"strings"
	"time"
)

// BenchmarkResult ベンチマーク結果
type BenchmarkResult struct {
	Name        string
	Duration    time.Duration
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

// CompareBenchmarks 2つのベンチマーク結果を比較
func CompareBenchmarks(noIndex, withIndex BenchmarkResult) {
	improvement := float64(noIndex.Duration) / float64(withIndex.Duration)
	timeSaved := noIndex.Duration - withIndex.Duration

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("改善結果: %s\n", noIndex.Name)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("インデックスなし: %v\n", noIndex.Duration.Round(time.Microsecond))
	fmt.Printf("インデックスあり: %v\n", withIndex.Duration.Round(time.Microsecond))
	fmt.Printf("改善率:           %.2f倍速くなりました\n", improvement)
	fmt.Printf("短縮時間:         %v\n", timeSaved.Round(time.Microsecond))
	fmt.Printf("取得件数:         %d件\n", withIndex.ResultCount)
	fmt.Println(strings.Repeat("=", 80))
}

// PrintBenchmarkResults ベンチマーク結果を一覧表示
func PrintBenchmarkResults(results []BenchmarkResult) {
	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("全体のベンチマーク結果")
	fmt.Println(strings.Repeat("=", 90))
	fmt.Printf("%-50s %20s %15s\n", "テスト名", "実行時間", "取得件数")
	fmt.Println(strings.Repeat("-", 90))

	for _, result := range results {
		fmt.Printf("%-50s %20s %15d\n",
			result.Name,
			result.Duration.Round(time.Microsecond),
			result.ResultCount)
	}
	fmt.Println(strings.Repeat("=", 90))
}

// PrintIndexInfo インデックス情報を表示
func PrintIndexInfo() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("インデックスの種類")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("1. 単一カラムインデックス")
	fmt.Println("   - 1つのカラムにインデックスを作成")
	fmt.Println("   - 例: category_id, price, status など")
	fmt.Println("   - 用途: WHERE句での単一条件検索、ソート処理")
	fmt.Println()
	fmt.Println("2. 複合インデックス")
	fmt.Println("   - 複数のカラムを組み合わせたインデックス")
	fmt.Println("   - 例: (first_name, last_name), (department, position)")
	fmt.Println("   - 用途: 複数条件での検索、カラムの順序が重要")
	fmt.Println("   - 注意: 左端のカラムから順に使用される")
	fmt.Println()
	fmt.Println("3. インデックスが効果的なケース")
	fmt.Println("   - WHERE句での検索")
	fmt.Println("   - ORDER BY句でのソート")
	fmt.Println("   - JOIN操作")
	fmt.Println("   - GROUP BY句での集約")
	fmt.Println()
	fmt.Println("4. インデックスが効果的でないケース")
	fmt.Println("   - LIKE '%keyword%' のような中間一致検索")
	fmt.Println("   - データの選択性が低い場合（例: 性別カラム）")
	fmt.Println("   - 小規模なテーブル")
	fmt.Println(strings.Repeat("=", 80))
}

// PrintSummary まとめを表示
func PrintSummary() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("まとめ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("【インデックスのメリット】")
	fmt.Println("  ✓ 検索速度の大幅な向上")
	fmt.Println("  ✓ ソート処理の高速化")
	fmt.Println("  ✓ 大量データでも高速なアクセス")
	fmt.Println()
	fmt.Println("【インデックスのデメリット】")
	fmt.Println("  ✗ ディスク容量の使用")
	fmt.Println("  ✗ INSERT/UPDATE/DELETE処理のオーバーヘッド")
	fmt.Println("  ✗ インデックスのメンテナンスコスト")
	fmt.Println()
	fmt.Println("【ベストプラクティス】")
	fmt.Println("  • WHERE句で頻繁に使用するカラムにインデックスを作成")
	fmt.Println("  • JOIN条件のカラムにインデックスを作成")
	fmt.Println("  • ソート処理に使用するカラムにインデックスを作成")
	fmt.Println("  • 複数カラムでの検索には複合インデックスを検討")
	fmt.Println("  • 不要なインデックスは削除する")
	fmt.Println("  • EXPLAIN/EXPLAIN ANALYZEで実行計画を確認")
	fmt.Println(strings.Repeat("=", 80))
}
