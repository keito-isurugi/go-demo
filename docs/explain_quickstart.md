# PostgreSQL EXPLAIN クイックスタート

5分で始める実践的なEXPLAINの使い方

## 今すぐ始める

### 1. データベースに接続

```bash
# PostgreSQLコンテナに接続
docker exec -it go-demo-db psql -U postgres -d go_demo
```

### 2. 最初のEXPLAINを実行

```sql
-- シンプルなクエリの実行計画を見る
EXPLAIN SELECT * FROM users;
```

**出力例:**
```
Seq Scan on users  (cost=0.00..10.50 rows=50 width=1576)
```

これだけ！あなたは今、実行計画を見ることができました。

### 3. 実際に実行してみる

```sql
-- クエリを実際に実行して、実測値も見る
EXPLAIN ANALYZE SELECT * FROM users;
```

**出力例:**
```
Seq Scan on users  (cost=0.00..10.50 rows=50 width=1576)
                   (actual time=0.007..0.009 rows=4 loops=1)
Planning Time: 0.398 ms
Execution Time: 0.041 ms
```

これで推定値と実測値の両方が見れます！

---

## 実行計画を読む（3つのポイント）

### ポイント1: スキャン方法

```sql
EXPLAIN SELECT * FROM todos WHERE user_id = 1;
```

- **Seq Scan** → 全件スキャン（遅いかも）
- **Index Scan** → インデックス使用（速い！）
- **Index Only Scan** → 最速！

### ポイント2: 実行時間

```sql
EXPLAIN ANALYZE SELECT * FROM todos;
```

**注目:**
```
Execution Time: 1.234 ms
```

この数字が重要！大きすぎる場合は最適化の余地あり。

### ポイント3: 推定 vs 実測

```
rows=100 width=...  (actual time=... rows=10000 loops=1)
 ^^^推定値^^^                         ^^^^実測値^^^^
```

推定と実測が大きく違う場合は要注意！

---

## 実践: インデックスで高速化

### Before（インデックスなし）

```sql
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;
```

**結果:**
```
Seq Scan on todos  (cost=0.00..268.00 ...)
Execution Time: 5.123 ms
```

遅い...

### After（インデックス作成）

```sql
-- インデックス作成
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- 同じクエリを再実行
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;
```

**結果:**
```
Index Scan using idx_todos_user_id on todos
Execution Time: 0.234 ms
```

速くなった！🚀

---

## よく使うコマンド集

### 基本

```sql
-- 実行計画のみ
EXPLAIN <クエリ>;

-- 実行計画 + 実測値
EXPLAIN ANALYZE <クエリ>;

-- 詳細情報も表示
EXPLAIN (ANALYZE, BUFFERS, VERBOSE) <クエリ>;
```

### インデックス管理

```sql
-- インデックス一覧
\di

-- インデックス作成
CREATE INDEX idx_column ON table_name(column);

-- インデックス削除
DROP INDEX idx_column;

-- 統計情報更新（インデックス作成後は必ず実行！）
ANALYZE table_name;
```

### テーブル情報

```sql
-- テーブル構造
\d table_name

-- レコード数
SELECT COUNT(*) FROM table_name;

-- テーブルサイズ
SELECT pg_size_pretty(pg_total_relation_size('table_name'));
```

---

## 5分間チャレンジ

実際に手を動かして理解を深めましょう！

### チャレンジ1: 基本

```sql
-- 1. usersテーブルの実行計画を見る
EXPLAIN SELECT * FROM users;

-- 2. 実際に実行してみる
EXPLAIN ANALYZE SELECT * FROM users;

-- 3. 実行時間を確認する
--    → Execution Time: ??? ms
```

### チャレンジ2: WHERE句

```sql
-- 1. user_id=1のTODOを検索
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;

-- 2. どのスキャン方法が使われている？
--    → Seq Scan? Index Scan?

-- 3. 実行時間は？
--    → Execution Time: ??? ms
```

### チャレンジ3: インデックスの効果

```sql
-- 1. インデックスが存在するか確認
SELECT indexname FROM pg_indexes
WHERE tablename = 'todos' AND indexname = 'idx_todos_user_id';

-- 2. 存在しない場合は作成
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- 3. 統計情報更新
ANALYZE todos;

-- 4. もう一度実行計画を確認
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;

-- 5. 実行計画は変わった？
```

---

## トラブルシューティング

### Q: 「relation does not exist」エラー

```sql
-- テーブルが存在するか確認
\dt

-- 存在しない場合はスキーマを適用
-- bash
make refresh-schema
```

### Q: インデックスを作ったのにSeq Scanが使われる

**原因1:** テーブルが小さすぎる
- 小規模データではSeq Scanの方が速いことがある
- これは正常な動作

**原因2:** 統計情報が古い
```sql
-- 統計情報を更新
ANALYZE todos;
```

**原因3:** 大部分の行を取得している
```sql
-- 例: 全行の90%を取得する場合、Seq Scanの方が効率的
SELECT * FROM todos WHERE user_id IN (1,2,3,4);
```

### Q: 推定行数と実測行数が全然違う

```sql
-- 統計情報を更新
ANALYZE todos;

-- それでも違う場合は、統計の精度を上げる
ALTER TABLE todos ALTER COLUMN user_id SET STATISTICS 1000;
ANALYZE todos;
```

---

## 次のステップ

### 📚 詳しく学ぶ

完全ガイドで深く学びましょう:
- [PostgreSQL EXPLAIN完全ガイド](./postgresql_explain_guide.md)

### 🎯 実践演習

SQLスクリプトで手を動かして学びましょう:
```bash
cd demo/explain_demo
cat examples.sql
```

### 🔧 プロジェクトで使う

1. 遅いクエリを見つける
2. `EXPLAIN ANALYZE`で原因を特定
3. インデックスやクエリを改善
4. 効果を測定

---

## チートシート

### 読み方

| 表示 | 意味 |
|------|------|
| `Seq Scan` | 全件スキャン |
| `Index Scan` | インデックス使用 |
| `Index Only Scan` | インデックスのみで完結（最速） |
| `cost=0.00..10.50` | 起動コスト..総コスト |
| `rows=100` | 推定行数 |
| `actual time=0.01..0.05` | 実測時間 |
| `rows=95` | 実際の行数 |
| `Execution Time: 1.23 ms` | 実行時間 |

### パフォーマンス判断

| 実行時間 | 評価 | 対応 |
|---------|------|------|
| < 10 ms | ✅ 良好 | OK |
| 10-100 ms | ⚠️ やや遅い | 改善検討 |
| > 100 ms | ❌ 遅い | 要改善 |

### よくある問題と解決策

| 問題 | 解決策 |
|------|--------|
| Seq Scanが使われている | インデックス作成 |
| 推定と実測が大きく違う | `ANALYZE`実行 |
| Nested Loopが遅い | 内側のテーブルにインデックス |
| Sortでディスク使用 | `work_mem`増やす or インデックス |

---

## まとめ

🎉 これで基本は完璧！

覚えておくべき3つのコマンド:

1. **`EXPLAIN ANALYZE <クエリ>`** - 実行計画を見る
2. **`CREATE INDEX`** - インデックス作成
3. **`ANALYZE`** - 統計情報更新

あとは実践あるのみ！

---

## 参考リソース

- [完全ガイド](./postgresql_explain_guide.md) - より詳しい説明
- [実践演習SQL](../demo/explain_demo/examples.sql) - 手を動かして学ぶ
- [PostgreSQL公式ドキュメント](https://www.postgresql.org/docs/current/sql-explain.html)
- [explain.depesz.com](https://explain.depesz.com/) - 実行計画の可視化ツール
