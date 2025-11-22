-- ============================================================
-- PostgreSQL EXPLAIN 実践演習用SQLスクリプト
-- ============================================================
--
-- このファイルには、EXPLAINの使い方を学ぶための実践的な例が含まれています。
--
-- 実行方法:
-- 1. PostgreSQLコンテナに接続:
--    docker exec -it go-demo-db psql -U postgres -d go_demo
--
-- 2. このファイルを読み込む（オプション）:
--    \i /docker-entrypoint-initdb.d/examples.sql
--
-- 3. または、このファイルから個別にコピー＆ペーストして実行
-- ============================================================


-- ============================================================
-- 準備: 現在のデータ状態を確認
-- ============================================================

\echo '■ データベースの状態確認'
\echo ''

-- テーブル一覧
\dt

-- レコード数確認
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'todos', COUNT(*) FROM todos;

-- インデックス一覧
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;


-- ============================================================
-- 例1: 基本的なEXPLAIN - 実行計画のみ
-- ============================================================

\echo ''
\echo '■ 例1: 基本的なEXPLAIN'
\echo '説明: クエリは実行せず、実行計画（プラン）のみを表示します'
\echo ''

-- 全件取得の実行計画
EXPLAIN SELECT * FROM users;

-- 解説:
-- Seq Scan = Sequential Scan（順次スキャン）
-- cost=0.00..10.50 → 起動コスト..総コスト
-- rows=50 → 推定される行数
-- width=1576 → 各行の推定バイト数


-- ============================================================
-- 例2: EXPLAIN ANALYZE - 実行計画 + 実測値
-- ============================================================

\echo ''
\echo '■ 例2: EXPLAIN ANALYZE（実際に実行）'
\echo '説明: クエリを実際に実行して、推定値と実測値を比較します'
\echo ''

EXPLAIN ANALYZE SELECT * FROM users;

-- 解説:
-- actual time=0.007..0.009 → 実際にかかった時間（ミリ秒）
-- rows=4 → 実際に返された行数
-- loops=1 → このノードの実行回数
-- Planning Time: クエリプランの作成時間
-- Execution Time: 実際の実行時間


-- ============================================================
-- 例3: WHERE句でのフィルタリング
-- ============================================================

\echo ''
\echo '■ 例3: WHERE句でのフィルタリング'
\echo '説明: 条件によってどのように実行計画が変わるか確認します'
\echo ''

-- 特定ユーザーのTODO検索
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 1;

-- 【出力の読み方】
-- Seq Scan on todos (cost=0.00..268.00 rows=10000 width=90)
--                   (actual time=0.009..1.309 rows=10000 loops=1)
--   Filter: (user_id = 1)
-- Planning Time: 2.441 ms
-- Execution Time: 1.656 ms
--
-- 📊 分析:
-- ・Seq Scan = 全件スキャン（改善の余地あり）
-- ・cost=0.00..268.00 = 起動コスト0、総コスト268
-- ・rows=10000 = 10,000行返すと推定
-- ・actual time=0.009..1.309 = 実際に1.3msかかった
-- ・Filter: (user_id = 1) = WHERE句で絞り込み
-- ・Execution Time: 1.656 ms = トータル実行時間
--
-- ❌ 問題点: 全行を読んでからフィルタリング（非効率）
-- ✅ 改善策: user_idにインデックスを作成


-- ============================================================
-- 例4: インデックスの効果
-- ============================================================

\echo ''
\echo '■ 例4: インデックスの効果を確認'
\echo ''

-- インデックスが存在するか確認
SELECT 'idx_todos_user_id' as index_name,
       EXISTS (
           SELECT 1 FROM pg_indexes
           WHERE tablename = 'todos' AND indexname = 'idx_todos_user_id'
       ) as exists;

-- インデックスが存在しない場合は作成
-- CREATE INDEX idx_todos_user_id ON todos(user_id);

-- インデックス作成後に再度実行
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 2;

-- 【出力の読み方 - インデックスあり】
-- Index Scan using idx_todos_user_id on todos
--     (cost=0.29..856.42 rows=25000 width=90)
--     (actual time=0.034..0.523 rows=25123 loops=1)
--   Index Cond: (user_id = 2)
-- Execution Time: 0.876 ms
--
-- 📊 分析:
-- ・Index Scan = インデックスを使用（効率的！）
-- ・Index Cond: (user_id = 2) = インデックスの条件
-- ・actual time=0.034..0.523 = 高速化された
--
-- 💡 比較:
-- インデックスなし: Execution Time: 4.123 ms
-- インデックスあり: Execution Time: 0.876 ms
-- → 約5倍高速化！
--
-- 注意: データの分布によって選択が変わります
-- ・少数の行を取得: Index Scan（効率的）
-- ・大部分の行を取得: Seq Scan（全件読む方が速い）


-- ============================================================
-- 例5: 大量データを追加してパフォーマンスを体感
-- ============================================================

\echo ''
\echo '■ 例5: 大量データでの実験'
\echo '説明: データ量を増やしてインデックスの効果を体感します'
\echo ''

-- 現在のレコード数
SELECT COUNT(*) as current_count FROM todos;

-- 大量データを追加（10万件）
-- 注意: すでにデータがある場合はスキップしてください
/*
INSERT INTO todos (user_id, title, note, done_flag)
SELECT
    (random() * 3 + 1)::int,  -- user_id 1-4
    'TODO ' || generate_series,
    'Note ' || generate_series,
    random() > 0.5
FROM generate_series(1, 100000);
*/

-- 統計情報を更新
ANALYZE todos;

-- インデックスなしで特定ユーザーを検索（遅い）
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 3;

-- インデックス作成
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);

-- インデックスありで同じクエリを実行（速い）
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = 3;

-- 実行時間を比較してください！


-- ============================================================
-- 例6: JOIN操作の実行計画
-- ============================================================

\echo ''
\echo '■ 例6: JOIN操作'
\echo '説明: テーブル結合時の実行計画を確認します'
\echo ''

-- 基本的なJOIN
EXPLAIN ANALYZE
SELECT t.title, u.name
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.done_flag = false
LIMIT 10;

-- 【出力の読み方 - JOIN】
-- Limit (cost=0.15..0.64 rows=10 width=534)
--   -> Nested Loop (cost=0.15..490.44 rows=10000 width=534)
--        -> Seq Scan on todos t (cost=0.00..243.00 rows=10000 width=22)
--             Filter: (done_flag = false)
--        -> Index Scan using users_pkey on users u
--             Index Cond: (id = t.user_id)
--
-- 📊 読み方（下から上に実行される）:
-- 1. todosテーブルを全件スキャン（Filter: done_flag = false）
-- 2. 各行について、usersテーブルをインデックスで検索
-- 3. Nested Loopで結合
-- 4. Limitで10件に絞る
--
-- 💡 JOIN方法の種類:
-- ・Nested Loop: 小規模データ、インデックスありで高速
-- ・Hash Join: 大規模データの結合に適している
-- ・Merge Join: ソート済みデータの結合に最適


-- ============================================================
-- 例7: 集約クエリ（GROUP BY）
-- ============================================================

\echo ''
\echo '■ 例7: 集約クエリ'
\echo '説明: GROUP BYの実行計画を確認します'
\echo ''

-- ユーザーごとのTODO件数
EXPLAIN ANALYZE
SELECT user_id, COUNT(*) as todo_count
FROM todos
GROUP BY user_id;

-- 解説:
-- HashAggregate: ハッシュテーブルを使った集約（高速）
-- GroupAggregate: ソート済みデータのグループ化


-- ユーザーごとのTODO件数（ソート付き）
EXPLAIN ANALYZE
SELECT user_id, COUNT(*) as todo_count
FROM todos
GROUP BY user_id
ORDER BY todo_count DESC;


-- ============================================================
-- 例8: 複合インデックス
-- ============================================================

\echo ''
\echo '■ 例8: 複合インデックスの活用'
\echo '説明: 複数カラムのインデックスの効果を確認します'
\echo ''

-- 複合インデックス作成
CREATE INDEX IF NOT EXISTS idx_todos_user_done ON todos(user_id, done_flag);

-- 統計情報更新
ANALYZE todos;

-- 複合条件での検索
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false
LIMIT 10;

-- 解説:
-- Index Cond: 両方の条件がインデックス条件として使われているか確認
-- Filter: 追加のフィルタリングが必要か確認


-- Index Only Scanの例（カラムを絞る）
EXPLAIN ANALYZE
SELECT user_id, done_flag
FROM todos
WHERE user_id = 2
LIMIT 10;

-- 解説:
-- Index Only Scan: テーブルにアクセスせずインデックスだけで結果を返す
-- 最も高速な方法！


-- ============================================================
-- 例9: サブクエリの実行計画
-- ============================================================

\echo ''
\echo '■ 例9: サブクエリ'
\echo '説明: サブクエリがどのように実行されるか確認します'
\echo ''

-- TODOが複数あるユーザーを検索
EXPLAIN ANALYZE
SELECT name FROM users
WHERE id IN (
    SELECT user_id
    FROM todos
    GROUP BY user_id
    HAVING COUNT(*) >= 3
);

-- EXISTS句での書き換え（場合によっては高速）
EXPLAIN ANALYZE
SELECT name FROM users u
WHERE EXISTS (
    SELECT 1
    FROM todos t
    WHERE t.user_id = u.id
    GROUP BY t.user_id
    HAVING COUNT(*) >= 3
);


-- ============================================================
-- 例10: 詳細オプション
-- ============================================================

\echo ''
\echo '■ 例10: 詳細オプション'
\echo '説明: より詳細な情報を表示します'
\echo ''

-- BUFFERS: バッファの使用状況も表示
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM todos WHERE user_id = 1 LIMIT 10;

-- 解説:
-- Shared Hit: 共有バッファからの読み取り（キャッシュヒット）
-- Shared Read: ディスクからの読み取り
-- バッファヒット率が高いほど高速


-- VERBOSE: より詳細な情報
EXPLAIN (ANALYZE, VERBOSE)
SELECT t.title, u.name
FROM todos t
JOIN users u ON t.user_id = u.id
LIMIT 5;


-- JSON形式で出力（ツールで分析しやすい）
EXPLAIN (ANALYZE, FORMAT JSON)
SELECT * FROM users;


-- ============================================================
-- 例11: パフォーマンスチューニングのヒント
-- ============================================================

\echo ''
\echo '■ 例11: パフォーマンスチューニング'
\echo ''

-- 現在の設定を確認
SHOW work_mem;
SHOW random_page_cost;

-- work_memを一時的に増やす（大きなソートやハッシュに効果的）
SET work_mem = '64MB';

-- 大きなソートを伴うクエリ
EXPLAIN ANALYZE
SELECT * FROM todos ORDER BY created_at DESC LIMIT 100;

-- 設定を戻す
RESET work_mem;


-- ============================================================
-- 例12: 統計情報の確認
-- ============================================================

\echo ''
\echo '■ 例12: 統計情報'
\echo '説明: PostgreSQLが保持している統計情報を確認します'
\echo ''

-- テーブルの統計情報
SELECT
    schemaname,
    tablename,
    n_live_tup as live_rows,
    n_dead_tup as dead_rows,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables
WHERE tablename IN ('users', 'todos');

-- カラムの統計情報
SELECT
    tablename,
    attname as column_name,
    n_distinct,
    most_common_vals,
    most_common_freqs
FROM pg_stats
WHERE tablename = 'todos' AND attname = 'user_id';


-- ============================================================
-- 例13: インデックスの使用状況
-- ============================================================

\echo ''
\echo '■ 例13: インデックスの使用状況'
\echo '説明: どのインデックスが使われているか確認します'
\echo ''

-- インデックスのスキャン回数
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

-- 使われていないインデックスを発見
SELECT
    schemaname || '.' || tablename AS table,
    indexname AS index,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size,
    idx_scan AS scans
FROM pg_stat_user_indexes
WHERE idx_scan = 0
    AND indexrelname NOT LIKE '%pkey'
ORDER BY pg_relation_size(indexrelid) DESC;


-- ============================================================
-- 例14: テーブルサイズの確認
-- ============================================================

\echo ''
\echo '■ 例14: テーブルとインデックスのサイズ'
\echo ''

-- テーブルサイズ
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;


-- ============================================================
-- 練習問題
-- ============================================================

\echo ''
\echo '■ 練習問題'
\echo '以下のクエリの実行計画を分析して、最適化してください'
\echo ''

-- 問題1: このクエリを高速化するには？
-- ヒント: どのカラムにインデックスを張るべきか考えてください
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE done_flag = false
ORDER BY created_at DESC
LIMIT 20;

-- 問題2: このJOINを最適化するには？
EXPLAIN ANALYZE
SELECT u.name, COUNT(t.id) as todo_count
FROM users u
LEFT JOIN todos t ON u.id = t.user_id
WHERE t.done_flag = false
GROUP BY u.id, u.name;

-- 問題3: この重いクエリを改善するには？
EXPLAIN ANALYZE
SELECT *
FROM todos
WHERE title LIKE '%TODO%'
AND done_flag = false;


-- ============================================================
-- クリーンアップ（必要に応じて実行）
-- ============================================================

/*
-- 作成したインデックスを削除
DROP INDEX IF EXISTS idx_todos_user_id;
DROP INDEX IF EXISTS idx_todos_user_done;

-- 追加したデータを削除
DELETE FROM todos WHERE id > 10;
*/


-- ============================================================
-- まとめ
-- ============================================================

\echo ''
\echo '■ まとめ'
\echo ''
\echo 'EXPLAINを使いこなすポイント:'
\echo '1. EXPLAIN ANALYZEで実測値を確認する'
\echo '2. 推定行数と実測行数の差に注目する'
\echo '3. Seq Scanが適切か判断する（小規模テーブルではOK）'
\echo '4. インデックスは万能ではない（メンテナンスコストも考慮）'
\echo '5. 定期的にANALYZEを実行して統計情報を更新する'
\echo ''
\echo '次のステップ:'
\echo '- 実際のプロジェクトで遅いクエリを見つける'
\echo '- EXPLAINで原因を特定する'
\echo '- インデックスやクエリの書き換えで改善する'
\echo '- 効果を測定する'
\echo ''
