package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// ParallelFetchHandler は複数URLを並列でGETし、最速レスポンスを返す
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
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	if len(req.URLs) == 0 {
		response := map[string]interface{}{
			"error": "No URLs provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
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
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}

// fetchURLsParallel は複数のURLを並列で取得し、最速のレスポンスを返す
func (h *ParallelFetchHandler) fetchURLsParallel(urls []string, timeoutSec int) FastestResponseResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// 最速レスポンス用のチャネル
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
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	if len(req.URLs) == 0 {
		response := map[string]interface{}{
			"error": "No URLs provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
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
		"results":       allResults,
		"total_urls":    len(req.URLs),
		"total_duration": time.Since(startTime).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
