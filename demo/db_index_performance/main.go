package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("データベースインデックスのパフォーマンス比較デモ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// インデックスの基本情報を表示
	PrintIndexInfo()

	// データ件数の設定
	productCount := 10000
	orderCount := 10000
	employeeCount := 5000

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("データベースのセットアップ")
	fmt.Println(strings.Repeat("=", 80))

	// インデックスなしのDBをセットアップ
	fmt.Println("\n[1/2] インデックスなしのDBをセットアップ中...")
	dbNoIndex, err := InitDB("no_index.db", false)
	if err != nil {
		fmt.Printf("データベース初期化エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SetupNoIndexDB(dbNoIndex); err != nil {
		fmt.Printf("マイグレーションエラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedProducts(dbNoIndex, productCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedOrders(dbNoIndex, orderCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedEmployees(dbNoIndex, employeeCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	// インデックスありのDBをセットアップ
	fmt.Println("\n[2/2] インデックスありのDBをセットアップ中...")
	dbWithIndex, err := InitDB("with_index.db", false)
	if err != nil {
		fmt.Printf("データベース初期化エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SetupIndexedDB(dbWithIndex); err != nil {
		fmt.Printf("マイグレーションエラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedProductsIndexed(dbWithIndex, productCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedOrdersIndexed(dbWithIndex, orderCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	if err := SeedEmployeesIndexed(dbWithIndex, employeeCount); err != nil {
		fmt.Printf("データ生成エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ セットアップ完了！")

	var results []BenchmarkResult

	// ============================================
	// ベンチマーク 1: カテゴリで商品を検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 1】カテゴリで商品を検索 (WHERE category_id = ?)")
	fmt.Println(strings.Repeat("=", 80))

	categoryID := 2
	fmt.Printf("\nカテゴリID %d の商品を検索...\n", categoryID)

	resultNoIndex1 := RunBenchmark("カテゴリ検索（インデックスなし）", func() (int, error) {
		products, err := FindProductsByCategory(dbNoIndex, categoryID)
		return len(products), err
	})
	results = append(results, resultNoIndex1)

	resultWithIndex1 := RunBenchmark("カテゴリ検索（インデックスあり）", func() (int, error) {
		products, err := FindProductsByCategoryIndexed(dbWithIndex, categoryID)
		return len(products), err
	})
	results = append(results, resultWithIndex1)

	CompareBenchmarks(resultNoIndex1, resultWithIndex1)

	// ============================================
	// ベンチマーク 2: 価格範囲で商品を検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 2】価格範囲で商品を検索 (WHERE price BETWEEN ? AND ?)")
	fmt.Println(strings.Repeat("=", 80))

	minPrice, maxPrice := 5000, 7000
	fmt.Printf("\n価格が %d〜%d の商品を検索...\n", minPrice, maxPrice)

	resultNoIndex2 := RunBenchmark("価格範囲検索（インデックスなし）", func() (int, error) {
		products, err := FindProductsByPriceRange(dbNoIndex, minPrice, maxPrice)
		return len(products), err
	})
	results = append(results, resultNoIndex2)

	resultWithIndex2 := RunBenchmark("価格範囲検索（インデックスあり）", func() (int, error) {
		products, err := FindProductsByPriceRangeIndexed(dbWithIndex, minPrice, maxPrice)
		return len(products), err
	})
	results = append(results, resultWithIndex2)

	CompareBenchmarks(resultNoIndex2, resultWithIndex2)

	// ============================================
	// ベンチマーク 3: ユーザーIDで注文を検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 3】ユーザーIDで注文を検索 (WHERE user_id = ?)")
	fmt.Println(strings.Repeat("=", 80))

	userID := 500
	fmt.Printf("\nユーザーID %d の注文を検索...\n", userID)

	resultNoIndex3 := RunBenchmark("ユーザーID検索（インデックスなし）", func() (int, error) {
		orders, err := FindOrdersByUserID(dbNoIndex, userID)
		return len(orders), err
	})
	results = append(results, resultNoIndex3)

	resultWithIndex3 := RunBenchmark("ユーザーID検索（インデックスあり）", func() (int, error) {
		orders, err := FindOrdersByUserIDIndexed(dbWithIndex, userID)
		return len(orders), err
	})
	results = append(results, resultWithIndex3)

	CompareBenchmarks(resultNoIndex3, resultWithIndex3)

	// ============================================
	// ベンチマーク 4: ステータスで注文を検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 4】ステータスで注文を検索 (WHERE status = ?)")
	fmt.Println(strings.Repeat("=", 80))

	status := "shipped"
	fmt.Printf("\nステータスが '%s' の注文を検索...\n", status)

	resultNoIndex4 := RunBenchmark("ステータス検索（インデックスなし）", func() (int, error) {
		orders, err := FindOrdersByStatus(dbNoIndex, status)
		return len(orders), err
	})
	results = append(results, resultNoIndex4)

	resultWithIndex4 := RunBenchmark("ステータス検索（インデックスあり）", func() (int, error) {
		orders, err := FindOrdersByStatusIndexed(dbWithIndex, status)
		return len(orders), err
	})
	results = append(results, resultWithIndex4)

	CompareBenchmarks(resultNoIndex4, resultWithIndex4)

	// ============================================
	// ベンチマーク 5: 複合インデックス - 名前で従業員を検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 5】名前で従業員を検索（複合インデックス）")
	fmt.Println("(WHERE first_name = ? AND last_name = ?)")
	fmt.Println(strings.Repeat("=", 80))

	firstName, lastName := "John", "Smith"
	fmt.Printf("\n名前が '%s %s' の従業員を検索...\n", firstName, lastName)

	resultNoIndex5 := RunBenchmark("名前検索（インデックスなし）", func() (int, error) {
		employees, err := FindEmployeesByName(dbNoIndex, firstName, lastName)
		return len(employees), err
	})
	results = append(results, resultNoIndex5)

	resultWithIndex5 := RunBenchmark("名前検索（複合インデックスあり）", func() (int, error) {
		employees, err := FindEmployeesByNameIndexed(dbWithIndex, firstName, lastName)
		return len(employees), err
	})
	results = append(results, resultWithIndex5)

	CompareBenchmarks(resultNoIndex5, resultWithIndex5)

	// ============================================
	// ベンチマーク 6: 複合インデックス - 部署と役職で検索
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 6】部署と役職で従業員を検索（複合インデックス）")
	fmt.Println("(WHERE department = ? AND position = ?)")
	fmt.Println(strings.Repeat("=", 80))

	department, position := "Engineering", "Senior"
	fmt.Printf("\n部署が '%s'、役職が '%s' の従業員を検索...\n", department, position)

	resultNoIndex6 := RunBenchmark("部署・役職検索（インデックスなし）", func() (int, error) {
		employees, err := FindEmployeesByDepartmentAndPosition(dbNoIndex, department, position)
		return len(employees), err
	})
	results = append(results, resultNoIndex6)

	resultWithIndex6 := RunBenchmark("部署・役職検索（複合インデックスあり）", func() (int, error) {
		employees, err := FindEmployeesByDepartmentAndPositionIndexed(dbWithIndex, department, position)
		return len(employees), err
	})
	results = append(results, resultWithIndex6)

	CompareBenchmarks(resultNoIndex6, resultWithIndex6)

	// ============================================
	// ベンチマーク 7: ソート処理 - 価格順
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 7】価格でソート (ORDER BY price DESC)")
	fmt.Println(strings.Repeat("=", 80))

	limit := 100
	fmt.Printf("\n商品を価格順にソート（上位%d件）...\n", limit)

	resultNoIndex7 := RunBenchmark("価格ソート（インデックスなし）", func() (int, error) {
		products, err := FindProductsSortedByPrice(dbNoIndex, limit)
		return len(products), err
	})
	results = append(results, resultNoIndex7)

	resultWithIndex7 := RunBenchmark("価格ソート（インデックスあり）", func() (int, error) {
		products, err := FindProductsSortedByPriceIndexed(dbWithIndex, limit)
		return len(products), err
	})
	results = append(results, resultWithIndex7)

	CompareBenchmarks(resultNoIndex7, resultWithIndex7)

	// ============================================
	// ベンチマーク 8: ソート処理 - 日付順
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【ベンチマーク 8】日付でソート (ORDER BY order_date DESC)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("\n注文を日付順にソート（上位%d件）...\n", limit)

	resultNoIndex8 := RunBenchmark("日付ソート（インデックスなし）", func() (int, error) {
		orders, err := FindOrdersSortedByDate(dbNoIndex, limit)
		return len(orders), err
	})
	results = append(results, resultNoIndex8)

	resultWithIndex8 := RunBenchmark("日付ソート（インデックスあり）", func() (int, error) {
		orders, err := FindOrdersSortedByDateIndexed(dbWithIndex, limit)
		return len(orders), err
	})
	results = append(results, resultWithIndex8)

	CompareBenchmarks(resultNoIndex8, resultWithIndex8)

	// ============================================
	// 全体の結果を表示
	// ============================================
	PrintBenchmarkResults(results)

	// まとめを表示
	PrintSummary()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("デモ完了！")
	fmt.Println(strings.Repeat("=", 80))
}
