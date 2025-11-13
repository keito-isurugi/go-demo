package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("二分探索木（BST）vs B-Tree の比較")
	fmt.Println(strings.Repeat("=", 80))

	// ========================================
	// 構造の違い
	// ========================================
	fmt.Println("\n【1】構造の違い")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n■ 二分探索木（Binary Search Tree）")
	fmt.Println("各ノードが最大2つの子を持つ")
	fmt.Println()
	fmt.Println("例: 7つの要素 [1, 2, 3, 4, 5, 6, 7] を挿入")
	fmt.Println()
	fmt.Println("        4")
	fmt.Println("      /   \\")
	fmt.Println("     2     6")
	fmt.Println("    / \\   / \\")
	fmt.Println("   1   3 5   7")
	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  - 各ノード: 1つのキーのみ")
	fmt.Println("  - 子ノード: 最大2つ（左と右）")
	fmt.Println("  - 高さ: log₂(n)")

	fmt.Println("\n■ B-Tree（B木）")
	fmt.Println("各ノードが複数のキーと複数の子を持つ")
	fmt.Println()
	fmt.Println("例: 同じ7つの要素、次数=4の場合")
	fmt.Println()
	fmt.Println("        [  4  ]")
	fmt.Println("       /       \\")
	fmt.Println("  [1, 2, 3]  [5, 6, 7]")
	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  - 各ノード: 複数のキー（最大 m-1 個）")
	fmt.Println("  - 子ノード: 複数（最大 m 個）")
	fmt.Println("  - 高さ: log_m(n)（mは次数）")

	// ========================================
	// ディスクアクセスの違い
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【2】ディスクアクセスの違い（これが最重要！）")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n💾 ディスクの特性:")
	fmt.Println("  - 1回の読み取り: 約10ms（機械的な遅延）")
	fmt.Println("  - メモリの読み取り: 約0.0001ms")
	fmt.Println("  → ディスクアクセスはメモリの100,000倍遅い！")
	fmt.Println("  → ディスクアクセスの回数を減らすことが最優先")

	fmt.Println("\n■ 100万件のデータから1件を検索する場合")
	fmt.Println()

	// 二分探索木の場合
	bstHeight := 20 // log₂(1,000,000) ≈ 20
	fmt.Printf("二分探索木（BST）:\n")
	fmt.Printf("  高さ: log₂(1,000,000) ≈ %d\n", bstHeight)
	fmt.Printf("  ディスクアクセス: %d回\n", bstHeight)
	fmt.Printf("  所要時間: %d回 × 10ms = %dms\n", bstHeight, bstHeight*10)
	fmt.Println()
	fmt.Println("  各ノードが1つのキーしか持たないため、")
	fmt.Println("  1回のディスクアクセスで1つのキーしか確認できない")

	// B-Treeの場合
	btreeOrder := 1000 // 次数1000（1ノードに999個のキー）
	btreeHeight := 3   // log₁₀₀₀(1,000,000) ≈ 2-3
	fmt.Printf("B-Tree（次数=%d）:\n", btreeOrder)
	fmt.Printf("  高さ: log₁₀₀₀(1,000,000) ≈ %d\n", btreeHeight)
	fmt.Printf("  ディスクアクセス: %d回\n", btreeHeight)
	fmt.Printf("  所要時間: %d回 × 10ms = %dms\n", btreeHeight, btreeHeight*10)
	fmt.Println()
	fmt.Println("  各ノードが999個のキーを持つため、")
	fmt.Println("  1回のディスクアクセスで999個のキーを確認できる")

	fmt.Println()
	fmt.Printf("📊 速度差: B-Treeは二分探索木の約 %.1f倍速い！\n",
		float64(bstHeight*10)/float64(btreeHeight*10))

	// ========================================
	// ノードサイズとディスクブロック
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【3】ノードサイズとディスクブロックの関係")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n💾 ディスクの物理的な制約:")
	fmt.Println("  - ディスクブロックサイズ: 4KB〜16KB")
	fmt.Println("  - 1回のディスクアクセスで1ブロック分読める")
	fmt.Println("  → 小さなデータを何度も読むより、")
	fmt.Println("    大きなデータをまとめて1回で読む方が効率的")

	fmt.Println("\n■ 二分探索木のノード（小さい）")
	fmt.Println()
	fmt.Println("  struct BSTNode {")
	fmt.Println("      int key;           // 4 bytes")
	fmt.Println("      void* left;        // 8 bytes")
	fmt.Println("      void* right;       // 8 bytes")
	fmt.Println("  }  // 合計: 20 bytes")
	fmt.Println()
	fmt.Println("  4KBのブロックを読んでも、1つのノード（20bytes）しか使わない")
	fmt.Println("  → 残りの約4000bytesが無駄！")

	fmt.Println("\n■ B-Treeのノード（大きい）")
	fmt.Println()
	fmt.Println("  struct BTreeNode {")
	fmt.Println("      int keys[999];     // 4KB")
	fmt.Println("      void* children[1000]; // 8KB")
	fmt.Println("  }  // 合計: 12KB")
	fmt.Println()
	fmt.Println("  ディスクブロックをほぼ埋め尽くす")
	fmt.Println("  → 1回のディスクアクセスで999個のキーを取得できる！")

	// ========================================
	// 実際の使い分け
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【4】実際の使い分け")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n✅ 二分探索木（BST）が向いている場面:")
	fmt.Println("  - メモリ上のデータ構造")
	fmt.Println("  - データ量が少ない（数千件程度）")
	fmt.Println("  - ディスクアクセスが不要")
	fmt.Println("  - 実装がシンプル")
	fmt.Println()
	fmt.Println("  例:")
	fmt.Println("    - Goのmap（実際はハッシュテーブル）")
	fmt.Println("    - プログラムのシンボルテーブル")
	fmt.Println("    - 一時的なソート済みデータ")

	fmt.Println("\n✅ B-Treeが向いている場面:")
	fmt.Println("  - ディスク上のデータ構造")
	fmt.Println("  - データ量が多い（数百万件以上）")
	fmt.Println("  - ディスクアクセスを最小化したい")
	fmt.Println("  - データベースのインデックス")
	fmt.Println()
	fmt.Println("  例:")
	fmt.Println("    - PostgreSQL、MySQL、SQLiteのインデックス")
	fmt.Println("    - ファイルシステム（ext4、NTFS、BtrFS）")
	fmt.Println("    - NoSQLデータベース（MongoDB）")

	// ========================================
	// バランスの問題
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【5】バランスの問題")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n■ 二分探索木の問題: 偏りやすい")
	fmt.Println()
	fmt.Println("昇順にデータを挿入した場合:")
	fmt.Println()
	fmt.Println("  挿入順: 1, 2, 3, 4, 5")
	fmt.Println()
	fmt.Println("  結果:")
	fmt.Println("    1")
	fmt.Println("     \\")
	fmt.Println("      2")
	fmt.Println("       \\")
	fmt.Println("        3")
	fmt.Println("         \\")
	fmt.Println("          4")
	fmt.Println("           \\")
	fmt.Println("            5")
	fmt.Println()
	fmt.Println("  → リンクリストと同じ（O(n)の探索時間）")
	fmt.Println("  → 対策: AVL木、赤黒木（バランス調整が複雑）")

	fmt.Println("\n■ B-Treeの利点: 常にバランスが保たれる")
	fmt.Println()
	fmt.Println("同じデータを挿入した場合:")
	fmt.Println()
	fmt.Println("  挿入順: 1, 2, 3, 4, 5")
	fmt.Println()
	fmt.Println("  結果（次数=3）:")
	fmt.Println("      [  3  ]")
	fmt.Println("     /       \\")
	fmt.Println("  [1, 2]   [4, 5]")
	fmt.Println()
	fmt.Println("  → 自動的にバランスが保たれる")
	fmt.Println("  → 全ての葉ノードが同じ深さ")
	fmt.Println("  → 常にO(log n)が保証される")

	// ========================================
	// 範囲検索の違い
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【6】範囲検索の違い")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n■ クエリ: SELECT * FROM users WHERE age BETWEEN 20 AND 30")
	fmt.Println()

	fmt.Println("二分探索木:")
	fmt.Println("  1. age=20 を探す")
	fmt.Println("  2. 中間順序走査（in-order traversal）で次々と訪問")
	fmt.Println("  3. age>30 になったら停止")
	fmt.Println("  → ツリーを上下に移動する必要がある")
	fmt.Println()

	fmt.Println("B+Tree（B-Treeの改良版）:")
	fmt.Println("  1. age=20 を探す（葉ノードへ）")
	fmt.Println("  2. 葉ノード同士がリンクされているので、")
	fmt.Println("     右にたどるだけで範囲内のデータを全て取得")
	fmt.Println("  3. age>30 になったら停止")
	fmt.Println("  → ツリーを上下に移動する必要がない（高速！）")
	fmt.Println()
	fmt.Println("B+Treeの葉ノード:")
	fmt.Println("  [18,19] → [20,21] → [22,23] → ... → [30,31]")
	fmt.Println("            ↑開始              ↑終了")

	// ========================================
	// まとめ
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【まとめ】なぜデータベースはB-Treeを使うのか？")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n理由1: ディスクアクセス回数が少ない")
	fmt.Println("  二分探索木: log₂(n) 回")
	fmt.Println("  B-Tree:    log_m(n) 回（mは次数、通常100〜1000）")
	fmt.Println("  → 100万件で約20回 vs 約3回")

	fmt.Println("\n理由2: ディスクブロックを効率的に使う")
	fmt.Println("  二分探索木: 1ノード ≈ 20 bytes（ブロックの0.5%）")
	fmt.Println("  B-Tree:    1ノード ≈ 4-16 KB（ブロックを最大活用）")

	fmt.Println("\n理由3: 常にバランスが保たれる")
	fmt.Println("  二分探索木: 偏る可能性あり（最悪O(n)）")
	fmt.Println("  B-Tree:    常にバランス（必ずO(log n)）")

	fmt.Println("\n理由4: 範囲検索が高速")
	fmt.Println("  二分探索木: ツリーを上下に移動")
	fmt.Println("  B+Tree:    葉ノードを横にスキャン")

	fmt.Println("\n理由5: キャッシュ効率が良い")
	fmt.Println("  二分探索木: 各ノードが小さく、キャッシュミスが多い")
	fmt.Println("  B-Tree:    各ノードが大きく、一度に多くのデータを取得")

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("結論:")
	fmt.Println("  メモリ上のデータ → 二分探索木でOK")
	fmt.Println("  ディスク上のデータ → B-Treeが圧倒的に有利！")
	fmt.Println(strings.Repeat("=", 80))
}
