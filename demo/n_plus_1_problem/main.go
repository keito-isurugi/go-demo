package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("N+1問題のパフォーマンス比較デモ")
	fmt.Println("================================\n")

	// データベースを初期化（SQLログは非表示）
	db, err := InitDB(false)
	if err != nil {
		fmt.Printf("データベース初期化エラー: %v\n", err)
		os.Exit(1)
	}

	// サンプルデータを生成
	// ユーザー: 10人、投稿/ユーザー: 5件、コメント/投稿: 3件
	if err := SeedData(db, 10, 5, 3); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	var results []BenchmarkResult

	// ============================================
	// ベンチマーク 1: ユーザーと投稿の取得
	// ============================================
	fmt.Println("【ベンチマーク 1】ユーザーと投稿の取得")
	fmt.Println(strings.Repeat("-", 80))

	// N+1問題あり
	fmt.Println("\n1-1. N+1問題のあるコード実行中...")
	resultBad1 := RunBenchmark("ユーザーと投稿 (N+1問題あり)", func() (int, error) {
		users, err := FetchUsersWithPostsBad(db)
		return len(users), err
	})
	results = append(results, resultBad1)

	// 最適化版
	fmt.Println("1-2. 最適化されたコード実行中...")
	resultGood1 := RunBenchmark("ユーザーと投稿 (Preload使用)", func() (int, error) {
		users, err := FetchUsersWithPostsGood(db)
		return len(users), err
	})
	results = append(results, resultGood1)

	CompareBenchmarks(resultBad1, resultGood1)

	// ============================================
	// ベンチマーク 2: 投稿とコメントの取得
	// ============================================
	fmt.Println("\n【ベンチマーク 2】投稿とコメントの取得")
	fmt.Println(strings.Repeat("-", 80))

	// N+1問題あり
	fmt.Println("\n2-1. N+1問題のあるコード実行中...")
	resultBad2 := RunBenchmark("投稿とコメント (N+1問題あり)", func() (int, error) {
		posts, err := FetchPostsWithCommentsBad(db)
		return len(posts), err
	})
	results = append(results, resultBad2)

	// 最適化版
	fmt.Println("2-2. 最適化されたコード実行中...")
	resultGood2 := RunBenchmark("投稿とコメント (Preload使用)", func() (int, error) {
		posts, err := FetchPostsWithCommentsGood(db)
		return len(posts), err
	})
	results = append(results, resultGood2)

	CompareBenchmarks(resultBad2, resultGood2)

	// ============================================
	// ベンチマーク 3: ユーザー、投稿、コメントの取得（3階層）
	// ============================================
	fmt.Println("\n【ベンチマーク 3】ユーザー、投稿、コメントの取得（3階層）")
	fmt.Println(strings.Repeat("-", 80))

	// N+1問題あり（多段階）
	fmt.Println("\n3-1. N+1問題のあるコード実行中（多段階）...")
	resultBad3 := RunBenchmark("3階層データ (N+1問題あり)", func() (int, error) {
		users, err := FetchUsersWithPostsAndCommentsBad(db)
		return len(users), err
	})
	results = append(results, resultBad3)

	// 最適化版（ネストしたPreload）
	fmt.Println("3-2. 最適化されたコード実行中（ネストしたPreload）...")
	resultGood3 := RunBenchmark("3階層データ (ネストしたPreload)", func() (int, error) {
		users, err := FetchUsersWithPostsAndCommentsGood(db)
		return len(users), err
	})
	results = append(results, resultGood3)

	CompareBenchmarks(resultBad3, resultGood3)

	// ============================================
	// 全体のベンチマーク結果を表示
	// ============================================
	PrintBenchmarkResults(results)

	// ============================================
	// SQLログ付きで実行例を表示
	// ============================================
	fmt.Println("\n【SQLログ付き実行例】")
	fmt.Println("N+1問題のあるコードで発行されるSQL:")
	fmt.Println(strings.Repeat("-", 80))

	// SQLログを有効にして再初期化
	dbWithLog, _ := InitDB(true)
	if err := SeedData(dbWithLog, 3, 2, 2); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n▼ N+1問題のあるコード（ユーザー3人、投稿2件/人）:")
	_, _ = FetchUsersWithPostsBad(dbWithLog)

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("▼ 最適化されたコード（同じデータ）:")
	_, _ = FetchUsersWithPostsGood(dbWithLog)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("デモ完了！")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\n【まとめ】")
	fmt.Println("・N+1問題: 各レコードごとに追加のクエリが発行される")
	fmt.Println("・Preloadを使用することで、必要最小限のクエリで全データを取得可能")
	fmt.Println("・パフォーマンスが大幅に改善される（特にレコード数が多い場合）")
}
