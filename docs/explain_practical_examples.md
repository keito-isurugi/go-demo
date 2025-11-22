# EXPLAIN実践ガイド - 実行計画の読み方と活用方法

実際にEXPLAINの出力をどう読んで、どう活用するかを具体例で学びます。

## 目次
1. [基本的な読み方](#基本的な読み方)
2. [実践例1: Seq Scanの問題を発見して解決](#実践例1-seq-scanの問題を発見して解決)
3. [実践例2: JOINの最適化](#実践例2-joinの最適化)
4. [実践例3: 推定値と実測値の乖離](#実践例3-推定値と実測値の乖離)
5. [実践例4: インデックスが使われない理由](#実践例4-インデックスが使われない理由)
6. [実践例5: 複合インデックスの活用](#実践例5-複合インデックスの活用)

---

## 基本的な読み方

### ステップ1: 実行計画を取得

```sql
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;
```

### ステップ2: 出力を理解する

```
Seq Scan on todos  (cost=0.00..268.00 rows=10000 width=90)
                   (actual time=0.009..1.309 rows=10000 loops=1)
  Filter: (user_id = 1)
Planning Time: 2.441 ms
Execution Time: 1.656 ms
```

**読み解き方:**

1. **スキャン方法**: `Seq Scan` → 全件スキャン（改善の余地あり）
2. **コスト**: `cost=0.00..268.00` → 総コスト268（相対値）
3. **推定行数**: `rows=10000` → 10,000行返すと推定
4. **実測時間**: `actual time=0.009..1.309` → 実際に1.3ms
5. **実際の行数**: `rows=10000` → 実際に10,000行返した
6. **フィルタ**: `Filter: (user_id = 1)` → WHERE句で絞り込み
7. **実行時間**: `Execution Time: 1.656 ms` → トータル1.66ms

### ステップ3: 問題を特定

❌ **問題点:**
- Seq Scan（全件スキャン）を使用 → 非効率
- Filter段階でフィルタリング → インデックスがあれば改善可能

✅ **改善策:**
- user_idカラムにインデックスを作成

---

## 実践例1: Seq Scanの問題を発見して解決

### シナリオ
「特定ユーザーのTODOリストが遅い」という報告を受けた。

### ステップ1: 現状を把握

```sql
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;
```

**出力:**
```
Seq Scan on todos  (cost=0.00..268.00 rows=25000 width=90)
                   (actual time=0.015..3.234 rows=25123 loops=1)
  Filter: (user_id = 2)
  Rows Removed by Filter: 74877
Planning Time: 0.554 ms
Execution Time: 4.123 ms
```

### 📊 分析

| 項目 | 値 | 意味 |
|-----|-----|-----|
| スキャン方法 | Seq Scan | ❌ 全100,000行を読んでいる |
| 実測時間 | 4.123 ms | やや遅い |
| Rows Removed | 74,877 | ❌ 約75%の行を読んで捨てている |
| 実際の行数 | 25,123 | 必要な行は25%のみ |

**問題:** 100,000行を読んで、75,000行を捨てている（無駄が多い）

### ステップ2: インデックスを作成

```sql
CREATE INDEX idx_todos_user_id ON todos(user_id);
ANALYZE todos;
```

### ステップ3: 改善を確認

```sql
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;
```

**出力:**
```
Index Scan using idx_todos_user_id on todos
    (cost=0.29..856.42 rows=25000 width=90)
    (actual time=0.034..0.523 rows=25123 loops=1)
  Index Cond: (user_id = 2)
Planning Time: 0.234 ms
Execution Time: 0.876 ms
```

### 📊 改善結果

| 項目 | Before | After | 改善率 |
|-----|--------|-------|--------|
| スキャン方法 | Seq Scan | Index Scan | ✅ |
| 実行時間 | 4.123 ms | 0.876 ms | **78%削減** |
| 読む行数 | 100,000 | 25,123 | **75%削減** |

### 💡 学んだこと

1. **Seq Scanは常に悪いわけではない** が、大きなテーブルで少数の行を取得する場合は問題
2. **Rows Removed by Filterに注目** → 無駄な読み込みの指標
3. **インデックス作成後は必ずANALYZE** を実行

---

## 実践例2: JOINの最適化

### シナリオ
「ユーザーとTODOを結合するクエリが遅い」

### ステップ1: 現状を把握

```sql
EXPLAIN ANALYZE
SELECT u.name, t.title
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.done_flag = false;
```

**出力:**
```
Hash Join  (cost=11.65..1893.65 rows=50000 width=534)
           (actual time=0.156..15.234 rows=49876 loops=1)
  Hash Cond: (t.user_id = u.id)
  ->  Seq Scan on todos t  (cost=0.00..1518.00 rows=50000 width=22)
                            (actual time=0.012..8.123 rows=49876 loops=1)
        Filter: (done_flag = false)
        Rows Removed by Filter: 50124
  ->  Hash  (cost=10.50..10.50 rows=50 width=520)
            (actual time=0.089..0.091 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on users u  (cost=0.00..10.50 rows=50 width=520)
                                  (actual time=0.008..0.011 rows=4 loops=1)
Planning Time: 0.445 ms
Execution Time: 16.789 ms
```

### 📊 分析（読み方を段階的に）

#### レベル1: 全体構造を把握

```
Hash Join                          ← 最終的な結合方法
  -> Seq Scan on todos (外側)      ← todosテーブルを全件読む
  -> Hash                          ← usersでハッシュテーブル作成
     -> Seq Scan on users (内側)   ← usersテーブルを全件読む
```

**実行順序:**
1. usersを全件読んでハッシュテーブル作成
2. todosを全件読む
3. Hash Joinで結合

#### レベル2: ボトルネックを特定

```
Seq Scan on todos  (actual time=0.012..8.123 rows=49876 loops=1)
  Filter: (done_flag = false)
  Rows Removed by Filter: 50124    ← ❌ 半分の行を捨てている
```

**問題:** done_flag=falseでフィルタリングしているが、インデックスがないため全件スキャン

#### レベル3: 実行時間の内訳

| 処理 | 時間 | 割合 |
|------|------|------|
| todos Seq Scan | 8.123 ms | 48% |
| Hash Join処理 | 7.111 ms | 42% |
| users Seq Scan | 0.089 ms | 1% |
| その他 | 1.466 ms | 9% |
| **合計** | **16.789 ms** | **100%** |

**ボトルネック:** todosのSeq Scanが最も遅い

### ステップ2: インデックスを作成

```sql
-- done_flagにインデックス
CREATE INDEX idx_todos_done_flag ON todos(done_flag);

-- さらに、user_idとの複合インデックスも検討
CREATE INDEX idx_todos_user_done ON todos(user_id, done_flag);

ANALYZE todos;
```

### ステップ3: 改善を確認

```sql
EXPLAIN ANALYZE
SELECT u.name, t.title
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.done_flag = false;
```

**出力:**
```
Nested Loop  (cost=0.43..1234.56 rows=50000 width=534)
             (actual time=0.023..3.456 rows=49876 loops=1)
  ->  Seq Scan on users u  (cost=0.00..10.50 rows=50 width=520)
                            (actual time=0.007..0.009 rows=4 loops=1)
  ->  Index Scan using idx_todos_user_done on todos t
          (cost=0.43..24.45 rows=12500 width=22)
          (actual time=0.004..0.123 rows=12469 loops=4)
        Index Cond: ((user_id = u.id) AND (done_flag = false))
Planning Time: 0.234 ms
Execution Time: 4.123 ms
```

### 📊 改善結果

| 項目 | Before | After | 改善 |
|-----|--------|-------|------|
| JOIN方法 | Hash Join | Nested Loop | ✅ |
| todos スキャン | Seq Scan | Index Scan | ✅ |
| 実行時間 | 16.789 ms | 4.123 ms | **75%削減** |

### 💡 読み方のポイント

1. **インデント（階層）に注目**
   - インデントが深い = 内側（先に実行される）
   - 下から上に実行されていく

2. **actual timeの合計を追う**
   - どの処理に時間がかかっているか特定
   - loops=1以外の場合は掛け算が必要

3. **JOIN方法の変化を確認**
   - Hash Join → Nested Loop に変わった理由は？
   - インデックスがあるとNested Loopが効率的になる

---

## 実践例3: 推定値と実測値の乖離

### シナリオ
「クエリプランナーが間違った判断をしている気がする」

### ステップ1: 問題のクエリ

```sql
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE title LIKE 'TODO 1%';
```

**出力:**
```
Seq Scan on todos  (cost=0.00..1768.00 rows=50 width=90)
                   (actual time=0.012..12.345 rows=11111 loops=1)
  Filter: (title ~~ 'TODO 1%'::text)
  Rows Removed by Filter: 88889
Planning Time: 0.123 ms
Execution Time: 12.789 ms
```

### 📊 分析: 推定と実測のギャップ

```
                推定値          実測値
rows=50        rows=11111    ← ❌ 222倍の差！
```

**問題:**
- PostgreSQLは50行しか返らないと推定
- 実際は11,111行返っている
- この誤差により、非効率な実行計画が選ばれる可能性

### ステップ2: 統計情報を確認

```sql
-- カラムの統計情報を確認
SELECT
    tablename,
    attname,
    n_distinct,
    correlation
FROM pg_stats
WHERE tablename = 'todos' AND attname = 'title';
```

**出力:**
```
 tablename | attname | n_distinct | correlation
-----------+---------+------------+-------------
 todos     | title   |         -1 |    0.123456
```

- `n_distinct = -1` → すべての値がユニークと推定
- `correlation = 0.123` → 物理的な順序との相関が低い

### ステップ3: 統計情報を更新

```sql
-- 統計情報の精度を上げる
ALTER TABLE todos ALTER COLUMN title SET STATISTICS 1000;

-- 統計情報を更新
ANALYZE todos;
```

### ステップ4: 改善を確認

```sql
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE title LIKE 'TODO 1%';
```

**出力:**
```
Seq Scan on todos  (cost=0.00..1768.00 rows=11000 width=90)
                   (actual time=0.012..12.234 rows=11111 loops=1)
  Filter: (title ~~ 'TODO 1%'::text)
  Rows Removed by Filter: 88889
```

### 📊 改善結果

| 項目 | Before | After |
|-----|--------|-------|
| 推定行数 | 50 | 11,000 |
| 実測行数 | 11,111 | 11,111 |
| 誤差 | 222倍 | 1.01倍 ✅ |

### 💡 学んだこと

1. **推定行数と実測行数の差が大きい場合**
   - まず`ANALYZE`を実行
   - それでも改善しない場合は`SET STATISTICS`で精度を上げる

2. **統計情報の重要性**
   - PostgreSQLは統計情報を基に実行計画を作成
   - 古い・不正確な統計情報は間違った判断を招く

---

## 実践例4: インデックスが使われない理由

### ケース1: テーブルが小さすぎる

```sql
CREATE INDEX idx_users_name ON users(name);

EXPLAIN ANALYZE
SELECT * FROM users WHERE name = '山田太郎';
```

**出力:**
```
Seq Scan on users  (cost=0.00..10.50 rows=1 width=1576)
                   (actual time=0.012..0.014 rows=1 loops=1)
  Filter: (name = '山田太郎'::text)
Planning Time: 0.234 ms
Execution Time: 0.034 ms
```

**なぜSeq Scanが選ばれた？**

```
テーブル: 4行（非常に小さい）
インデックススキャンのコスト: インデックス読み込み + テーブル読み込み
全件スキャンのコスト: テーブル1回読み込み

→ 全件スキャンの方が速い！
```

💡 **これは正しい判断** - 小さなテーブルではSeq Scanが効率的

### ケース2: 関数を使っている

```sql
CREATE INDEX idx_users_email ON users(email);

-- ❌ インデックスが使われない
EXPLAIN ANALYZE
SELECT * FROM users WHERE LOWER(email) = 'yamada@example.com';
```

**出力:**
```
Seq Scan on users  (cost=0.00..10.75 rows=1 width=1576)
  Filter: (lower(email) = 'yamada@example.com'::text)
```

**理由:** `LOWER(email)`に対するインデックスがない

**解決策:**

```sql
-- 関数インデックスを作成
CREATE INDEX idx_users_email_lower ON users(LOWER(email));

EXPLAIN ANALYZE
SELECT * FROM users WHERE LOWER(email) = 'yamada@example.com';
```

**出力:**
```
Index Scan using idx_users_email_lower on users
  Index Cond: (lower(email) = 'yamada@example.com'::text)
```

### ケース3: データの選択性が低い

```sql
CREATE INDEX idx_todos_done_flag ON todos(done_flag);

-- done_flag = false が全体の90%を占める場合
EXPLAIN ANALYZE
SELECT * FROM todos WHERE done_flag = false;
```

**出力:**
```
Seq Scan on todos  (cost=0.00..1768.00 rows=90000 width=90)
  Filter: (done_flag = false)
```

**理由:**
- 90%の行を取得する場合、全件スキャンの方が効率的
- インデックスを使うと、インデックス + テーブルの両方を読む必要がある

💡 **選択性が高い（少数の行を返す）場合にインデックスが有効**

---

## 実践例5: 複合インデックスの活用

### シナリオ
「複数の条件で検索するクエリを最適化したい」

### ステップ1: 現状（単一インデックス）

```sql
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_done_flag ON todos(done_flag);

EXPLAIN ANALYZE
SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false;
```

**出力:**
```
Bitmap Heap Scan on todos
    (cost=12.34..234.56 rows=6250 width=90)
    (actual time=0.234..1.234 rows=6234 loops=1)
  Recheck Cond: ((user_id = 1) AND (done_flag = false))
  ->  BitmapAnd
        ->  Bitmap Index Scan on idx_todos_user_id
        ->  Bitmap Index Scan on idx_todos_done_flag
```

### 📊 分析

**何が起きているか:**
1. user_idのインデックスをスキャン → ビットマップ作成
2. done_flagのインデックスをスキャン → ビットマップ作成
3. 2つのビットマップをAND演算
4. ビットマップを使ってテーブルから行を取得

**問題:**
- 2つのインデックスを両方読む必要がある
- Recheck（再確認）が必要

### ステップ2: 複合インデックスを作成

```sql
-- カラムの順序が重要！
CREATE INDEX idx_todos_user_done ON todos(user_id, done_flag);

EXPLAIN ANALYZE
SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false;
```

**出力:**
```
Index Scan using idx_todos_user_done on todos
    (cost=0.43..245.67 rows=6250 width=90)
    (actual time=0.034..0.456 rows=6234 loops=1)
  Index Cond: ((user_id = 1) AND (done_flag = false))
```

### 📊 改善結果

| 項目 | Before | After |
|-----|--------|-------|
| スキャン方法 | Bitmap Heap Scan | Index Scan |
| 使用インデックス | 2個 | 1個 |
| 実測時間 | 1.234 ms | 0.456 ms |
| 改善率 | - | **63%削減** |

### ステップ3: カラム順序の重要性を理解

```sql
-- パターンA: (user_id, done_flag)
CREATE INDEX idx_a ON todos(user_id, done_flag);

-- パターンB: (done_flag, user_id)
CREATE INDEX idx_b ON todos(done_flag, user_id);
```

**以下のクエリでどちらが使われる？**

```sql
-- クエリ1: 両方のカラムを指定
SELECT * FROM todos WHERE user_id = 1 AND done_flag = false;
→ idx_a ✅ / idx_b ✅ （どちらも使える）

-- クエリ2: user_idのみ
SELECT * FROM todos WHERE user_id = 1;
→ idx_a ✅ / idx_b ❌ （idx_aのみ使える）

-- クエリ3: done_flagのみ
SELECT * FROM todos WHERE done_flag = false;
→ idx_a ❌ / idx_b ✅ （idx_bのみ使える）
```

**ルール: 複合インデックスは左から順に使える**

```
idx_todos_user_done (user_id, done_flag)

使える:
  WHERE user_id = 1
  WHERE user_id = 1 AND done_flag = false

使えない:
  WHERE done_flag = false （先頭カラムがない）
```

### 💡 複合インデックスの設計指針

1. **選択性の高い（ユニークに近い）カラムを先頭に**
2. **よく使う検索パターンを優先**
3. **カラム数は3〜4個まで**（多すぎると効果が薄れる）

---

## まとめ: EXPLAINの活用フロー

### 📋 チェックリスト

実行計画を見たら、以下を確認:

#### 1. スキャン方法は適切か？
- [ ] Seq Scanが使われている場合、それは正当か？
- [ ] インデックスがあるのに使われていないか？
- [ ] 使われているインデックスは最適か？

#### 2. 推定値と実測値は近いか？
- [ ] `rows=推定` と `rows=実測` の差は10%以内か？
- [ ] 大きく乖離している場合、ANALYZEを実行したか？

#### 3. 無駄な処理はないか？
- [ ] `Rows Removed by Filter`が多すぎないか？
- [ ] 不要なソートは入っていないか？
- [ ] JOINの順序は適切か？

#### 4. 実行時間は許容範囲か？
- [ ] `Execution Time`は期待値以内か？
- [ ] どの処理に時間がかかっているか特定できたか？

### 🔄 改善サイクル

```
1. EXPLAIN ANALYZEで現状把握
   ↓
2. ボトルネックを特定
   ↓
3. 改善策を実施
   - インデックス作成
   - クエリ書き換え
   - 統計情報更新
   ↓
4. EXPLAIN ANALYZEで効果測定
   ↓
5. 改善が見られなければ別の手法を試す
```

### 📊 パフォーマンス目標

| クエリタイプ | 目標実行時間 |
|------------|------------|
| 主キー検索 | < 1 ms |
| インデックス検索 | < 10 ms |
| 小規模JOIN | < 50 ms |
| 集計クエリ | < 100 ms |

---

## 次のステップ

1. **[クイックスタート](./explain_quickstart.md)** で基本を復習
2. **[完全ガイド](./postgresql_explain_guide.md)** で詳細を学ぶ
3. **[実践演習SQL](../demo/explain_demo/examples.sql)** で手を動かす
4. 実際のプロジェクトで実践する

実行計画を読めるようになると、データベースのパフォーマンス問題の80%は解決できます！
