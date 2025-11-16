package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 環境変数から接続情報を取得（デフォルト値あり）
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "go_demo")

	// データベース接続
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PostgreSQL EXPLAIN デモプログラム")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	// 例1: 基本的なSELECT
	example1(db)

	// 例2: WHERE句でのフィルタリング
	example2(db)

	// 例3: インデックスの効果を比較
	example3(db)

	// 例4: JOIN操作
	example4(db)

	// 例5: 集約クエリ
	example5(db)

	// 例6: 大量データでのパフォーマンス比較
	example6(db)

	// 例7: 複合インデックス
	example7(db)
}

// printResult は実行計画の結果を見やすく表示
func printResult(title string, rows *sql.Rows) {
	fmt.Println("【" + title + "】")
	fmt.Println(strings.Repeat("-", 80))

	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		fmt.Println(plan)
	}
	fmt.Println()
}

// executeExplain は EXPLAIN クエリを実行して結果を表示
func executeExplain(db *sql.DB, title, query string, analyze bool) {
	explainQuery := "EXPLAIN "
	if analyze {
		explainQuery = "EXPLAIN ANALYZE "
	}
	explainQuery += query

	rows, err := db.Query(explainQuery)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	printResult(title, rows)
}

// example1: 基本的なSELECT
func example1(db *sql.DB) {
	fmt.Println("\n■ 例1: 基本的なSELECT - Seq Scan の確認")
	fmt.Println("=" + strings.Repeat("=", 80))

	// EXPLAIN のみ（実行しない）
	executeExplain(db,
		"実行計画のみ（EXPLAIN）",
		"SELECT * FROM users",
		false)

	// EXPLAIN ANALYZE（実際に実行）
	executeExplain(db,
		"実行計画 + 実測値（EXPLAIN ANALYZE）",
		"SELECT * FROM users",
		true)
}

// example2: WHERE句でのフィルタリング
func example2(db *sql.DB) {
	fmt.Println("\n■ 例2: WHERE句でのフィルタリング")
	fmt.Println("=" + strings.Repeat("=", 80))

	executeExplain(db,
		"特定ユーザーのTODO検索",
		"SELECT * FROM todos WHERE user_id = 1",
		true)

	executeExplain(db,
		"完了フラグでのフィルタリング",
		"SELECT * FROM todos WHERE done_flag = false",
		true)
}

// example3: インデックスの効果を比較
func example3(db *sql.DB) {
	fmt.Println("\n■ 例3: インデックスの効果を確認")
	fmt.Println("=" + strings.Repeat("=", 80))

	// インデックスが存在するか確認
	var indexExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE tablename = 'todos' AND indexname = 'idx_todos_user_id'
		)
	`).Scan(&indexExists)
	if err != nil {
		log.Printf("Error checking index: %v", err)
		return
	}

	if indexExists {
		fmt.Println("✓ インデックス idx_todos_user_id が存在します")
		fmt.Println()

		executeExplain(db,
			"インデックスを使った検索",
			"SELECT * FROM todos WHERE user_id = 2",
			true)
	} else {
		fmt.Println("✗ インデックス idx_todos_user_id が存在しません")
		fmt.Println("以下のコマンドでインデックスを作成できます:")
		fmt.Println("  CREATE INDEX idx_todos_user_id ON todos(user_id);")
		fmt.Println()

		executeExplain(db,
			"インデックスなしの検索（Seq Scan）",
			"SELECT * FROM todos WHERE user_id = 2",
			true)
	}
}

// example4: JOIN操作
func example4(db *sql.DB) {
	fmt.Println("\n■ 例4: JOIN操作の実行計画")
	fmt.Println("=" + strings.Repeat("=", 80))

	executeExplain(db,
		"TODOとユーザーのJOIN",
		`SELECT t.title, u.name
		 FROM todos t
		 JOIN users u ON t.user_id = u.id
		 WHERE t.done_flag = false`,
		true)

	executeExplain(db,
		"LEFT JOINの例",
		`SELECT u.name, COUNT(t.id) as todo_count
		 FROM users u
		 LEFT JOIN todos t ON u.id = t.user_id
		 GROUP BY u.id, u.name`,
		true)
}

// example5: 集約クエリ
func example5(db *sql.DB) {
	fmt.Println("\n■ 例5: 集約クエリの実行計画")
	fmt.Println("=" + strings.Repeat("=", 80))

	executeExplain(db,
		"ユーザーごとのTODO件数",
		`SELECT user_id, COUNT(*) as count
		 FROM todos
		 GROUP BY user_id`,
		true)

	executeExplain(db,
		"集約結果のソート",
		`SELECT user_id, COUNT(*) as count
		 FROM todos
		 GROUP BY user_id
		 ORDER BY count DESC`,
		true)
}

// example6: 大量データでのパフォーマンス比較
func example6(db *sql.DB) {
	fmt.Println("\n■ 例6: 大量データでのパフォーマンス")
	fmt.Println("=" + strings.Repeat("=", 80))

	// テーブル内のレコード数を確認
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	if err != nil {
		log.Printf("Error counting rows: %v", err)
		return
	}

	fmt.Printf("現在のTODOレコード数: %d件\n\n", count)

	if count < 1000 {
		fmt.Println("💡 ヒント: 大量データでの違いを体感するには、以下のコマンドでデータを追加してください:")
		fmt.Println()
		fmt.Println("  psql -U postgres -d go_demo -c \"")
		fmt.Println("    INSERT INTO todos (user_id, title, note, done_flag)")
		fmt.Println("    SELECT")
		fmt.Println("      (random() * 3 + 1)::int,")
		fmt.Println("      'TODO ' || generate_series,")
		fmt.Println("      'Note ' || generate_series,")
		fmt.Println("      random() > 0.5")
		fmt.Println("    FROM generate_series(1, 100000);")
		fmt.Println("  \"")
		fmt.Println()
	}

	executeExplain(db,
		"大量データからの検索",
		"SELECT * FROM todos WHERE user_id = 3",
		true)
}

// example7: 複合インデックス
func example7(db *sql.DB) {
	fmt.Println("\n■ 例7: 複合インデックスの活用")
	fmt.Println("=" + strings.Repeat("=", 80))

	// 複合インデックスの存在確認
	var indexExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE tablename = 'todos' AND indexname = 'idx_todos_user_done'
		)
	`).Scan(&indexExists)
	if err != nil {
		log.Printf("Error checking index: %v", err)
		return
	}

	if indexExists {
		fmt.Println("✓ 複合インデックス idx_todos_user_done が存在します")
		fmt.Println()
	} else {
		fmt.Println("✗ 複合インデックス idx_todos_user_done が存在しません")
		fmt.Println("以下のコマンドで作成できます:")
		fmt.Println("  CREATE INDEX idx_todos_user_done ON todos(user_id, done_flag);")
		fmt.Println()
	}

	executeExplain(db,
		"複合条件での検索",
		"SELECT * FROM todos WHERE user_id = 1 AND done_flag = false",
		true)

	executeExplain(db,
		"Index Only Scan の例（カラムを絞る）",
		"SELECT user_id, done_flag FROM todos WHERE user_id = 1",
		true)
}
