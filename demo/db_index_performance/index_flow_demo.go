package main

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OrderWithIndex インデックス付きのテーブル
type OrderWithIndex struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     int       `gorm:"index:idx_user_id"` // user_idにインデックス
	TotalPrice int
	Status     string
	OrderDate  time.Time
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("インデックスを使ったクエリの実行フロー")
	fmt.Println(strings.Repeat("=", 80))

	// データベースのセットアップ
	db, _ := gorm.Open(sqlite.Open("index_flow_demo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.AutoMigrate(&OrderWithIndex{})
	db.Exec("DELETE FROM order_with_indices")

	// サンプルデータを挿入
	fmt.Println("\n【ステップ1】データを挿入")
	fmt.Println(strings.Repeat("-", 80))

	orders := []OrderWithIndex{
		{ID: 1, UserID: 500, TotalPrice: 1000, Status: "pending", OrderDate: time.Now()},
		{ID: 2, UserID: 123, TotalPrice: 2000, Status: "shipped", OrderDate: time.Now()},
		{ID: 3, UserID: 500, TotalPrice: 3000, Status: "delivered", OrderDate: time.Now()},
		{ID: 4, UserID: 789, TotalPrice: 4000, Status: "pending", OrderDate: time.Now()},
		{ID: 5, UserID: 500, TotalPrice: 5000, Status: "cancelled", OrderDate: time.Now()},
		{ID: 6, UserID: 123, TotalPrice: 6000, Status: "shipped", OrderDate: time.Now()},
	}

	for _, order := range orders {
		db.Create(&order)
		fmt.Printf("挿入: ID=%d, UserID=%d, Price=%d, Status=%s\n",
			order.ID, order.UserID, order.TotalPrice, order.Status)
	}

	// メインテーブルの状態
	fmt.Println("\n【ステップ2】メインテーブル（order_with_indices）の状態")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("ディスク上のデータ配置（挿入順）:")
	fmt.Println()
	fmt.Println("┌────┬────────┬────────────┬───────────┐")
	fmt.Println("│ ID │ UserID │ TotalPrice │ Status    │")
	fmt.Println("├────┼────────┼────────────┼───────────┤")
	for _, order := range orders {
		fmt.Printf("│ %2d │ %6d │ %10d │ %-9s │\n",
			order.ID, order.UserID, order.TotalPrice, order.Status)
	}
	fmt.Println("└────┴────────┴────────────┴───────────┘")

	// インデックス（B+Tree）の状態
	fmt.Println("\n【ステップ3】インデックス（idx_user_id）の内部構造")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("B+Treeの構造（user_idでソートされている）:")
	fmt.Println()
	fmt.Println("                  [500]")
	fmt.Println("                 /     \\")
	fmt.Println("           [123]        [789]")
	fmt.Println("            ↓            ↓")
	fmt.Println("       ┌─────────┐   ┌──────┐")
	fmt.Println("       │ UserID  │   │UserID│")
	fmt.Println("       │   123   │   │ 500  │   789")
	fmt.Println("       ├─────────┤   ├──────┤   ├──────┤")
	fmt.Println("       │RowID: 2 │   │RowID:│   │RowID:│")
	fmt.Println("       │ RowID:6 │   │  1   │   │  4   │")
	fmt.Println("       └─────────┘   │  3   │   └──────┘")
	fmt.Println("            ↓        │  5   │       ↓")
	fmt.Println("       [行2, 行6]    └──────┘   [行4]")
	fmt.Println("                        ↓")
	fmt.Println("                  [行1, 行3, 行5]")

	// クエリの実行フロー
	fmt.Println("\n【ステップ4】クエリ実行: WHERE user_id = 500")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("\nSQLクエリ:")
	fmt.Println("  SELECT * FROM order_with_indices WHERE user_id = 500;")

	fmt.Println("\n実行フロー:")
	fmt.Println()
	fmt.Println("1️⃣  クエリオプティマイザがインデックス idx_user_id を選択")
	fmt.Println("    → EXPLAIN QUERY PLAN を見ると:")
	fmt.Println("      SEARCH order_with_indices USING INDEX idx_user_id (user_id=?)")
	fmt.Println()
	fmt.Println("2️⃣  B+Treeで user_id = 500 を検索")
	fmt.Println("    ルートノード [500] を見る → 一致！")
	fmt.Println("    → 対応するリーフノードへ移動")
	fmt.Println()
	fmt.Println("3️⃣  リーフノードから行IDのリストを取得")
	fmt.Println("    user_id = 500 → [RowID: 1, 3, 5]")
	fmt.Println()
	fmt.Println("4️⃣  各行IDを使ってメインテーブルから実データを取得")
	fmt.Println("    RowID 1 → ID=1, UserID=500, Price=1000, Status=pending")
	fmt.Println("    RowID 3 → ID=3, UserID=500, Price=3000, Status=delivered")
	fmt.Println("    RowID 5 → ID=5, UserID=500, Price=5000, Status=cancelled")
	fmt.Println()
	fmt.Println("5️⃣  結果を返す")

	// 実際にクエリを実行
	fmt.Println("\n【ステップ5】実際のクエリ実行結果")
	fmt.Println(strings.Repeat("-", 80))

	var result []OrderWithIndex
	db.Where("user_id = ?", 500).Find(&result)

	fmt.Println("\n取得されたレコード:")
	fmt.Println()
	fmt.Println("┌────┬────────┬────────────┬───────────┐")
	fmt.Println("│ ID │ UserID │ TotalPrice │ Status    │")
	fmt.Println("├────┼────────┼────────────┼───────────┤")
	for _, order := range result {
		fmt.Printf("│ %2d │ %6d │ %10d │ %-9s │\n",
			order.ID, order.UserID, order.TotalPrice, order.Status)
	}
	fmt.Println("└────┴────────┴────────────┴───────────┘")
	fmt.Printf("\n件数: %d件\n", len(result))

	// インデックスなしの場合との比較
	fmt.Println("\n【比較】インデックスなしの場合")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("\n実行フロー:")
	fmt.Println()
	fmt.Println("1️⃣  クエリオプティマイザがフルテーブルスキャンを選択")
	fmt.Println("    → EXPLAIN QUERY PLAN を見ると:")
	fmt.Println("      SCAN order_with_indices")
	fmt.Println()
	fmt.Println("2️⃣  メインテーブルの全行をスキャン")
	fmt.Println("    行1: UserID=500 → ✅ 一致（結果に追加）")
	fmt.Println("    行2: UserID=123 → ❌ 不一致（スキップ）")
	fmt.Println("    行3: UserID=500 → ✅ 一致（結果に追加）")
	fmt.Println("    行4: UserID=789 → ❌ 不一致（スキップ）")
	fmt.Println("    行5: UserID=500 → ✅ 一致（結果に追加）")
	fmt.Println("    行6: UserID=123 → ❌ 不一致（スキップ）")
	fmt.Println()
	fmt.Println("3️⃣  結果を返す")
	fmt.Println()
	fmt.Println("📊 比較:")
	fmt.Println("  インデックスあり: 3回のアクセス（RowID 1, 3, 5 のみ）")
	fmt.Println("  インデックスなし: 6回のアクセス（全行をスキャン）")
	fmt.Println("  → データ量が多いほど差が大きくなる！")

	// 実際のEXPLAIN QUERY PLANを表示
	fmt.Println("\n【ステップ6】実際のクエリプランを確認")
	fmt.Println(strings.Repeat("-", 80))

	sqlDB, _ := db.DB()
	rows, _ := sqlDB.Query("EXPLAIN QUERY PLAN SELECT * FROM order_with_indices WHERE user_id = 500")

	fmt.Println("\nEXPLAIN QUERY PLAN の結果:")
	var id, parent, notused int
	var detail string
	for rows.Next() {
		rows.Scan(&id, &parent, &notused, &detail)
		fmt.Printf("  %s\n", detail)
	}
	rows.Close()

	if strings.Contains(detail, "USING INDEX") {
		fmt.Println("\n✅ インデックス idx_user_id が使われている！")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("まとめ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nインデックスを使ったクエリの流れ:")
	fmt.Println("  1. WHERE句で user_id が指定される")
	fmt.Println("  2. データベースが idx_user_id インデックスを選択")
	fmt.Println("  3. B+Treeで user_id の値を検索（高速）")
	fmt.Println("  4. 該当する行IDのリストを取得")
	fmt.Println("  5. 行IDを使ってメインテーブルから実データを取得")
	fmt.Println("  6. 結果を返す")
	fmt.Println("\n→ 全行をスキャンする必要がないので高速！")
}
