# Go Demo

Go言語の様々な概念やライブラリを学ぶためのデモプロジェクトです。

## 開発コマンド

```bash
# サーバー起動 (port 8080)
go run main.go

# テスト実行
go test ./...

# データベースリフレッシュ
make refresh-schema
```

## デモパッケージ

各パッケージは独立した `go.mod` を持ち、個別に実行できます。

| カテゴリ | パッケージ | 説明 |
|---------|-----------|------|
| **アルゴリズム** | [algorithm](demo/algorithm/) | ソート・探索アルゴリズム |
| **データベース** | [explain_demo](demo/explain_demo/) | PostgreSQL EXPLAIN解析 |
| | [db_index_performance](demo/db_index_performance/) | インデックス性能検証 |
| | [n_plus_1_problem](demo/n_plus_1_problem/) | N+1問題と解決策 |
| **設計パターン** | [designpattern](demo/designpattern/) | GoFデザインパターン |
| | [ddd](demo/ddd/) | ドメイン駆動設計 |
| | [di_demo](demo/di_demo/) | 依存性注入 |
| **セキュリティ** | [crypto_demo](demo/crypto_demo/) | 暗号化 (Caesar, RSA等) |
| | [security](demo/security/) | セキュリティ対策 |
| | [jwt_demo](demo/jwt_demo/) | JWT認証 |
| | [oauth](demo/oauth/) | OAuth実装 |
| | [casbin_demo](demo/casbin_demo/) | アクセス制御 |
| **ネットワーク** | [net_http](demo/net_http/) | net/httpパッケージ |
| | [graphql](demo/graphql/) | GraphQLサーバー |
| | [scraping](demo/scraping/) | Webスクレイピング |
| **並行処理** | [goroutine_channel_basics](demo/goroutine_channel_basics/) | Goroutine・Channel基礎 |
| | [performance](demo/performance/) | パフォーマンス計測 |
| **その他** | [zap](demo/zap/) | 構造化ログ |
| | [rabbitmq](demo/rabbitmq/) | メッセージキュー |
| | [error_handling](demo/error_handling/) | エラーハンドリング |
| | [go101](demo/go101/) | Go言語基礎 |

## 学習教材

### PostgreSQL EXPLAIN完全ガイド

- [クイックスタート](docs/explain_quickstart.md) - 5分で始める
- [実践的な読み方ガイド](docs/explain_practical_examples.md) - 出力の読み方と活用方法
- [完全ガイド](docs/postgresql_explain_guide.md) - 詳細解説
- [総合案内](docs/README_EXPLAIN.md) - すべての教材の案内

```bash
# すぐに試す
cd demo/explain_demo
./quick_demo.sh
```

## 主要な依存ライブラリ

- **GORM** - ORM (PostgreSQL)
- **gqlgen** - GraphQL
- **Cobra** - CLI
- **testify** - テスト
