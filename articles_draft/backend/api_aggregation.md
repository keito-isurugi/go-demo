## 概要

- Goで複数のAPIを並列で実行して結果を集約する実装
- マイクロサービスアーキテクチャにおけるAPIゲートウェイパターン
- sync.WaitGroupとmutexを使った安全な並行処理

### 特徴

- **並列API実行**: 複数のAPIを同時に呼び出して処理時間を短縮
- **結果の集約**: すべてのAPIレスポンスを1つのレスポンスにまとめる
- **柔軟な設定**: URL、メソッド、ヘッダー、ボディを個別に指定可能
- **タイムアウト制御**: 全体タイムアウトで長時間待機を防止
- **成功率の計算**: 成功・失敗をカウントして統計情報を提供
- **エラー分離**: 一部のAPIが失敗しても他の結果は返す

## API集約とは

### 基本概念

API集約（API Aggregation）は、複数の独立したAPIを並列で呼び出し、それらのレスポンスを1つのレスポンスにまとめて返すパターン。

**従来のアプローチ（クライアント側で複数リクエスト）**:
```
クライアント → API1 (2秒)
クライアント → API2 (2秒)
クライアント → API3 (2秒)
合計: 6秒 + ネットワーク往復 × 3
```

**API集約アプローチ（サーバー側で並列実行）**:
```
クライアント → 集約API → [API1, API2, API3] 並列実行 (2秒)
                      ← [結果1, 結果2, 結果3]
            ← 集約レスポンス
合計: 2秒 + ネットワーク往復 × 1
```

### メリット

1. **レスポンス時間の短縮**: 並列実行により最も遅いAPIの時間で完了
2. **ネットワーク効率**: クライアントからのリクエスト回数を削減
3. **データ整形**: サーバー側で統一フォーマットに変換
4. **エラーハンドリング**: 部分的な失敗を適切に処理
5. **認証の一元化**: 各APIへの認証をサーバー側で管理

### 使用場面

1. **ダッシュボード**: 複数のデータソースを集約して表示
2. **マイクロサービス**: 複数サービスのデータを1つのエンドポイントで提供
3. **BFF (Backend for Frontend)**: フロントエンド専用の集約レイヤー
4. **サードパーティAPI統合**: 複数の外部APIを組み合わせる

## サンプルコード

### 1. 基本的な構造体定義

```go
package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// AggregateAPIHandler は複数APIを並列で叩いて結果をまとめる
type AggregateAPIHandler struct{}

// APIRequest は個別APIリクエストの定義
type APIRequest struct {
	Name    string            `json:"name"`    // APIの識別名
	URL     string            `json:"url"`     // リクエストURL
	Method  string            `json:"method"`  // HTTPメソッド (GET, POST, etc.)
	Headers map[string]string `json:"headers"` // リクエストヘッダー
	Body    string            `json:"body"`    // リクエストボディ (POSTなど)
}

// AggregateRequest は集約リクエストの構造
type AggregateRequest struct {
	APIs    []APIRequest `json:"apis"`    // 並列実行するAPIのリスト
	Timeout int          `json:"timeout"` // タイムアウト（秒）
}

// APIResult は個別APIの結果
type APIResult struct {
	Name         string                 `json:"name"`
	URL          string                 `json:"url"`
	Method       string                 `json:"method"`
	StatusCode   int                    `json:"status_code"`
	ResponseTime int64                  `json:"response_time_ms"`
	Body         map[string]interface{} `json:"body,omitempty"`
	RawBody      string                 `json:"raw_body,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Success      bool                   `json:"success"`
}

// AggregateResponse は集約レスポンス
type AggregateResponse struct {
	Results       []APIResult `json:"results"`
	TotalAPIs     int         `json:"total_apis"`
	SuccessCount  int         `json:"success_count"`
	FailureCount  int         `json:"failure_count"`
	TotalDuration int64       `json:"total_duration_ms"`
}
```

### 2. メインハンドラー

```go
// AggregateHandler は複数APIを並列で叩いて結果を集約するハンドラ
func (h *AggregateAPIHandler) AggregateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var req AggregateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := map[string]interface{}{
			"error":   "Invalid JSON format",
			"message": "Please provide a JSON body with 'apis' array",
			"example": map[string]interface{}{
				"apis": []map[string]interface{}{
					{
						"name":   "google",
						"url":    "https://www.google.com",
						"method": "GET",
					},
					{
						"name":   "github",
						"url":    "https://api.github.com",
						"method": "GET",
						"headers": map[string]string{
							"Accept": "application/json",
						},
					},
				},
				"timeout": 10,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.APIs) == 0 {
		response := map[string]interface{}{
			"error": "No APIs provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// デフォルトのタイムアウトは30秒
	timeout := 30
	if req.Timeout > 0 && req.Timeout <= 120 {
		timeout = req.Timeout
	}

	startTime := time.Now()

	// 複数APIを並列実行
	results := h.executeAPIsParallel(req.APIs, timeout)

	// 成功・失敗をカウント
	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	response := AggregateResponse{
		Results:       results,
		TotalAPIs:     len(req.APIs),
		SuccessCount:  successCount,
		FailureCount:  failureCount,
		TotalDuration: time.Since(startTime).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

### 3. 並列実行ロジック

```go
// executeAPIsParallel は複数APIを並列で実行
func (h *AggregateAPIHandler) executeAPIsParallel(apis []APIRequest, timeoutSec int) []APIResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	results := make([]APIResult, len(apis))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, apiReq := range apis {
		wg.Add(1)
		go func(index int, api APIRequest) {
			defer wg.Done()
			result := h.executeAPI(ctx, api)

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, apiReq)
	}

	wg.Wait()
	return results
}
```

**ポイント**:
- `sync.WaitGroup`: すべてのgoroutineの完了を待つ
- `sync.Mutex`: スライスへの並行書き込みを保護
- `context.WithTimeout`: 全体のタイムアウトを制御
- インデックス保持: リクエスト順序とレスポンス順序を一致させる

### 4. 単一API実行

```go
// executeAPI は単一のAPIを実行
func (h *AggregateAPIHandler) executeAPI(ctx context.Context, apiReq APIRequest) APIResult {
	startTime := time.Now()

	// メソッドのデフォルトはGET
	method := apiReq.Method
	if method == "" {
		method = "GET"
	}

	// リクエストボディの設定
	var bodyReader io.Reader
	if apiReq.Body != "" {
		bodyReader = strings.NewReader(apiReq.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiReq.URL, bodyReader)
	if err != nil {
		return APIResult{
			Name:         apiReq.Name,
			URL:          apiReq.URL,
			Method:       method,
			ResponseTime: time.Since(startTime).Milliseconds(),
			Error:        err.Error(),
			Success:      false,
		}
	}

	// ヘッダーの設定
	for key, value := range apiReq.Headers {
		req.Header.Set(key, value)
	}

	// デフォルトでJSONを受け入れる
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return APIResult{
			Name:         apiReq.Name,
			URL:          apiReq.URL,
			Method:       method,
			ResponseTime: time.Since(startTime).Milliseconds(),
			Error:        err.Error(),
			Success:      false,
		}
	}
	defer resp.Body.Close()

	responseTime := time.Since(startTime).Milliseconds()

	// レスポンスボディを読み込み（最大10KB）
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024))
	if err != nil {
		return APIResult{
			Name:         apiReq.Name,
			URL:          apiReq.URL,
			Method:       method,
			StatusCode:   resp.StatusCode,
			ResponseTime: responseTime,
			Error:        "Failed to read response body: " + err.Error(),
			Success:      false,
		}
	}

	result := APIResult{
		Name:         apiReq.Name,
		URL:          apiReq.URL,
		Method:       method,
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
	}

	// JSONとしてパース試行
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
		result.Body = jsonBody
	} else {
		// JSON以外の場合はrawテキストとして保存
		result.RawBody = string(bodyBytes)
		// 長すぎる場合は省略
		if len(result.RawBody) > 500 {
			result.RawBody = result.RawBody[:500] + "... (truncated)"
		}
	}

	return result
}
```

**ポイント**:
- `io.LimitReader`: 大きなレスポンスからメモリを保護（最大10KB）
- JSON自動パース: レスポンスがJSONなら構造化データとして保存
- エラーの詳細記録: 各APIのエラーを個別に記録
- レスポンス時間計測: 各APIの所要時間を記録

### 5. プリセット集約エンドポイント

```go
// PresetAggregateHandler は事前定義された複数APIを並列で叩くハンドラ
func (h *AggregateAPIHandler) PresetAggregateHandler(w http.ResponseWriter, r *http.Request) {
	// 例: 公開APIを組み合わせて情報収集
	presetAPIs := []APIRequest{
		{
			Name:   "github_status",
			URL:    "https://www.githubstatus.com/api/v2/status.json",
			Method: "GET",
		},
		{
			Name:   "httpbin_uuid",
			URL:    "https://httpbin.org/uuid",
			Method: "GET",
		},
		{
			Name:   "httpbin_user_agent",
			URL:    "https://httpbin.org/user-agent",
			Method: "GET",
		},
		{
			Name:   "httpbin_headers",
			URL:    "https://httpbin.org/headers",
			Method: "GET",
		},
	}

	startTime := time.Now()
	results := h.executeAPIsParallel(presetAPIs, 30)

	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	response := AggregateResponse{
		Results:       results,
		TotalAPIs:     len(presetAPIs),
		SuccessCount:  successCount,
		FailureCount:  failureCount,
		TotalDuration: time.Since(startTime).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

## 動作原理

### 並列実行のフロー

```
クライアント → POST /api/aggregate
    ↓
リクエスト解析
    ↓
タイムアウト付きコンテキスト作成
    ↓
各APIに対してgoroutine起動
    ↓
┌─────────────┬─────────────┬─────────────┐
│ goroutine 0 │ goroutine 1 │ goroutine 2 │
│  API1実行   │  API2実行   │  API3実行   │
│  (1.5秒)    │  (2.3秒)    │  (0.8秒)    │
│      ↓      │      ↓      │      ↓      │
│  results[0] │  results[1] │  results[2] │
└──────┬──────┴──────┬──────┴──────┬──────┘
       │             │             │
       └──────→ WaitGroup.Wait() ←──┘
                     ↓
              統計情報の計算
                     ↓
              JSONレスポンス返却
```

### sync.WaitGroupの仕組み

```go
var wg sync.WaitGroup

// カウンタを増やす（起動するgoroutineの数だけ）
for i := 0; i < 3; i++ {
	wg.Add(1)  // カウンタ: 0 → 1 → 2 → 3
	go func() {
		defer wg.Done()  // 完了時にカウンタを減らす
		// 処理
	}()
}

// すべてのgoroutineの完了を待つ（カウンタが0になるまで）
wg.Wait()
```

**動作イメージ**:
```
時刻    カウンタ    イベント
t=0s    0          初期状態
t=0s    1          wg.Add(1) - goroutine 0起動
t=0s    2          wg.Add(1) - goroutine 1起動
t=0s    3          wg.Add(1) - goroutine 2起動
t=0.8s  2          goroutine 2完了 - wg.Done()
t=1.5s  1          goroutine 0完了 - wg.Done()
t=2.3s  0          goroutine 1完了 - wg.Done()
t=2.3s  -          wg.Wait()のブロック解除
```

### Mutexによる排他制御

```go
var mu sync.Mutex
results := make([]APIResult, len(apis))

go func(index int) {
	result := executeAPI(...)

	// 複数のgoroutineが同時にスライスに書き込む
	// → データ競合を防ぐためロック
	mu.Lock()
	results[index] = result  // 保護された書き込み
	mu.Unlock()
}(i)
```

**なぜMutexが必要か**:
```go
// データ競合の例（Mutexなし）
goroutine 0: results[0] = result0  ← 同時実行
goroutine 1: results[1] = result1  ← 同時実行
// スライスの内部構造が破損する可能性

// Mutexあり
goroutine 0: Lock() → results[0] = result0 → Unlock()
goroutine 1: (待機) → Lock() → results[1] = result1 → Unlock()
// 安全に書き込み
```

## 実装パターン

### 1. 条件付き実行（依存関係あり）

一部のAPIの結果に基づいて次のAPIを実行。

```go
func (h *AggregateAPIHandler) ConditionalAggregate(apis []APIRequest) []APIResult {
	results := make([]APIResult, 0)

	// 第1段階: ユーザー情報を取得
	userResult := h.executeAPI(context.Background(), apis[0])
	results = append(results, userResult)

	if !userResult.Success {
		return results // ユーザー取得失敗で終了
	}

	// 第2段階: ユーザーIDを使って並列で追加情報を取得
	var wg sync.WaitGroup
	var mu sync.Mutex

	userID := userResult.Body["id"].(string)

	for _, api := range apis[1:] {
		// URLにユーザーIDを埋め込み
		api.URL = strings.Replace(api.URL, "{user_id}", userID, -1)

		wg.Add(1)
		go func(a APIRequest) {
			defer wg.Done()
			result := h.executeAPI(context.Background(), a)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(api)
	}

	wg.Wait()
	return results
}
```

### 2. フォールバック（失敗時の代替）

プライマリAPIが失敗したらセカンダリAPIを試行。

```go
func (h *AggregateAPIHandler) FetchWithFallback(primary, secondary APIRequest) APIResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// プライマリAPIを試行
	result := h.executeAPI(ctx, primary)
	if result.Success {
		return result
	}

	// 失敗したらセカンダリAPIを試行
	result = h.executeAPI(ctx, secondary)
	result.Name = primary.Name + " (fallback)"
	return result
}

// 複数APIでフォールバック付き集約
func (h *AggregateAPIHandler) AggregateWithFallback(apiPairs [][2]APIRequest) []APIResult {
	results := make([]APIResult, len(apiPairs))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, pair := range apiPairs {
		wg.Add(1)
		go func(index int, primary, secondary APIRequest) {
			defer wg.Done()
			result := h.FetchWithFallback(primary, secondary)

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, pair[0], pair[1])
	}

	wg.Wait()
	return results
}
```

### 3. キャッシュ統合

よく使うAPIの結果をキャッシュ。

```go
type CachedAggregator struct {
	cache map[string]CacheEntry
	mu    sync.RWMutex
}

type CacheEntry struct {
	Result    APIResult
	ExpiresAt time.Time
}

func (ca *CachedAggregator) ExecuteWithCache(ctx context.Context, api APIRequest, ttl time.Duration) APIResult {
	cacheKey := api.URL

	// キャッシュチェック
	ca.mu.RLock()
	entry, exists := ca.cache[cacheKey]
	ca.mu.RUnlock()

	if exists && time.Now().Before(entry.ExpiresAt) {
		// キャッシュヒット
		entry.Result.Name = api.Name + " (cached)"
		return entry.Result
	}

	// キャッシュミス: APIを実行
	handler := &AggregateAPIHandler{}
	result := handler.executeAPI(ctx, api)

	// 成功した場合のみキャッシュに保存
	if result.Success {
		ca.mu.Lock()
		ca.cache[cacheKey] = CacheEntry{
			Result:    result,
			ExpiresAt: time.Now().Add(ttl),
		}
		ca.mu.Unlock()
	}

	return result
}
```

### 4. レート制限付き実行

同時実行数を制限。

```go
import "golang.org/x/sync/semaphore"

func (h *AggregateAPIHandler) AggregateWithRateLimit(apis []APIRequest, maxConcurrent int64) []APIResult {
	ctx := context.Background()
	sem := semaphore.NewWeighted(maxConcurrent)

	results := make([]APIResult, len(apis))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, api := range apis {
		wg.Add(1)
		go func(index int, a APIRequest) {
			defer wg.Done()

			// セマフォ取得（最大maxConcurrent個まで）
			if err := sem.Acquire(ctx, 1); err != nil {
				mu.Lock()
				results[index] = APIResult{
					Name:    a.Name,
					Error:   err.Error(),
					Success: false,
				}
				mu.Unlock()
				return
			}
			defer sem.Release(1)

			// API実行
			result := h.executeAPI(ctx, a)

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, api)
	}

	wg.Wait()
	return results
}
```

## 使いどころ

### 向いている場面

1. **ユーザーダッシュボード**: 複数サービスのデータを1画面に表示
   ```go
   apis := []APIRequest{
       {Name: "profile", URL: userServiceURL + "/profile"},
       {Name: "orders", URL: orderServiceURL + "/orders"},
       {Name: "notifications", URL: notificationServiceURL + "/unread"},
       {Name: "recommendations", URL: recommendServiceURL + "/suggestions"},
   }
   ```

2. **価格比較サイト**: 複数ECサイトの価格を並列取得
   ```go
   apis := []APIRequest{
       {Name: "amazon", URL: "https://api.amazon.com/product/123"},
       {Name: "ebay", URL: "https://api.ebay.com/item/123"},
       {Name: "walmart", URL: "https://api.walmart.com/products/123"},
   }
   ```

3. **監視ダッシュボード**: 複数システムの状態を集約
   ```go
   apis := []APIRequest{
       {Name: "web_server", URL: "http://web1/health"},
       {Name: "api_server", URL: "http://api1/health"},
       {Name: "database", URL: "http://db1/health"},
       {Name: "cache", URL: "http://redis1/ping"},
   }
   ```

4. **BFF (Backend for Frontend)**: モバイルアプリ用の最適化エンドポイント
   ```go
   // 1回のリクエストでアプリ起動に必要なすべてのデータを取得
   apis := []APIRequest{
       {Name: "user", URL: "/api/user"},
       {Name: "settings", URL: "/api/settings"},
       {Name: "feed", URL: "/api/feed"},
       {Name: "notifications_count", URL: "/api/notifications/count"},
   }
   ```

### 向いていない場面

1. **順序依存の処理**: 前のAPIの結果が次のAPIに必須
2. **大量のAPI**: 数百以上のAPIを並列実行（リソース制約）
3. **ストリーミング**: 長時間接続が必要な処理

## 実用例

### 1. Eコマースの商品詳細ページ

```go
func (h *AggregateAPIHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")

	apis := []APIRequest{
		{
			Name:   "product_info",
			URL:    fmt.Sprintf("http://product-service/api/products/%s", productID),
			Method: "GET",
		},
		{
			Name:   "reviews",
			URL:    fmt.Sprintf("http://review-service/api/reviews?product=%s&limit=5", productID),
			Method: "GET",
		},
		{
			Name:   "inventory",
			URL:    fmt.Sprintf("http://inventory-service/api/stock/%s", productID),
			Method: "GET",
		},
		{
			Name:   "recommendations",
			URL:    fmt.Sprintf("http://recommend-service/api/similar?product=%s&limit=4", productID),
			Method: "GET",
		},
		{
			Name:   "pricing",
			URL:    fmt.Sprintf("http://pricing-service/api/price/%s", productID),
			Method: "GET",
		},
	}

	results := h.executeAPIsParallel(apis, 5)

	// フロントエンド用に整形
	response := map[string]interface{}{
		"product":         findResult(results, "product_info"),
		"reviews":         findResult(results, "reviews"),
		"in_stock":        findResult(results, "inventory"),
		"recommendations": findResult(results, "recommendations"),
		"pricing":         findResult(results, "pricing"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func findResult(results []APIResult, name string) interface{} {
	for _, r := range results {
		if r.Name == name && r.Success {
			return r.Body
		}
	}
	return map[string]string{"error": "not available"}
}
```

### 2. 天気情報の統合

```go
func (h *AggregateAPIHandler) GetWeatherAggregate(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")

	// 複数の天気APIから情報を取得して信頼性向上
	apis := []APIRequest{
		{
			Name:   "openweather",
			URL:    fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=KEY", location),
			Method: "GET",
		},
		{
			Name:   "weatherapi",
			URL:    fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=KEY&q=%s", location),
			Method: "GET",
		},
		{
			Name:   "meteo",
			URL:    fmt.Sprintf("https://api.open-meteo.com/v1/forecast?location=%s", location),
			Method: "GET",
		},
	}

	results := h.executeAPIsParallel(apis, 10)

	// 複数ソースから平均値を計算
	temps := []float64{}
	for _, r := range results {
		if r.Success {
			if temp, ok := extractTemperature(r.Body); ok {
				temps = append(temps, temp)
			}
		}
	}

	avgTemp := average(temps)

	response := map[string]interface{}{
		"location":       location,
		"temperature":    avgTemp,
		"sources":        len(temps),
		"raw_results":    results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

### 3. ソーシャルメディア統合ダッシュボード

```go
func (h *AggregateAPIHandler) GetSocialMediaDashboard(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Get("Authorization")

	apis := []APIRequest{
		{
			Name:   "twitter",
			URL:    "https://api.twitter.com/2/users/me",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": userToken,
			},
		},
		{
			Name:   "facebook",
			URL:    "https://graph.facebook.com/me",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": userToken,
			},
		},
		{
			Name:   "instagram",
			URL:    "https://graph.instagram.com/me",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": userToken,
			},
		},
		{
			Name:   "linkedin",
			URL:    "https://api.linkedin.com/v2/me",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": userToken,
			},
		},
	}

	results := h.executeAPIsParallel(apis, 15)

	// 各SNSのフォロワー数など統計を集約
	dashboard := map[string]interface{}{
		"connected_accounts": successCount(results),
		"accounts":           results,
		"last_updated":       time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func successCount(results []APIResult) int {
	count := 0
	for _, r := range results {
		if r.Success {
			count++
		}
	}
	return count
}
```

## まとめ

### メリット

- **高速化**: 並列実行により劇的に高速（順次実行の1/n）
- **ネットワーク効率**: クライアントのリクエスト数削減
- **部分的失敗への耐性**: 一部が失敗しても他の結果は返せる
- **統一インターフェース**: 複数APIを1つのエンドポイントで提供
- **データ整形**: サーバー側で統一フォーマットに変換可能

### デメリット

- **複雑性**: 並行処理特有のバグ（データ競合など）
- **デバッグ困難**: 複数のAPIが絡むとトラブルシューティングが難しい
- **遅延の伝播**: 最も遅いAPIが全体の速度を決定
- **メモリ使用量**: 全結果を保持するため大量データ時は注意

### 使うべき時

- **マイクロサービス**: 複数サービスのデータを集約
- **モバイルBFF**: アプリ向けに最適化したエンドポイント
- **ダッシュボード**: 複数ソースのデータを1画面に表示
- **外部API統合**: サードパーティAPIを組み合わせる

### ベストプラクティス

1. **タイムアウト設定**: 必ずcontext.WithTimeoutを使用
2. **エラー分離**: 一部の失敗が全体に影響しないよう設計
3. **レート制限**: 大量APIはセマフォで同時実行数を制限
4. **キャッシュ活用**: 頻繁に呼ばれるAPIはキャッシュ
5. **監視とログ**: 各APIの成功率・レスポンス時間を記録
6. **リトライ戦略**: 一時的な失敗に対してリトライ

### パターン選択

| ユースケース | パターン | 理由 |
| --- | --- | --- |
| 独立したAPI群 | 基本並列実行 | シンプルで高速 |
| 依存関係あり | 条件付き実行 | 段階的に処理 |
| 高可用性重視 | フォールバック | 冗長性確保 |
| 頻繁なアクセス | キャッシュ統合 | レスポンス時間短縮 |
| 大量API | レート制限付き | リソース保護 |

API集約はマイクロサービスアーキテクチャにおいて重要なパターン。Goのgoroutineとsync.WaitGroupを使えば、安全で効率的な並列API実行を簡単に実装できる。クライアントからのリクエスト数を削減し、レスポンス時間を短縮することで、優れたユーザー体験を提供できる。
