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

// PresetAggregateHandler は事前定義された複数APIを並列で叩くハンドラ
func (h *AggregateAPIHandler) PresetAggregateHandler(w http.ResponseWriter, r *http.Request) {
	// 例: 天気、時刻、GitHub APIなど
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
