#!/bin/bash

# PostgreSQL EXPLAIN デモを実行するスクリプト

set -e

echo "=================================================="
echo "PostgreSQL EXPLAIN デモプログラム"
echo "=================================================="
echo ""

# データベース接続確認
echo "🔍 データベース接続確認中..."
if ! docker exec go-demo-db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ PostgreSQLコンテナが起動していません"
    echo ""
    echo "以下のコマンドで起動してください:"
    echo "  docker compose up -d db"
    exit 1
fi

echo "✅ データベース接続OK"
echo ""

# レコード数確認
echo "📊 データ確認中..."
USERS_COUNT=$(docker exec go-demo-db psql -U postgres -d go_demo -t -c "SELECT COUNT(*) FROM users;")
TODOS_COUNT=$(docker exec go-demo-db psql -U postgres -d go_demo -t -c "SELECT COUNT(*) FROM todos;")

echo "  - users: ${USERS_COUNT} 件"
echo "  - todos: ${TODOS_COUNT} 件"
echo ""

if [ "${TODOS_COUNT// /}" -lt 1000 ]; then
    echo "💡 ヒント: より効果的な学習のため、大量データを投入することをお勧めします"
    echo "  以下のコマンドで10万件のデータを追加できます:"
    echo ""
    echo "  docker exec go-demo-db psql -U postgres -d go_demo -c \""
    echo "    INSERT INTO todos (user_id, title, note, done_flag)"
    echo "    SELECT"
    echo "      (random() * 3 + 1)::int,"
    echo "      'TODO ' || generate_series,"
    echo "      'Note ' || generate_series,"
    echo "      random() > 0.5"
    echo "    FROM generate_series(1, 100000);"
    echo "  \""
    echo ""
    read -p "データを追加しますか？ (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "📝 データ追加中..."
        docker exec go-demo-db psql -U postgres -d go_demo -c "
            INSERT INTO todos (user_id, title, note, done_flag)
            SELECT
              (random() * 3 + 1)::int,
              'TODO ' || generate_series,
              'Note ' || generate_series,
              random() > 0.5
            FROM generate_series(1, 100000);
        " > /dev/null

        echo "✅ データ追加完了"

        # 統計情報更新
        echo "📊 統計情報更新中..."
        docker exec go-demo-db psql -U postgres -d go_demo -c "ANALYZE todos;" > /dev/null
        echo "✅ 統計情報更新完了"
        echo ""
    fi
fi

echo "=================================================="
echo "実行方法を選択してください:"
echo "=================================================="
echo ""
echo "1) PostgreSQLに接続して対話的に実行（推奨）"
echo "2) 基本的なEXPLAIN例を自動実行"
echo "3) Goデモプログラムを実行"
echo "0) 終了"
echo ""
read -p "選択 (1-3): " choice

case $choice in
    1)
        echo ""
        echo "PostgreSQLに接続します..."
        echo ""
        echo "💡 ヒント:"
        echo "  - examples.sqlの内容をコピー＆ペーストして実行できます"
        echo "  - または、以下のような基本コマンドを試してください:"
        echo "    EXPLAIN ANALYZE SELECT * FROM users;"
        echo "    \\q で終了"
        echo ""
        docker exec -it go-demo-db psql -U postgres -d go_demo
        ;;
    2)
        echo ""
        echo "基本的なEXPLAIN例を実行します..."
        echo ""

        docker exec go-demo-db psql -U postgres -d go_demo << 'EOF'
\echo '==================================================================='
\echo '例1: 基本的なSELECT'
\echo '==================================================================='
EXPLAIN ANALYZE SELECT * FROM users;

\echo ''
\echo '==================================================================='
\echo '例2: WHEREでのフィルタリング'
\echo '==================================================================='
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 1 LIMIT 10;

\echo ''
\echo '==================================================================='
\echo '例3: JOIN操作'
\echo '==================================================================='
EXPLAIN ANALYZE
SELECT t.title, u.name
FROM todos t
JOIN users u ON t.user_id = u.id
LIMIT 10;

\echo ''
\echo '==================================================================='
\echo '例4: 集約クエリ'
\echo '==================================================================='
EXPLAIN ANALYZE
SELECT user_id, COUNT(*) as count
FROM todos
GROUP BY user_id;

\echo ''
\echo '==================================================================='
\echo '例5: インデックスの効果確認'
\echo '==================================================================='

-- インデックスが存在するか確認
\echo 'インデックス確認:'
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'todos' AND indexname = 'idx_todos_user_id';

-- インデックスがない場合は作成を提案
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE tablename = 'todos' AND indexname = 'idx_todos_user_id'
    ) THEN
        RAISE NOTICE 'インデックスが存在しません。作成するには以下を実行:';
        RAISE NOTICE 'CREATE INDEX idx_todos_user_id ON todos(user_id);';
    END IF;
END $$;

\echo ''
\echo '==================================================================='
\echo '完了！'
\echo '==================================================================='
\echo ''
\echo '次のステップ:'
\echo '  1. PostgreSQLに接続: docker exec -it go-demo-db psql -U postgres -d go_demo'
\echo '  2. examples.sqlを試す'
\echo '  3. docs/explain_quickstart.md を読む'
EOF
        ;;
    3)
        echo ""
        echo "Goデモプログラムを実行します..."
        echo ""

        # Dockerネットワークを使用して実行
        docker run --rm \
          --network go-demo_go-demo-network \
          -v "$(pwd)":/app \
          -w /app \
          -e DB_HOST=db \
          golang:1.23 \
          sh -c "go mod download > /dev/null 2>&1 && go run main.go"
        ;;
    0)
        echo "終了します"
        exit 0
        ;;
    *)
        echo "無効な選択です"
        exit 1
        ;;
esac

echo ""
echo "=================================================="
echo "お疲れさまでした！"
echo "=================================================="
echo ""
echo "📚 さらに学ぶには:"
echo "  - クイックスタート: docs/explain_quickstart.md"
echo "  - 完全ガイド: docs/postgresql_explain_guide.md"
echo "  - 実践演習: demo/explain_demo/examples.sql"
echo ""
