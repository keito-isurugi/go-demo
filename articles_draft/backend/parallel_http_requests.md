## 概要

- Goで複数のHTTPリクエストを並列処理する実装
- goroutineとchannelを使った並行プログラミングのパターン
- 最速レスポンスの取得、全結果の収集など実用的なユースケース

### 特徴

- **並列リクエスト**: 複数URLを同時にGETして処理時間を短縮
- **最速レスポンス取得**: 最初に成功したレスポンスを即座に返す
- **タイムアウト制御**: context.Contextによる一括タイムアウト管理
- **全結果収集**: すべてのgoroutineの完了を待って結果を集約
- **エラーハンドリング**: 個別のエラーを適切に処理

## 並列HTTPリクエストとは

### 基本概念

複数のHTTPリクエストを順次実行する代わりに、goroutineを使って並列実行することで、処理時間を大幅に短縮する技術。

**順次実行 vs 並列実行**:
```
順次実行（3つのURL、各2秒）:
URL1 [==] → URL2 [==] → URL3 [==]
合計: 6秒

並列実行（3つのURL、各2秒）:
URL1 [==]
URL2 [==]
URL3 [==]
合計: 2秒（最も遅いリクエスト）
```

### 使用場面

1. **複数APIの集約**: マイクロサービスから並列でデータ取得
2. **ミラーサーバー**: 複数サーバーから最速レスポンスを取得
3. **ヘルスチェック**: 複数エンドポイントを並列監視
4. **データ収集**: Webスクレイピングの高速化
5. **フェイルオーバー**: プライマリが遅い場合にセカンダリを使用

## サンプルコード

### 1. 基本的な構造体定義

```go
package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// ParallelFetchHandler は並列フェッチのハンドラ
type ParallelFetchHandler struct{}

// URLRequest はリクエストボディの構造
type URLRequest struct {
	URLs    []string `json:"urls"`
	Timeout int      `json:"timeout"` // タイムアウト（秒）
}

// URLResponse は個別のURL取得結果
type URLResponse struct {
	URL          string `json:"url"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	Body         string `json:"body,omitempty"`
	Error        string `json:"error,omitempty"`
}

// FastestResponseResult は最終的なレスポンス
type FastestResponseResult struct {
	Fastest       URLResponse   `json:"fastest"`
	AllResults    []URLResponse `json:"all_results"`
	TotalDuration int64         `json:"total_duration_ms"`
}
```

### 2. 最速レスポンスを取得

複数URLから最初に成功したレスポンスを返すパターン。

```go
// FetchFastestHandler は最速のレスポンスを返すハンドラ
func (h *ParallelFetchHandler) FetchFastestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := map[string]interface{}{
			"error":   "Invalid JSON format",
			"message": "Please provide a JSON body with 'urls' array and optional 'timeout' (in seconds)",
			"example": map[string]interface{}{
				"urls":    []string{"https://example.com", "https://google.com"},
				"timeout": 10,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.URLs) == 0 {
		response := map[string]interface{}{
			"error": "No URLs provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// デフォルトのタイムアウトは10秒
	timeout := 10
	if req.Timeout > 0 && req.Timeout <= 60 {
		timeout = req.Timeout
	}

	startTime := time.Now()

	// 最速のレスポンスを取得
	result := h.fetchURLsParallel(req.URLs, timeout)
	result.TotalDuration = time.Since(startTime).Milliseconds()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
```

### 3. 並列フェッチの実装

```go
// fetchURLsParallel は複数のURLを並列で取得し、最速のレスポンスを返す
func (h *ParallelFetchHandler) fetchURLsParallel(urls []string, timeoutSec int) FastestResponseResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// 最速レスポンス用のチャネル（バッファサイズ1）
	fastestChan := make(chan URLResponse, 1)
	// 全結果収集用のチャネル
	resultsChan := make(chan URLResponse, len(urls))

	// 全URLを並列で取得
	for _, url := range urls {
		go func(u string) {
			result := h.fetchURL(ctx, u)

			// 最初の成功レスポンスを最速チャネルに送信（試行のみ）
			if result.Error == "" {
				select {
				case fastestChan <- result:
					// 最速レスポンスとして送信成功
				default:
					// すでに最速レスポンスが送信済み
				}
			}

			// 全結果チャネルには必ず送信
			resultsChan <- result
		}(url)
	}

	// 最速レスポンスを待機
	var fastest URLResponse
	select {
	case fastest = <-fastestChan:
		// 最速レスポンスを取得
	case <-ctx.Done():
		// タイムアウト
		fastest = URLResponse{
			Error: "All requests timed out or failed",
		}
	}

	// 全結果を収集（最速取得後も他のgoroutineの完了を待つ）
	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		select {
		case result := <-resultsChan:
			allResults = append(allResults, result)
		case <-time.After(time.Duration(timeoutSec) * time.Second):
			// タイムアウト
		}
	}

	return FastestResponseResult{
		Fastest:    fastest,
		AllResults: allResults,
	}
}
```

**ポイント**:
- `fastestChan`: バッファサイズ1で最初の成功レスポンスのみ受信
- `resultsChan`: バッファサイズlen(urls)で全結果を受信
- `select`と`default`: 既に送信済みの場合はブロックしない
- `ctx.Done()`: タイムアウト時の処理

### 4. 単一URLの取得

```go
// fetchURL は単一のURLを取得
func (h *ParallelFetchHandler) fetchURL(ctx context.Context, url string) URLResponse {
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return URLResponse{
			URL:          url,
			ResponseTime: time.Since(startTime).Milliseconds(),
			Error:        err.Error(),
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return URLResponse{
			URL:          url,
			ResponseTime: time.Since(startTime).Milliseconds(),
			Error:        err.Error(),
		}
	}
	defer resp.Body.Close()

	responseTime := time.Since(startTime).Milliseconds()

	// レスポンスボディを読み込み（最大1KB）
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	body := string(bodyBytes)
	if err != nil {
		body = ""
	}

	return URLResponse{
		URL:          url,
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
		Body:         body,
	}
}
```

**ポイント**:
- `http.NewRequestWithContext`: コンテキストでタイムアウト制御
- `io.LimitReader`: 大きなレスポンスからメモリを保護（最大1KB）
- レスポンス時間の計測

### 5. 全結果を取得

すべてのURLの結果を並列で取得して返すパターン。

```go
// FetchAllHandler は全URLを並列で取得し、全結果を返すハンドラ
func (h *ParallelFetchHandler) FetchAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := map[string]interface{}{
			"error":   "Invalid JSON format",
			"message": "Please provide a JSON body with 'urls' array and optional 'timeout' (in seconds)",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.URLs) == 0 {
		response := map[string]interface{}{
			"error": "No URLs provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	timeout := 10
	if req.Timeout > 0 && req.Timeout <= 60 {
		timeout = req.Timeout
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 全結果収集用のチャネル
	resultsChan := make(chan URLResponse, len(req.URLs))

	// 全URLを並列で取得
	for _, url := range req.URLs {
		go func(u string) {
			resultsChan <- h.fetchURL(ctx, u)
		}(url)
	}

	// 全結果を収集
	allResults := make([]URLResponse, 0, len(req.URLs))
	for i := 0; i < len(req.URLs); i++ {
		select {
		case result := <-resultsChan:
			allResults = append(allResults, result)
		case <-ctx.Done():
			break
		}
	}

	response := map[string]interface{}{
		"results":        allResults,
		"total_urls":     len(req.URLs),
		"total_duration": time.Since(startTime).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

## 動作原理

### 並列処理のフロー

```
リクエスト受信
    ↓
タイムアウト付きコンテキスト作成
    ↓
各URLに対してgoroutine起動
    ↓
┌─────────────┬─────────────┬─────────────┐
│ goroutine 1 │ goroutine 2 │ goroutine 3 │
│  GET URL1   │  GET URL2   │  GET URL3   │
│   (2秒)     │   (1秒)     │   (3秒)     │
└──────┬──────┴──────┬──────┴──────┬──────┘
       │             │ (最速)      │
       │             ↓             │
       │      fastestChan          │
       │             ↓             │
       ├────→ resultsChan ←────────┤
       ↓             ↓             ↓
     全結果を収集して返す
```

### チャネルの使い分け

**最速レスポンス用チャネル**:
```go
fastestChan := make(chan URLResponse, 1)

// 送信側（各goroutine）
if result.Error == "" {
	select {
	case fastestChan <- result:
		// 最初の送信は成功
	default:
		// 2番目以降は無視
	}
}

// 受信側（メインgoroutine）
select {
case fastest = <-fastestChan:
	// 最速レスポンスを取得
case <-ctx.Done():
	// タイムアウト
}
```

**全結果収集用チャネル**:
```go
resultsChan := make(chan URLResponse, len(urls))

// 送信側（各goroutine）
resultsChan <- result  // 必ず送信

// 受信側（メインgoroutine）
for i := 0; i < len(urls); i++ {
	result := <-resultsChan
	allResults = append(allResults, result)
}
```

### タイムアウトの仕組み

```go
// コンテキストによるタイムアウト制御
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 各goroutineでコンテキストを使用
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

// タイムアウト時の動作
select {
case result := <-resultsChan:
	// 正常終了
case <-ctx.Done():
	// タイムアウトまたはキャンセル
}
```

**タイムアウトの流れ**:
```
t=0s:  コンテキスト作成（10秒タイムアウト）
t=1s:  URL2のgoroutineが完了 → fastestChanに送信
t=2s:  URL1のgoroutineが完了 → resultsChanに送信
t=10s: タイムアウト発生
       → ctx.Done()がクローズ
       → 未完了のgoroutineは自動キャンセル
```

## 実装パターン

### 1. Fan-Out / Fan-In パターン

複数のgoroutineを起動（Fan-Out）して結果を集約（Fan-In）。

```go
// Fan-Out: 複数のgoroutineを起動
func FanOut(urls []string) []<-chan URLResponse {
	channels := make([]<-chan URLResponse, len(urls))
	for i, url := range urls {
		ch := make(chan URLResponse, 1)
		channels[i] = ch
		go func(u string, c chan URLResponse) {
			c <- fetchURL(u)
			close(c)
		}(url, ch)
	}
	return channels
}

// Fan-In: 複数のチャネルを1つに集約
func FanIn(channels []<-chan URLResponse) <-chan URLResponse {
	out := make(chan URLResponse)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan URLResponse) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// 使用例
channels := FanOut(urls)
results := FanIn(channels)
for result := range results {
	fmt.Println(result)
}
```

### 2. ワーカープールパターン

並列度を制限してリソースを保護。

```go
type WorkerPool struct {
	workerCount int
}

func (wp *WorkerPool) FetchURLs(urls []string) []URLResponse {
	jobs := make(chan string, len(urls))
	results := make(chan URLResponse, len(urls))

	// ワーカー起動（最大workerCount個）
	for i := 0; i < wp.workerCount; i++ {
		go func() {
			for url := range jobs {
				results <- fetchURL(url)
			}
		}()
	}

	// ジョブ投入
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// 結果収集
	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		allResults = append(allResults, <-results)
	}

	return allResults
}

// 使用例: 最大10並列
pool := &WorkerPool{workerCount: 10}
results := pool.FetchURLs(urls)
```

### 3. セマフォパターン

並列度を制限する別の方法。

```go
import "golang.org/x/sync/semaphore"

func FetchURLsWithSemaphore(ctx context.Context, urls []string, maxConcurrent int64) []URLResponse {
	sem := semaphore.NewWeighted(maxConcurrent)
	results := make(chan URLResponse, len(urls))

	for _, url := range urls {
		// セマフォを取得（最大maxConcurrent個まで）
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}

		go func(u string) {
			defer sem.Release(1)
			results <- fetchURL(u)
		}(url)
	}

	// 全結果を収集
	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		allResults = append(allResults, <-results)
	}

	return allResults
}
```

### 4. エラーグループパターン

エラーハンドリングを簡潔に。

```go
import "golang.org/x/sync/errgroup"

func FetchURLsWithErrGroup(ctx context.Context, urls []string) ([]URLResponse, error) {
	g, ctx := errgroup.WithContext(ctx)
	results := make([]URLResponse, len(urls))

	for i, url := range urls {
		i, url := i, url // キャプチャ
		g.Go(func() error {
			result := fetchURL(ctx, url)
			if result.Error != "" {
				return fmt.Errorf("failed to fetch %s: %s", url, result.Error)
			}
			results[i] = result
			return nil
		})
	}

	// すべてのgoroutineの完了を待つ
	// 1つでもエラーがあればエラーを返す
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
```

## 計算量とパフォーマンス

### 時間計算量

| 方式 | 時間計算量 | 説明 |
| --- | --- | --- |
| 順次実行 | O(n × t) | n = URL数、t = 平均レスポンス時間 |
| 並列実行 | O(t_max) | t_max = 最も遅いリクエストの時間 |
| ワーカープール | O(⌈n/w⌉ × t) | w = ワーカー数 |

**具体例**:
```
10個のURL、各2秒の場合:

順次実行:     10 × 2秒 = 20秒
並列実行:     max(2秒) = 2秒  ← 10倍高速！
ワーカー(3):  ⌈10/3⌉ × 2秒 = 8秒
```

### メモリ使用量

| 項目 | メモリ使用量 | 説明 |
| --- | --- | --- |
| goroutine | 約2KB/個 | スタックサイズ（初期値） |
| チャネル | O(バッファサイズ) | バッファ付きチャネル |
| レスポンス | O(n × レスポンスサイズ) | 全結果を保持 |

**メモリ推定**:
```
100個のURL並列取得:
- goroutine: 100 × 2KB = 200KB
- チャネル: 100個 × 数バイト = 数KB
- レスポンス: 100 × 1KB = 100KB
合計: 約300KB（非常に軽量）
```

### ワーカープールの選択基準

```go
// 小規模（< 100リクエスト）
// → 無制限並列でOK
results := fetchAllParallel(urls)

// 中規模（100-1000リクエスト）
// → ワーカープールで制限
pool := &WorkerPool{workerCount: 50}
results := pool.FetchURLs(urls)

// 大規模（> 1000リクエスト）
// → セマフォ + バッチ処理
sem := semaphore.NewWeighted(100)
for _, batch := range splitIntoBatches(urls, 100) {
	processBatch(ctx, sem, batch)
}
```

## 使いどころ

### 向いている場面

1. **マイクロサービス集約**: 複数サービスからデータを並列取得
   ```go
   urls := []string{
       "http://user-service/api/user/123",
       "http://order-service/api/orders?user=123",
       "http://payment-service/api/payments?user=123",
   }
   results := fetchAllParallel(urls)
   // 3つのサービスを並列に呼び出して集約
   ```

2. **ヘルスチェック**: 複数エンドポイントの監視
   ```go
   endpoints := []string{
       "https://api1.example.com/health",
       "https://api2.example.com/health",
       "https://db.example.com/health",
   }
   results := fetchAllParallel(endpoints)
   // すべて200 OKか確認
   ```

3. **CDN/ミラー選択**: 最速サーバーを自動選択
   ```go
   mirrors := []string{
       "https://cdn1.example.com/file.zip",
       "https://cdn2.example.com/file.zip",
       "https://cdn3.example.com/file.zip",
   }
   fastest := fetchFastest(mirrors)
   // 最速ミラーからダウンロード
   ```

4. **価格比較**: 複数ECサイトの価格を並列取得
   ```go
   sites := []string{
       "https://store-a.com/api/product/123",
       "https://store-b.com/api/product/123",
       "https://store-c.com/api/product/123",
   }
   prices := fetchAllParallel(sites)
   // 最安値を見つける
   ```

### 向いていない場面

1. **順序依存**: 前のリクエスト結果が次に必要
2. **レート制限**: API呼び出し制限がある場合
3. **リソース制約**: メモリやネットワーク帯域が限られている

## 実用例

### 1. マイクロサービス集約API

```go
type AggregateHandler struct {
	userServiceURL    string
	orderServiceURL   string
	paymentServiceURL string
}

func (h *AggregateHandler) GetUserDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3つのサービスを並列呼び出し
	type result struct {
		name string
		data interface{}
		err  error
	}

	results := make(chan result, 3)

	// ユーザー情報
	go func() {
		data, err := h.fetchUserInfo(ctx, userID)
		results <- result{"user", data, err}
	}()

	// 注文履歴
	go func() {
		data, err := h.fetchOrders(ctx, userID)
		results <- result{"orders", data, err}
	}()

	// 支払い情報
	go func() {
		data, err := h.fetchPayments(ctx, userID)
		results <- result{"payments", data, err}
	}()

	// 結果を集約
	dashboard := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		res := <-results
		if res.err != nil {
			// エラーでも他のデータは返す（部分的成功）
			dashboard[res.name] = map[string]string{"error": res.err.Error()}
		} else {
			dashboard[res.name] = res.data
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}
```

### 2. リトライ付き並列フェッチ

```go
func FetchWithRetry(ctx context.Context, url string, maxRetries int) URLResponse {
	var lastError error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result := fetchURL(ctx, url)
		if result.Error == "" {
			return result
		}

		lastError = fmt.Errorf(result.Error)

		// 指数バックオフ
		if attempt < maxRetries {
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return URLResponse{URL: url, Error: ctx.Err().Error()}
			}
		}
	}

	return URLResponse{
		URL:   url,
		Error: fmt.Sprintf("failed after %d retries: %v", maxRetries, lastError),
	}
}

// 並列フェッチ + リトライ
func FetchAllWithRetry(urls []string, maxRetries int) []URLResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make(chan URLResponse, len(urls))

	for _, url := range urls {
		go func(u string) {
			results <- FetchWithRetry(ctx, u, maxRetries)
		}(url)
	}

	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		allResults = append(allResults, <-results)
	}

	return allResults
}
```

### 3. プログレス表示付き並列フェッチ

```go
type ProgressTracker struct {
	total     int
	completed int
	mu        sync.Mutex
}

func (p *ProgressTracker) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed++
	fmt.Printf("Progress: %d/%d (%.1f%%)\n",
		p.completed, p.total, float64(p.completed)/float64(p.total)*100)
}

func FetchWithProgress(urls []string) []URLResponse {
	tracker := &ProgressTracker{total: len(urls)}
	results := make(chan URLResponse, len(urls))

	for _, url := range urls {
		go func(u string) {
			result := fetchURL(context.Background(), u)
			results <- result
			tracker.Increment()
		}(url)
	}

	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		allResults = append(allResults, <-results)
	}

	return allResults
}
```

### 4. キャッシュ付き並列フェッチ

```go
type CachedFetcher struct {
	cache map[string]URLResponse
	mu    sync.RWMutex
}

func NewCachedFetcher() *CachedFetcher {
	return &CachedFetcher{
		cache: make(map[string]URLResponse),
	}
}

func (cf *CachedFetcher) FetchURLs(urls []string) []URLResponse {
	results := make(chan URLResponse, len(urls))
	fetchNeeded := []string{}

	// キャッシュチェック
	cf.mu.RLock()
	for _, url := range urls {
		if cached, exists := cf.cache[url]; exists {
			results <- cached
		} else {
			fetchNeeded = append(fetchNeeded, url)
		}
	}
	cf.mu.RUnlock()

	// キャッシュミスのURLを並列フェッチ
	for _, url := range fetchNeeded {
		go func(u string) {
			result := fetchURL(context.Background(), u)

			// キャッシュに保存
			cf.mu.Lock()
			cf.cache[u] = result
			cf.mu.Unlock()

			results <- result
		}(url)
	}

	// 全結果を収集
	allResults := make([]URLResponse, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		allResults = append(allResults, <-results)
	}

	return allResults
}
```

## まとめ

### メリット

- **高速化**: 順次実行と比べて劇的に高速（n倍）
- **効率的**: I/O待ち時間を有効活用
- **スケーラブル**: goroutineは軽量で大量起動可能
- **柔軟性**: 最速取得、全結果取得など用途に応じた実装
- **タイムアウト制御**: context.Contextで簡単に制御

### デメリット

- **複雑性**: 並行処理特有のバグ（競合状態など）
- **リソース消費**: 大量の並列リクエストはサーバー負荷に
- **デバッグ困難**: goroutineのデバッグは難しい
- **エラーハンドリング**: 複数のエラーを適切に処理する必要

### 使うべき時

- **I/Oバウンド**: ネットワークI/Oが支配的な処理
- **独立したリクエスト**: 各リクエストが互いに依存しない
- **レイテンシ重視**: 応答速度が重要
- **多数のエンドポイント**: 複数APIを集約

### ベストプラクティス

1. **タイムアウト設定**: 必ずcontext.WithTimeoutを使用
2. **並列度制限**: ワーカープールで制限（特に大規模時）
3. **エラーハンドリング**: 各goroutineのエラーを適切に処理
4. **リソース管理**: http.Response.Bodyは必ずClose
5. **テスト**: httptest.Serverでテスト容易に

### パターン選択

| ユースケース | パターン | 理由 |
| --- | --- | --- |
| 最速サーバー選択 | Fan-Out + 最速チャネル | 最初の成功を返す |
| 全データ収集 | Fan-Out / Fan-In | すべての結果が必要 |
| 大量リクエスト | ワーカープール | リソース保護 |
| エラー重視 | errgroup | 1つでも失敗したら中止 |

Goの並行処理機能を活用すれば、複数のHTTPリクエストを効率的に並列実行できる。goroutineとchannelを適切に組み合わせることで、わずか数十行のコードで強力な並列フェッチ機能を実装可能。マイクロサービスアーキテクチャやAPI集約など、モダンなWeb開発において必須のテクニック。
