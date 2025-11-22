# PostgreSQL EXPLAIN完全ガイド - 実務で使えるクエリ最適化入門

## 目次
1. [EXPLAINとは](#explainとは)
2. [基本的な使い方](#基本的な使い方)
3. [EXPLAIN ANALYZEの違い](#explain-analyzeの違い)
4. [実行計画の読み方](#実行計画の読み方)
5. [主要なノードタイプ](#主要なノードタイプ)
6. [実践例](#実践例)
7. [パフォーマンスチューニング](#パフォーマンスチューニング)

---

## EXPLAINとは

`EXPLAIN`は、PostgreSQLがSQLクエリをどのように実行するかを示すコマンドです。クエリの実行計画（クエリプラン）を確認することで、以下のことができます:

- クエリが遅い原因を特定
- インデックスが正しく使われているか確認
- JOIN操作の効率性を評価
- クエリの最適化ポイントを発見

**重要な概念:**
- **実行計画**: PostgreSQLがクエリを実行するために選択した手順
- **コスト**: クエリ実行の相対的な計算量（単位は任意）
- **rows**: 推定される行数
- **width**: 各行の平均バイト数

---

## 基本的な使い方

### 1. EXPLAIN（実行計画のみ表示）

```sql
EXPLAIN SELECT * FROM users;
```

**出力例:**
```
Seq Scan on users  (cost=0.00..15.20 rows=520 width=48)
```

クエリは実際には実行されず、計画だけが表示されます。

### 2. EXPLAIN ANALYZE（実際に実行して結果も表示）

```sql
EXPLAIN ANALYZE SELECT * FROM users;
```

**出力例:**
```
Seq Scan on users  (cost=0.00..15.20 rows=520 width=48)
                   (actual time=0.012..0.015 rows=4 loops=1)
Planning Time: 0.089 ms
Execution Time: 0.032 ms
```

**注意:** `EXPLAIN ANALYZE`は実際にクエリを実行するため、INSERT/UPDATE/DELETEでは注意が必要です！
（トランザクション内で実行してROLLBACKすると安全です）

### 3. 詳細オプション

```sql
-- より詳細な情報を表示
EXPLAIN (ANALYZE, BUFFERS, VERBOSE) SELECT * FROM users;

-- JSON形式で出力
EXPLAIN (FORMAT JSON) SELECT * FROM users;

-- YAML形式で出力
EXPLAIN (FORMAT YAML) SELECT * FROM users;
```

**オプション一覧:**
- `ANALYZE`: 実際に実行して実測値も表示
- `BUFFERS`: バッファの使用状況を表示
- `VERBOSE`: 詳細情報を表示
- `COSTS`: コスト情報を表示（デフォルトON）
- `TIMING`: 実行時間を計測（デフォルトON）
- `FORMAT`: 出力形式を指定（TEXT/JSON/YAML/XML）

---

## EXPLAIN ANALYZEの違い

| コマンド | クエリ実行 | 表示内容 | 用途 |
|---------|----------|---------|------|
| `EXPLAIN` | しない | 推定値のみ | 実行計画の事前確認 |
| `EXPLAIN ANALYZE` | する | 推定値 + 実測値 | 実際のパフォーマンス計測 |

**推定値と実測値の比較:**
```sql
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 1;
```

```
Seq Scan on todos  (cost=0.00..18.50 rows=4 width=...)  -- 推定
                   (actual time=0.015..0.018 rows=4 loops=1)  -- 実測
                    ^^^^^^^^^推定^^^^^^^   ^^^^^^^^実測^^^^^^^^
```

---

## 実行計画の読み方

### コスト表示の読み方

```
Seq Scan on users  (cost=0.00..15.20 rows=520 width=48)
                           ^^^^  ^^^^^ ^^^^^^^^ ^^^^^^^^
                           起動  総計  推定行数  平均幅
```

- **起動コスト (0.00)**: 最初の行を返すまでのコスト
- **総コスト (15.20)**: すべての行を返すまでのコスト
- **rows (520)**: このノードが返すと推定される行数
- **width (48)**: 各行の平均バイト数

### 実測時間の読み方

```
(actual time=0.012..0.015 rows=4 loops=1)
             ^^^^^  ^^^^^ ^^^^^^ ^^^^^^^
             開始   終了  実際の  実行
             時間   時間  行数    回数
```

- **actual time**: 実際にかかった時間（ミリ秒）
  - 最初の値: 最初の行を返すまでの時間
  - 2番目の値: すべての行を返し終わるまでの時間
- **rows**: 実際に返された行数
- **loops**: このノードが実行された回数

---

## 主要なノードタイプ

### 1. Seq Scan（順次スキャン）

テーブル全体を先頭から順番に読む。

```sql
EXPLAIN SELECT * FROM todos;
```

```
Seq Scan on todos  (cost=0.00..18.50 rows=850 width=...)
```

**特徴:**
- インデックスを使わない
- 小規模テーブルや、大部分の行を取得する場合は効率的
- 大規模テーブルで少数の行を取得する場合は遅い

### 2. Index Scan（インデックススキャン）

インデックスを使って必要な行を効率的に取得。

```sql
-- 前提: user_idにインデックスが存在
EXPLAIN SELECT * FROM todos WHERE user_id = 1;
```

```
Index Scan using todos_user_id_idx on todos
    (cost=0.15..8.17 rows=4 width=...)
    Index Cond: (user_id = 1)
```

**特徴:**
- WHERE句の条件に合致する行を効率的に検索
- インデックスとテーブルの両方にアクセス
- 少数の行を取得する場合に最適

### 3. Index Only Scan（インデックスオンリースキャン）

インデックスだけで結果を返せる場合の最速スキャン。

```sql
-- インデックスに含まれるカラムのみ取得
CREATE INDEX idx_user_id ON todos(user_id);
EXPLAIN SELECT user_id FROM todos WHERE user_id > 1;
```

```
Index Only Scan using idx_user_id on todos
    (cost=0.15..4.17 rows=4 width=4)
    Index Cond: (user_id > 1)
```

**特徴:**
- テーブルにアクセスせずインデックスだけで完結
- 最も高速
- VACUUMが適切に実行されていることが重要

### 4. Bitmap Index Scan / Bitmap Heap Scan

複数のインデックスを組み合わせる場合や、多数の行を取得する場合に使用。

```sql
EXPLAIN SELECT * FROM todos WHERE user_id IN (1, 2, 3);
```

```
Bitmap Heap Scan on todos
    -> Bitmap Index Scan on todos_user_id_idx
         Index Cond: (user_id = ANY ('{1,2,3}'::integer[]))
```

**特徴:**
- まずビットマップを作成してから行を取得
- ランダムアクセスを減らし、効率的にディスクを読む
- 中規模の結果セットに適している

### 5. Nested Loop（ネステッドループ）

JOINの基本的な方法。外側のテーブルの各行について、内側のテーブルを検索。

```sql
EXPLAIN SELECT * FROM todos t
JOIN users u ON t.user_id = u.id;
```

```
Nested Loop  (cost=0.15..45.23 rows=850 width=...)
    -> Seq Scan on users u
    -> Index Scan using todos_user_id_idx on todos t
         Index Cond: (user_id = u.id)
```

**特徴:**
- 内側のループにインデックスがある場合に効率的
- 結合する行数が少ない場合に最適
- 外側の行数 × 内側の検索コスト

### 6. Hash Join（ハッシュジョイン）

大量のデータを結合する場合に効率的。

```sql
EXPLAIN SELECT * FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.done_flag = false;
```

```
Hash Join  (cost=15.35..35.58 rows=425 width=...)
    Hash Cond: (t.user_id = u.id)
    -> Seq Scan on todos t
         Filter: (done_flag = false)
    -> Hash
         -> Seq Scan on users u
```

**特徴:**
- 一方のテーブルでハッシュテーブルを作成
- 大量のデータの結合に適している
- メモリを使用（work_mem設定に依存）

### 7. Merge Join（マージジョイン）

ソート済みのデータを結合する場合に効率的。

```sql
EXPLAIN SELECT * FROM todos t
JOIN users u ON t.user_id = u.id
ORDER BY t.user_id;
```

```
Merge Join  (cost=45.12..65.34 rows=850 width=...)
    Merge Cond: (t.user_id = u.id)
    -> Sort
         Sort Key: t.user_id
         -> Seq Scan on todos t
    -> Sort
         Sort Key: u.id
         -> Seq Scan on users u
```

**特徴:**
- 両方のテーブルをソートしてマージ
- すでにソートされている（インデックスがある）場合に最適
- 大量のデータでも安定したパフォーマンス

### 8. Aggregate（集約）

GROUP BYやCOUNT、SUM等の集約関数。

```sql
EXPLAIN SELECT user_id, COUNT(*)
FROM todos
GROUP BY user_id;
```

```
HashAggregate  (cost=20.63..22.63 rows=200 width=12)
    Group Key: user_id
    -> Seq Scan on todos
```

**特徴:**
- HashAggregate: ハッシュテーブルを使用（高速）
- GroupAggregate: ソート済みデータをグループ化

### 9. Sort（ソート）

ORDER BY句で使用。

```sql
EXPLAIN SELECT * FROM todos ORDER BY created_at DESC;
```

```
Sort  (cost=52.33..54.45 rows=850 width=...)
    Sort Key: created_at DESC
    -> Seq Scan on todos
```

**特徴:**
- メモリ内でソート（work_mem設定に依存）
- メモリ不足の場合はディスクを使用（遅い）
- インデックスがあればスキップ可能

---

## 実践例

このセクションでは、実際にプロジェクトのデータベースを使って練習できます。

### 準備: データベースの起動とデータ投入

```bash
# Dockerコンテナ起動
docker compose up -d db

# スキーマ適用とダミーデータ投入
make refresh-schema

# PostgreSQLに接続
docker exec -it go-demo-db psql -U postgres -d go_demo
```

### 例1: 基本的なSELECTの分析

```sql
-- 全件取得の実行計画
EXPLAIN ANALYZE SELECT * FROM users;
```

**読み解き方:**
1. Seq Scanが使われている → テーブル全体を読んでいる
2. rows=4（実測値）→ 4件のデータが存在
3. Execution Time → 実際にかかった時間

### 例2: WHERE句でのフィルタリング

```sql
-- 特定ユーザーのTODOを取得
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;
```

**最適化前（インデックスなし）:**
```
Seq Scan on todos  (cost=0.00..18.50 rows=4 width=...)
    Filter: (user_id = 1)
    Rows Removed by Filter: 6
```
→ 全行をスキャンしてフィルタリング（非効率）

**最適化: インデックス作成**
```sql
CREATE INDEX idx_todos_user_id ON todos(user_id);

EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;
```

**最適化後:**
```
Index Scan using idx_todos_user_id on todos
    (cost=0.15..8.17 rows=4 width=...)
    Index Cond: (user_id = 1)
```
→ インデックスを使用して効率的に検索

### 例3: JOIN操作の分析

```sql
-- TODOとユーザーを結合
EXPLAIN ANALYZE
SELECT t.title, u.name
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.done_flag = false;
```

**分析ポイント:**
- どのJOIN方式が使われているか（Nested Loop / Hash / Merge）
- 結合の順序は適切か
- インデックスは使われているか

### 例4: 集約クエリの分析

```sql
-- ユーザーごとのTODO件数
EXPLAIN ANALYZE
SELECT u.name, COUNT(t.id) as todo_count
FROM users u
LEFT JOIN todos t ON u.id = t.user_id
GROUP BY u.id, u.name
ORDER BY todo_count DESC;
```

**分析ポイント:**
- HashAggregate vs GroupAggregate
- Sortが必要か
- インデックスでソートを省略できるか

### 例5: サブクエリの分析

```sql
-- TODOが3件以上あるユーザーを検索
EXPLAIN ANALYZE
SELECT name FROM users
WHERE id IN (
    SELECT user_id
    FROM todos
    GROUP BY user_id
    HAVING COUNT(*) >= 3
);
```

**分析ポイント:**
- サブクエリがどのように実行されているか
- EXISTS句に書き換えると改善するか

### 例6: 複合インデックスの活用

```sql
-- 複合インデックスの作成
CREATE INDEX idx_todos_user_done ON todos(user_id, done_flag);

-- 複合条件での検索
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false;
```

**インデックスの使われ方を確認:**
- Index Condに両方の条件が含まれているか
- Filterで追加フィルタリングされていないか

### 例7: LIKE検索の最適化

```sql
-- 前方一致検索（インデックス使用可能）
EXPLAIN ANALYZE
SELECT * FROM users WHERE name LIKE '山田%';

-- 部分一致検索（インデックス使用不可）
EXPLAIN ANALYZE
SELECT * FROM users WHERE name LIKE '%太郎%';
```

**最適化方法:**
```sql
-- 全文検索用のインデックス（GIN）
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_users_name_gin ON users USING gin(name gin_trgm_ops);

EXPLAIN ANALYZE
SELECT * FROM users WHERE name LIKE '%太郎%';
```

### 例8: 大量データでの実験

実際の性能差を体感するため、大量データを投入してみましょう。

```sql
-- 大量のTODOデータを投入（10万件）
INSERT INTO todos (user_id, title, note, done_flag)
SELECT
    (random() * 3 + 1)::int,  -- user_id 1-4
    'TODO ' || generate_series,
    'Note ' || generate_series,
    random() > 0.5
FROM generate_series(1, 100000);

-- インデックスなしでの検索
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;

-- インデックスありでの検索
CREATE INDEX idx_todos_user_id ON todos(user_id);
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;
```

**実行時間の比較:**
- インデックスなし: 数十ミリ秒〜
- インデックスあり: 数ミリ秒

---

## パフォーマンスチューニング

### 1. よくある問題と対策

#### 問題1: Seq Scanが使われている

**原因:**
- インデックスがない
- 統計情報が古い
- テーブルが小さすぎる（Seq Scanの方が速い）
- WHERE句の条件が複雑すぎる

**対策:**
```sql
-- インデックス作成
CREATE INDEX idx_column ON table_name(column);

-- 統計情報の更新
ANALYZE table_name;

-- 統計情報の確認
SELECT * FROM pg_stats WHERE tablename = 'todos';
```

#### 問題2: 推定行数と実測行数が大きく異なる

**例:**
```
Seq Scan on todos  (cost=0.00..18.50 rows=10 width=...)
                   (actual time=0.015..0.018 rows=100000 loops=1)
                    ^^^^^^^^^推定10行^^^^^^^   ^^^^^実際は10万行^^^^^
```

**原因:**
- 統計情報が古い
- データの偏りが大きい

**対策:**
```sql
-- 統計情報を更新
ANALYZE todos;

-- 統計情報の詳細度を上げる
ALTER TABLE todos ALTER COLUMN user_id SET STATISTICS 1000;
ANALYZE todos;
```

#### 問題3: Nested Loopが遅い

**原因:**
- 内側のテーブルにインデックスがない
- 外側のテーブルの行数が多すぎる

**対策:**
```sql
-- 内側のテーブルにインデックスを作成
CREATE INDEX idx_inner_column ON inner_table(join_column);

-- JOINの順序を変更（ヒント句を使用）
SET enable_nestloop = off;  -- 一時的に無効化して別の方式を試す
```

#### 問題4: ソートでメモリ不足

**メモリ不足の兆候:**
```
Sort  (cost=...)
    Sort Key: created_at
    Sort Method: external merge  Disk: 12345kB  ← ディスク使用
```

**対策:**
```sql
-- work_memを増やす（セッション単位）
SET work_mem = '64MB';

-- インデックスを使ってソートを回避
CREATE INDEX idx_created_at ON todos(created_at);
```

### 2. インデックス戦略

#### 単一カラムインデックス

```sql
-- 頻繁にWHERE句で使うカラム
CREATE INDEX idx_user_id ON todos(user_id);
CREATE INDEX idx_done_flag ON todos(done_flag);
```

#### 複合インデックス

```sql
-- 複数カラムを同時に使う場合
CREATE INDEX idx_user_done ON todos(user_id, done_flag);

-- 注意: カラムの順序が重要！
-- このインデックスは以下のクエリに有効:
-- WHERE user_id = 1
-- WHERE user_id = 1 AND done_flag = false
-- しかし以下には無効:
-- WHERE done_flag = false（user_idがない）
```

#### 部分インデックス

```sql
-- 特定の条件のみインデックス化
CREATE INDEX idx_active_todos
ON todos(user_id)
WHERE done_flag = false;

-- 未完了のTODOのみを扱うクエリで効果的
EXPLAIN SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false;
```

#### カバリングインデックス

```sql
-- クエリで使うすべてのカラムを含める
CREATE INDEX idx_covering
ON todos(user_id, done_flag, title, created_at);

-- Index Only Scanが可能になる
EXPLAIN SELECT title, created_at
FROM todos
WHERE user_id = 1 AND done_flag = false;
```

### 3. 実行計画を改善する設定パラメータ

```sql
-- 現在の設定を確認
SHOW work_mem;
SHOW random_page_cost;
SHOW effective_cache_size;

-- セッション単位で変更
SET work_mem = '64MB';              -- ソート・ハッシュ用メモリ
SET random_page_cost = 1.1;         -- SSD使用時は低めに
SET effective_cache_size = '4GB';   -- システムキャッシュサイズ
```

### 4. 実行計画のチェックリスト

クエリを最適化する際は、以下をチェック:

- [ ] Seq Scanが適切か？（大部分の行を取得する場合はOK）
- [ ] インデックスが使われているか？
- [ ] 推定行数と実測行数は近いか？
- [ ] JOINの順序は適切か？
- [ ] 不要なソートはないか？
- [ ] メモリ不足でディスクを使っていないか？
- [ ] 実行時間は許容範囲内か？

---

## 実践的なワークフロー

### 1. 遅いクエリを見つける

```sql
-- 実行時間が長いクエリをログに記録
ALTER DATABASE go_demo SET log_min_duration_statement = 1000;  -- 1秒以上

-- スロークエリログを確認（Docker環境）
docker logs go-demo-db | grep "duration:"
```

### 2. 実行計画を分析

```sql
EXPLAIN (ANALYZE, BUFFERS)
<遅いクエリ>;
```

### 3. 問題点を特定

- どのノードに時間がかかっているか
- Seq Scanが不適切に使われていないか
- JOINの方法は適切か

### 4. 改善策を実施

- インデックス作成
- クエリの書き換え
- 統計情報の更新

### 5. 効果を検証

```sql
EXPLAIN ANALYZE
<改善後のクエリ>;
```

実行時間とコストが改善されたか確認。

---

## まとめ

### EXPLAINを使いこなすポイント

1. **まずはEXPLAIN ANALYZEで実測する**
   - 推測ではなく実測値で判断

2. **推定行数と実測行数の差に注目**
   - 大きく乖離している場合はANALYZEを実行

3. **コストだけでなく実行時間を見る**
   - 最終的には実行時間が重要

4. **インデックスは万能ではない**
   - 小規模テーブルではSeq Scanの方が速い
   - 更新のオーバーヘッドも考慮

5. **定期的な統計情報の更新**
   - データ量が変わったらANALYZEを実行

### 次のステップ

- pg_stat_statementsでクエリ統計を収集
- pgBadgerでログを分析
- EXPLAIN (ANALYZE, BUFFERS)でバッファ使用状況を分析
- auto_explainモジュールで自動的に遅いクエリを記録

---

## 参考リンク

- [PostgreSQL公式ドキュメント - EXPLAIN](https://www.postgresql.org/docs/current/sql-explain.html)
- [PostgreSQL公式ドキュメント - パフォーマンスチューニング](https://www.postgresql.org/docs/current/performance-tips.html)
- [pgMustard](https://www.pgmustard.com/) - 実行計画の可視化ツール
- [explain.depesz.com](https://explain.depesz.com/) - 実行計画の分析ツール

---

## 付録: 便利なSQLスニペット

### 現在のインデックス一覧

```sql
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename IN ('todos', 'users')
ORDER BY tablename, indexname;
```

### テーブルサイズの確認

```sql
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### インデックスの使用状況

```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;  -- 使われていないインデックスを発見
```

### 未使用のインデックスを削除

```sql
-- 使用されていないインデックスを確認
SELECT
    schemaname || '.' || tablename AS table,
    indexname AS index,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size,
    idx_scan AS scans
FROM pg_stat_user_indexes
WHERE idx_scan = 0
    AND indexrelname NOT LIKE '%pkey'  -- 主キーは除外
ORDER BY pg_relation_size(indexrelid) DESC;
```
