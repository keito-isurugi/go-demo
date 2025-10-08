package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

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
