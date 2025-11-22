#!/bin/bash

# シンプルなEXPLAINデモスクリプト

echo "=========================================="
echo "PostgreSQL EXPLAIN クイックデモ"
echo "=========================================="
echo ""

# データベース接続確認
if ! docker exec go-demo-db pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ PostgreSQLコンテナが起動していません"
    echo "   docker compose up -d db"
    exit 1
fi

echo "✅ データベース接続OK"
echo ""

# デモ実行
docker exec go-demo-db psql -U postgres -d go_demo << 'EOF'
\echo '■ 例1: 基本的なSELECT'
\echo '---'
EXPLAIN ANALYZE SELECT * FROM users;

\echo ''
\echo '■ 例2: WHEREでのフィルタリング'
\echo '---'
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 1 LIMIT 5;

\echo ''
\echo '■ 例3: JOIN操作'
\echo '---'
EXPLAIN ANALYZE
SELECT t.title, u.name
FROM todos t
JOIN users u ON t.user_id = u.id
LIMIT 5;

\echo ''
\echo '■ 例4: 集約クエリ'
\echo '---'
EXPLAIN ANALYZE
SELECT user_id, COUNT(*) as count
FROM todos
GROUP BY user_id;
EOF

echo ""
echo "=========================================="
echo "📚 次のステップ"
echo "=========================================="
echo ""
echo "1. PostgreSQLに接続して対話的に実行:"
echo "   docker exec -it go-demo-db psql -U postgres -d go_demo"
echo ""
echo "2. クイックスタートを読む:"
echo "   cat ../../docs/explain_quickstart.md"
echo ""
echo "3. 実践演習を試す:"
echo "   cat examples.sql"
echo ""
