## 概要

- GoとRedisを使ったキャッシュのパフォーマンス比較デモ
- キャッシュの有無による応答速度の違いを実測し、実際のアプリケーションで使われるキャッシング戦略を学ぶ

### 特徴

- **実用的なユースケース**: 管理画面のログビューアを想定した実装
- **明確なパフォーマンス差**: キャッシュあり（0-2ms）vs キャッシュなし（50-100ms）
- **現実的なデータ量**: 10万件のログデータから最新100件を取得
- **RESTful API**: エンドポイントで簡単にテスト可能
- **Docker Compose対応**: Redis環境を簡単に構築

## キャッシュとは

キャッシュは頻繁にアクセスされるデータを高速なストレージ（メモリ）に一時保存し、データ取得を高速化する技術。

### キャッシュのメリット

1. **レスポンス速度の向上**: データベースアクセスを削減
2. **データベース負荷の軽減**: 同じクエリの繰り返し実行を防ぐ
3. **スケーラビリティ向上**: システム全体のスループット向上
4. **コスト削減**: データベースのリソース使用量を削減

### Redisの特徴

- **インメモリデータストア**: 全データをメモリに保持し超高速アクセス
- **豊富なデータ構造**: String、Hash、List、Set、Sorted Setなど
- **TTL（有効期限）サポート**: 自動的に古いデータを削除
- **永続化オプション**: メモリデータをディスクに保存可能

## 実装パターン

### 1. Cache-Aside（Lazy Loading）パターン

このデモで採用している最も一般的なパターン：

```
1. アプリケーションがキャッシュを確認
2. キャッシュヒット → キャッシュからデータを返す
3. キャッシュミス → データベースから取得してキャッシュに保存
```

#### メリット
- 実装がシンプル
- 必要なデータだけがキャッシュされる
- キャッシュ障害時もデータベースで動作継続

#### デメリット
- 初回アクセスは遅い（キャッシュミス）
- データの整合性管理が必要

### 2. その他の主要パターン

#### Write-Through
- データ書き込み時に同時にキャッシュも更新
- データの一貫性が高い
- 書き込みが遅くなる

#### Write-Behind
- データをキャッシュに書き込み、非同期でDBに反映
- 書き込みが高速
- データロストのリスク

## コード解説

### 定数定義

```go
const (
    cacheKey       = "logs:latest"
    cacheTTL       = 5 * time.Minute
    testDataCount  = 100000  // テストデータ10万件
    fetchLimit     = 100     // 取得件数（実用的なサイズ）
)
```

- `cacheKey`: Redisのキー名
- `cacheTTL`: キャッシュの有効期限（5分）
- `testDataCount`: テスト用に作成するログレコード数
- `fetchLimit`: 一度に取得するログ件数

### データ構造

```go
type LogRecord struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Message   string    `json:"message"`
    Level     string    `json:"level"`
    Timestamp time.Time `json:"timestamp"`
}

type PerformanceResult struct {
    Source     string      `json:"source"`
    Count      int         `json:"count"`
    DurationMs int64       `json:"duration_ms"`
    Records    []LogRecord `json:"records,omitempty"`
}
```

### キャッシュありの実装

```go
func (h *CacheHandler) CacheWithHandler(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    start := time.Now()

    // Redisからキャッシュを取得
    cached, err := h.Redis.Get(ctx, cacheKey).Result()
    if err == nil {
        // キャッシュヒット
        var logs []LogRecord
        if err := json.Unmarshal([]byte(cached), &logs); err == nil {
            duration := time.Since(start)
            result := PerformanceResult{
                Source:     "cache (Redis)",
                Count:      len(logs),
                DurationMs: duration.Milliseconds(),
                Records:    logs,
            }
            writeJSONResponse(w, result)
            return
        }
    }

    // キャッシュミス - DBから最新100件を取得
    var logs []LogRecord
    if err := h.DB.Order("timestamp DESC").Limit(fetchLimit).Find(&logs).Error; err != nil {
        http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
        return
    }

    // Redisにキャッシュを保存
    data, _ := json.Marshal(logs)
    h.Redis.Set(ctx, cacheKey, data, cacheTTL)

    duration := time.Since(start)
    result := PerformanceResult{
        Source:     "database (PostgreSQL) - cached for next request",
        Count:      len(logs),
        DurationMs: duration.Milliseconds(),
        Records:    logs,
    }
    writeJSONResponse(w, result)
}
```

#### 処理フロー

1. **キャッシュ確認**: `Redis.Get()`でキャッシュをチェック
2. **キャッシュヒット**: JSONをデシリアライズして即座に返す
3. **キャッシュミス**: PostgreSQLから最新100件取得
4. **キャッシュ保存**: 取得したデータをJSON化してRedisに保存
5. **レスポンス**: データとパフォーマンス情報を返す

### キャッシュなしの実装

```go
func (h *CacheHandler) CacheWithoutHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // DBから最新100件を取得（キャッシュなし）
    var logs []LogRecord
    if err := h.DB.Order("timestamp DESC").Limit(fetchLimit).Find(&logs).Error; err != nil {
        http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
        return
    }

    duration := time.Since(start)
    result := PerformanceResult{
        Source:     "database (PostgreSQL) - no cache",
        Count:      len(logs),
        DurationMs: duration.Milliseconds(),
        Records:    logs,
    }
    writeJSONResponse(w, result)
}
```

毎回データベースに問い合わせるため、キャッシュありと比較して遅い。

## 環境構築

### Docker Composeの設定

```yaml
redis:
  container_name: go-demo-redis
  image: redis:7-alpine
  ports:
    - 6379:6379
  volumes:
    - ./persist/redis:/data
  networks:
    - go-demo-network
```

### Redis接続

```go
// main.go
rdb := redis.NewClient(&redis.Options{
    Addr:     "redis:6379",
    Password: "",
    DB:       0,
})

// Redisの接続確認
ctx := context.Background()
if err := rdb.Ping(ctx).Err(); err != nil {
    log.Printf("Warning: Redis connection failed: %v", err)
}
```

### 依存パッケージ

```bash
go get github.com/redis/go-redis/v9
```

## API エンドポイント

### 1. テストデータ作成
```bash
GET /demo/cache/init
```

10万件のログレコードを作成

**レスポンス例:**
```json
{
  "message": "Test data created successfully",
  "record_count": 100000,
  "duration_ms": 5432
}
```

### 2. キャッシュあり
```bash
GET /demo/cache/with
```

**初回レスポンス（キャッシュミス）:**
```json
{
  "source": "database (PostgreSQL) - cached for next request",
  "count": 100,
  "duration_ms": 52,
  "records": [...]
}
```

**2回目以降（キャッシュヒット）:**
```json
{
  "source": "cache (Redis)",
  "count": 100,
  "duration_ms": 0,
  "records": [...]
}
```

### 3. キャッシュなし
```bash
GET /demo/cache/without
```

**レスポンス例:**
```json
{
  "source": "database (PostgreSQL) - no cache",
  "count": 100,
  "duration_ms": 48,
  "records": [...]
}
```

### 4. キャッシュクリア
```bash
GET /demo/cache/clear
```

**レスポンス例:**
```json
{
  "message": "Cache cleared successfully"
}
```

## パフォーマンス測定結果

### 実測値

| 項目 | レスポンス時間 | 備考 |
|-----|--------------|------|
| キャッシュあり（Redis） | 0-2ms | インメモリアクセス |
| キャッシュあり（初回） | 50-70ms | DB取得+キャッシュ保存 |
| キャッシュなし（PostgreSQL） | 50-100ms | 毎回DB検索 |

### 速度向上率

**キャッシュあり（2回目以降）は約50倍高速**

```
50ms（DB） → 1ms（Redis）
```

### 実際の本番環境での傾向

```
Twitter: タイムライン取得
- キャッシュあり: ~5ms
- キャッシュなし: ~100ms

Facebook: ニュースフィード
- キャッシュあり: ~10ms
- キャッシュなし: ~200ms

ECサイト: 商品一覧
- キャッシュあり: ~2ms
- キャッシュなし: ~80ms
```

## 実行方法

### 1. Redisコンテナを起動

```bash
docker-compose up -d redis
```

### 2. アプリケーション起動

```bash
go run main.go
```

### 3. テストデータ作成

```bash
curl http://localhost:8080/demo/cache/init
```

### 4. パフォーマンス比較

```bash
# キャッシュクリア
curl http://localhost:8080/demo/cache/clear

# キャッシュなしで測定
curl http://localhost:8080/demo/cache/without | jq '.duration_ms'

# キャッシュありで測定（初回）
curl http://localhost:8080/demo/cache/with | jq '.duration_ms'

# キャッシュありで測定（2回目 - 高速）
curl http://localhost:8080/demo/cache/with | jq '.duration_ms'
```

## 実用的なユースケース

### 1. 管理画面のログビューア

```go
// 最新100件のログを表示
// 頻繁にアクセスされるがデータはそこまで頻繁に変わらない
// → キャッシュに最適
```

### 2. ダッシュボードの統計情報

```go
// 今日のアクセス数、売上など
// 1分ごとに更新で十分
// → TTL 1分でキャッシュ
```

### 3. APIレート制限

```go
// ユーザーごとのリクエスト回数をRedisで管理
// INCR + EXPIRE でシンプルに実装可能
```

### 4. セッション管理

```go
// ログイン状態をRedisに保存
// 複数サーバーで共有可能
```

## キャッシュ戦略のベストプラクティス

### 1. 適切なTTL設定

```go
// データの性質に応じて設定
- ほぼ静的なデータ: 1時間〜1日
- 頻繁に更新されるデータ: 1分〜5分
- リアルタイム性が重要: 10秒〜30秒
```

### 2. キャッシュするデータの選択

**キャッシュに適したデータ:**
- 読み取り頻度が高い
- 更新頻度が低い
- 計算コストが高い
- サイズが小さい〜中程度

**キャッシュに不適切なデータ:**
- リアルタイム性が最重要
- 頻繁に更新される
- サイズが巨大（数MB以上）
- ユーザーごとに異なる（カスタマイズ性が高い）

### 3. キャッシュキーの命名規則

```go
// 明確で衝突しない命名
const (
    userProfileKey = "user:profile:%d"        // user:profile:123
    productListKey = "products:list:page:%d"  // products:list:page:1
    statsKey = "stats:daily:%s"               // stats:daily:2024-01-15
)
```

### 4. キャッシュの無効化戦略

```go
// データ更新時にキャッシュをクリア
func UpdateUser(userID int, data UserData) error {
    // DBを更新
    if err := db.Update(userID, data); err != nil {
        return err
    }

    // キャッシュを削除
    cacheKey := fmt.Sprintf("user:profile:%d", userID)
    redis.Del(ctx, cacheKey)

    return nil
}
```

### 5. エラーハンドリング

```go
// Redisが落ちてもアプリケーションは動作継続
cached, err := redis.Get(ctx, key).Result()
if err != nil {
    // キャッシュミスとして扱い、DBから取得
    log.Printf("Cache miss or error: %v", err)
    return fetchFromDB()
}
```

## よくある問題と対策

### 1. Thundering Herd（キャッシュスタンピード）

**問題:** キャッシュ期限切れ時に大量リクエストが同時にDBにアクセス

**対策:**
```go
// ロックを使って1つのリクエストだけがDBアクセス
func GetWithLock(key string) (data, error) {
    // キャッシュ確認
    if cached := redis.Get(key); cached != nil {
        return cached, nil
    }

    // 分散ロック取得
    lock := redis.SetNX(key+":lock", "1", 10*time.Second)
    if !lock {
        // 他のリクエストがDBアクセス中 → 少し待ってリトライ
        time.Sleep(100 * time.Millisecond)
        return GetWithLock(key)
    }

    // DBから取得してキャッシュ
    data := fetchFromDB()
    redis.Set(key, data, cacheTTL)
    redis.Del(key + ":lock")

    return data, nil
}
```

### 2. キャッシュとDBの不整合

**問題:** データ更新時にキャッシュ削除を忘れる

**対策:**
- トランザクション内でキャッシュ削除
- イベント駆動でキャッシュ無効化
- TTLを短めに設定

### 3. メモリ不足

**問題:** Redisのメモリが枯渇

**対策:**
```bash
# redis.conf で最大メモリと削除ポリシー設定
maxmemory 2gb
maxmemory-policy allkeys-lru  # LRUで古いキーを削除
```

## まとめ

### 学んだこと

1. **Redisキャッシュの効果**: 50倍以上の速度向上を実現
2. **Cache-Asideパターン**: 最も一般的で実用的なキャッシング戦略
3. **適切なデータサイズ**: 100万件ではなく100件が現実的
4. **TTL管理**: データの性質に応じた有効期限設定の重要性
5. **エラーハンドリング**: キャッシュ障害時の可用性確保

### 実際の開発での応用

- **管理画面**: ログ、監査ログ、アクティビティ履歴
- **ダッシュボード**: 統計情報、KPI、チャート
- **API**: よくアクセスされるエンドポイント
- **セッション**: ログイン状態、ショッピングカート
- **レート制限**: API制限、DDOS対策

### パフォーマンスチューニングの原則

1. **測定第一**: 推測せず実測する
2. **ボトルネック特定**: プロファイリングで遅い箇所を見つける
3. **適切な最適化**: 効果の大きい箇所から改善
4. **トレードオフ理解**: メモリ、一貫性、複雑さのバランス

Redisキャッシュは現代のWebアプリケーションに不可欠な技術です。このデモを通じて、キャッシングの基本と実践的なパターンを理解できました。
