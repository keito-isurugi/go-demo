## 概要

- Goでレートリミッター（Rate Limiter）を実装
- スライディングウィンドウ方式を使用したHTTPミドルウェア
- 特定の時間枠内のリクエスト数を制限し、サービスを過負荷から保護

### 特徴

- **スライディングウィンドウ方式**: 時間枠を固定せず、常に直近N秒のリクエストを計算
- **IPアドレスベース**: クライアントごとに独立したレート制限
- **スレッドセーフ**: sync.Mutexによる並行処理対応
- **自動クリーンアップ**: 古いエントリを定期的に削除してメモリリークを防止
- **HTTPミドルウェア**: 既存のHTTPハンドラーに簡単に組み込み可能

## レートリミットとは

### 基本概念

レートリミット（Rate Limiting）は、特定の時間内に実行できるリクエスト数を制限する技術。APIやWebサービスの過負荷を防ぎ、公平なリソース利用を実現する。

**典型的な制限例**:
```
- 1分間に10リクエストまで
- 1時間に100リクエストまで
- 1日に1000リクエストまで
```

### レートリミットの目的

1. **DoS攻撃対策**: 大量リクエストによるサービス停止を防止
2. **リソース保護**: データベースやAPIの過負荷を防止
3. **公平性の確保**: 特定ユーザーによるリソース独占を防止
4. **コスト管理**: 外部APIの従量課金コストを制御
5. **サービス品質維持**: すべてのユーザーに安定したレスポンスを提供

### レートリミットのアプローチ

| アルゴリズム | 説明 | メリット | デメリット |
| --- | --- | --- | --- |
| **固定ウィンドウ** | 固定の時間枠でカウント | 実装が簡単 | 境界でバースト発生 |
| **スライディングウィンドウ** | 常に直近N秒を計算 | 正確な制限 | メモリ使用量が多い |
| **トークンバケット** | トークンを消費して実行 | バースト許容 | 実装が複雑 |
| **リーキーバケット** | 一定レートで処理 | 平滑化が得意 | 柔軟性が低い |

今回実装するのは**スライディングウィンドウ方式**。

## サンプルコード

### 1. レートリミッター構造体

```go
package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time  // IPアドレス -> タイムスタンプのリスト
	mu       sync.Mutex              // 並行アクセス制御
	limit    int                     // 制限回数
	window   time.Duration           // 時間枠
}
```

**構造体の役割**:
- `requests`: 各IPアドレスのリクエスト履歴を保存
- `mu`: マップへの並行アクセスを保護
- `limit`: 許可するリクエスト数（例: 10回）
- `window`: 時間枠（例: 1分間）

### 2. 初期化とクリーンアップ

```go
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// クリーンアップ処理を定期的に実行
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, timestamps := range rl.requests {
			// 古いエントリを削除
			validTimestamps := []time.Time{}
			for _, t := range timestamps {
				if now.Sub(t) < rl.window {
					validTimestamps = append(validTimestamps, t)
				}
			}
			if len(validTimestamps) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validTimestamps
			}
		}
		rl.mu.Unlock()
	}
}
```

**ポイント**:
- `NewRateLimiter`: コンストラクタで初期化
- `cleanup()`: バックグラウンドで定期的に古いエントリを削除
- `time.NewTicker`: ウィンドウ時間ごとにクリーンアップ実行
- メモリリークを防ぐため、不要なIPアドレスのエントリを削除

### 3. HTTPミドルウェア実装

```go
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()

		// IPアドレスのリクエスト履歴を取得
		timestamps, exists := rl.requests[ip]
		if !exists {
			timestamps = []time.Time{}
		}

		// ウィンドウ内のリクエストのみをフィルタ
		validTimestamps := []time.Time{}
		for _, t := range timestamps {
			if now.Sub(t) < rl.window {
				validTimestamps = append(validTimestamps, t)
			}
		}

		// レート制限チェック
		if len(validTimestamps) >= rl.limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Rate limit exceeded. Maximum 10 requests per minute."}`))
			return
		}

		// 新しいリクエストを記録
		validTimestamps = append(validTimestamps, now)
		rl.requests[ip] = validTimestamps

		next.ServeHTTP(w, r)
	})
}
```

**処理の流れ**:
1. クライアントのIPアドレスを取得
2. ロックを取得して並行アクセスを制御
3. 現在時刻と過去のリクエスト履歴を比較
4. ウィンドウ内のリクエストのみをフィルタリング
5. リクエスト数が制限を超えている場合、429エラーを返す
6. 制限内であれば、タイムスタンプを記録して次のハンドラーに処理を委譲

### 4. 使用例

```go
package main

import (
	"fmt"
	"net/http"
	"time"
	"your-project/middleware"
)

func main() {
	// レートリミッター作成: 1分間に10リクエストまで
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)

	// ハンドラー作成
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Success"}`))
	})

	// レートリミッターをミドルウェアとして適用
	handler := rateLimiter.Middleware(mux)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}
```

### 5. テストコード

```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// 3リクエスト/秒の制限
	rl := NewRateLimiter(3, time.Second)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 同じIPからのリクエストをシミュレート
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// 最初の3リクエストは成功
	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("Request %d should succeed, got status %d", i+1, rec.Code)
		}
	}

	// 4番目のリクエストは制限される
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Request should be rate limited, got status %d", rec.Code)
	}

	// 1秒待つとリセットされる
	time.Sleep(time.Second)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Request after window should succeed, got status %d", rec.Code)
	}
}
```

## 動作原理

### スラ��ディングウィンドウの仕組み

固定ウィンドウとの違いを理解するために比較してみましょう。

#### 固定ウィンドウの問題点

```
制限: 10リクエスト/分

時刻    0:00              0:30              1:00              1:30
        |-----------------|-----------------|-----------------|
        [    ウィンドウ1    ][    ウィンドウ2    ]

0:59に10リクエスト → OK
1:00に10リクエスト → OK （新しいウィンドウ）
→ 1秒間に20リクエストが通過！
```

**問題**: ウィンドウの境界で大量のリクエストが集中する可能性がある。

#### スライディングウィンドウの解決

```
制限: 10リクエスト/分

現在時刻: 12:05:30
スライディングウィンドウ: 12:04:30 ~ 12:05:30 の直近60秒

リクエスト履歴:
12:04:35 ✓
12:04:40 ✓
12:04:45 ✗ (60秒以上前なので除外)
12:04:50 ✓
12:05:00 ✓
12:05:10 ✓
12:05:25 ✓
→ 有効なリクエスト数: 6個 → 制限内
```

**利点**: 常に正確な直近N秒を計算するため、境界問題が発生しない。

### 実際の動作例

**シナリオ**: 1分間に3リクエストまでの制限

```go
// 初期状態
requests["192.168.1.1"] = []

// 12:00:00 - リクエスト1
timestamps: [12:00:00]
→ 1個 < 3 → 許可

// 12:00:20 - リクエスト2
timestamps: [12:00:00, 12:00:20]
→ 2個 < 3 → 許可

// 12:00:40 - リクエスト3
timestamps: [12:00:00, 12:00:20, 12:00:40]
→ 3個 = 3 → 許可

// 12:00:50 - リクエスト4
validTimestamps: [12:00:00, 12:00:20, 12:00:40] (すべて60秒以内)
→ 3個 >= 3 → 拒否（429エラー）

// 12:01:10 - リクエスト5
validTimestamps: [12:00:20, 12:00:40] (12:00:00は60秒以上前なので除外)
→ 2個 < 3 → 許可
timestamps: [12:00:20, 12:00:40, 12:01:10]
```

### クリーンアップの動作

```go
// メモリ使用量を抑えるため定期的にクリーンアップ

現在時刻: 12:05:00
ウィンドウ: 60秒

IPアドレス: 192.168.1.1
timestamps: [12:03:30, 12:04:50]
→ 12:03:30は60秒以上前 → 削除
→ validTimestamps: [12:04:50]

IPアドレス: 192.168.1.2
timestamps: [12:03:00, 12:03:20]
→ すべて60秒以上前 → IPごと削除
```

## 計算量

### 時間計算量

| 操作 | 計算量 | 説明 |
| --- | --- | --- |
| リクエスト処理 | O(n) | n = ウィンドウ内のリクエスト数（最大limit個） |
| クリーンアップ | O(m × n) | m = IPアドレス数、n = 平均タイムスタンプ数 |

**実用的には**:
- リクエスト処理: limitが小さい（10~100程度）ためほぼO(1)
- クリーンアップ: バックグラウンドで実行されるため影響小

### 空間計算量

| 項目 | 計算量 | 説明 |
| --- | --- | --- |
| 各IPのタイムスタンプ | O(limit) | 最大limit個のタイムスタンプ |
| 全体 | O(m × limit) | m = アクティブなIPアドレス数 |

**メモリ使用量の推定**:
```
仮定:
- 同時アクティブIP: 1,000個
- limit: 10リクエスト/分
- タイムスタンプ: 24バイト（time.Time）

メモリ使用量 ≈ 1,000 × 10 × 24バイト = 240KB
→ 非常に軽量
```

## 改善と最適化

### 1. 分散環境対応（Redis使用）

単一サーバーではメモリ上のマップで十分だが、複数サーバーでは共有ストレージが必要。

```go
type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func (rl *RedisRateLimiter) Allow(ip string) bool {
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	key := fmt.Sprintf("rate_limit:%s", ip)

	// ウィンドウ外のエントリを削除
	rl.client.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", windowStart))

	// 現在のカウント取得
	count, _ := rl.client.ZCard(key).Result()

	if count >= int64(rl.limit) {
		return false
	}

	// 新しいリクエストを記録
	rl.client.ZAdd(key, redis.Z{Score: float64(now), Member: now})
	rl.client.Expire(key, rl.window)

	return true
}
```

**メリット**:
- 複数サーバー間で制限を共有
- 永続化により再起動時もデータ保持
- Redisの高速性を活用

### 2. トークンバケットアルゴリズム

よりスムーズなレート制限を実現。

```go
type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64  // トークン/秒
	lastRefill time.Time
	mu         sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 経過時間に応じてトークンを補充
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now

	// トークンが1つ以上あれば許可
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}
```

**メリット**:
- バースト的なリクエストを許容
- より自然なトラフィック制御
- メモリ効率が良い（タイムスタンプ不要）

### 3. ユーザー認証ベースの制限

IPアドレスではなくユーザーIDで制限。

```go
func (rl *RateLimiter) MiddlewareWithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 認証トークンからユーザーIDを取得
		userID := getUserIDFromToken(r.Header.Get("Authorization"))
		if userID == "" {
			userID = r.RemoteAddr // フォールバック
		}

		// userIDをキーとして制限
		if !rl.checkLimit(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

**メリット**:
- プロキシやVPN経由でも正確に制限
- ユーザーごとに異なる制限を設定可能
- より公平なリソース配分

### 4. レスポンスヘッダーの追加

クライアントに制限情報を提供。

```go
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		remaining, resetTime := rl.checkLimitWithInfo(ip)

		// レート制限情報をヘッダーに追加
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		if remaining <= 0 {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Until(resetTime).Seconds())))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

**ヘッダーの意味**:
- `X-RateLimit-Limit`: 最大リクエスト数
- `X-RateLimit-Remaining`: 残りリクエスト数
- `X-RateLimit-Reset`: リセット時刻（Unixタイムスタンプ）
- `Retry-After`: 再試行までの秒数

## 使いどころ

### 向いている場面

1. **公開API**: 外部クライアントからの過剰アクセスを防止
   ```go
   // 無料プラン: 60リクエスト/時間
   freeLimiter := NewRateLimiter(60, time.Hour)

   // 有料プラン: 1000リクエスト/時間
   paidLimiter := NewRateLimiter(1000, time.Hour)
   ```

2. **ログインエンドポイント**: ブルートフォース攻撃対策
   ```go
   // 5分間に5回までの失敗ログインを許可
   loginLimiter := NewRateLimiter(5, 5*time.Minute)
   ```

3. **データ集約API**: 高コストな処理の保護
   ```go
   // 大規模レポート生成: 1時間に3回まで
   reportLimiter := NewRateLimiter(3, time.Hour)
   ```

4. **WebSocket接続**: 接続数の制限
   ```go
   // 同時接続数を制限
   connectionLimiter := NewRateLimiter(100, time.Second)
   ```

### 向いていない���面

1. **内部API**: 信頼できる内部サービス間の通信
2. **静的ファイル**: CDNやWebサーバーレベルで処理すべき
3. **リアルタイム通信**: ゲームやチャットなど遅延が致命的な場合

## 実例

### 1. GitHub API風のレート制限

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIServer struct {
	limiter *RateLimiter
}

func NewAPIServer() *APIServer {
	return &APIServer{
		limiter: NewRateLimiter(5000, time.Hour), // 5000リクエスト/時間
	}
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	if !s.limiter.Allow(ip) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "API rate limit exceeded",
			"documentation_url": "https://docs.example.com/rate-limiting",
		})
		return
	}

	// 実際のAPI処理
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success",
	})
}
```

### 2. エンドポイントごとの制限

```go
type MultiRateLimiter struct {
	limiters map[string]*RateLimiter
}

func NewMultiRateLimiter() *MultiRateLimiter {
	return &MultiRateLimiter{
		limiters: map[string]*RateLimiter{
			"/api/search":  NewRateLimiter(10, time.Minute),   // 検索: 10/分
			"/api/create":  NewRateLimiter(5, time.Minute),    // 作成: 5/分
			"/api/delete":  NewRateLimiter(3, time.Minute),    // 削除: 3/分
			"/api/default": NewRateLimiter(100, time.Minute),  // デフォルト: 100/分
		},
	}
}

func (m *MultiRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter, exists := m.limiters[r.URL.Path]
		if !exists {
			limiter = m.limiters["/api/default"]
		}

		// エンドポイント固有のレート制限を適用
		wrappedHandler := limiter.Middleware(next)
		wrappedHandler.ServeHTTP(w, r)
	})
}
```

### 3. 段階的な制限（ティア制）

```go
type TieredRateLimiter struct {
	tiers map[string]*RateLimiter
}

func NewTieredRateLimiter() *TieredRateLimiter {
	return &TieredRateLimiter{
		tiers: map[string]*RateLimiter{
			"free":       NewRateLimiter(100, time.Hour),   // 無料: 100/時間
			"basic":      NewRateLimiter(1000, time.Hour),  // ベーシック: 1000/時間
			"premium":    NewRateLimiter(10000, time.Hour), // プレミアム: 10000/時間
			"enterprise": NewRateLimiter(100000, time.Hour), // エンタープライズ: 100000/時間
		},
	}
}

func (t *TieredRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// APIキーからティアを判定
		apiKey := r.Header.Get("X-API-Key")
		tier := getUserTier(apiKey) // "free", "basic", "premium", "enterprise"

		limiter, exists := t.tiers[tier]
		if !exists {
			limiter = t.tiers["free"]
		}

		wrappedHandler := limiter.Middleware(next)
		wrappedHandler.ServeHTTP(w, r)
	})
}
```

## まとめ

### メリット

- **DoS攻撃対策**: 大量リクエストからサービスを保護
- **リソース保護**: データベースやAPIの負荷を制御
- **公平性**: すべてのユーザーに均等なアクセスを保証
- **実装が簡単**: 標準ライブラリのみで実装可能
- **柔軟性**: ウィンドウサイズや制限数を自由に設定

### デメリット

- **メモリ使用量**: タイムスタンプを保存するため、アクティブユーザーが多いと増加
- **単一サーバー制限**: 分散環境ではRedisなど外部ストレージが必要
- **時刻同期**: サーバー間で時刻がずれると正確な制限が困難
- **正当なユーザーへの影響**: 同じIPを共有する場合、他のユーザーの影響を受ける

### 使うべき時

- **公開API提供時**: 必須の機能
- **認証エンドポイント**: ブルートフォース攻撃対策
- **高コストな処理**: レポート生成、データ集約など
- **マイクロサービス**: 各サービスの保護

### 選択基準

| シナリオ | 推奨アルゴリズム | 理由 |
| --- | --- | --- |
| 単一サーバーAPI | スライディングウィンドウ | 正確で実装が簡単 |
| 分散環境 | Redis + スライディングウィンドウ | サーバー間で共有可能 |
| バースト許容 | トークンバケット | 一時的な高負荷を吸収 |
| 厳密な平滑化 | リーキーバケット | 常に一定レートを維持 |
| 低メモリ環境 | 固定ウィンドウ | メモリ効率が最高 |

レートリミットはモダンなWebサービスに不可欠な機能。適切に実装することで、サービスの安定性と公平性を大幅に向上させることができる。Goの並行処理機能とHTTPミドルウェアパターンを活用すれば、わずか100行程度のコードで本格的なレート制限を実装できる。
