package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LogRecord struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
}

type CacheHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

type PerformanceResult struct {
	Source     string      `json:"source"`
	Count      int         `json:"count"`
	DurationMs int64       `json:"duration_ms"`
	Records    []LogRecord `json:"records,omitempty"`
}

const (
	cacheKey       = "logs:latest"
	cacheTTL       = 5 * time.Minute
	testDataCount  = 100000  // テストデータ10万件
	fetchLimit     = 100     // 取得件数（実用的なサイズ）
)

// CacheWithHandler - キャッシュありのAPI（最新100件のログをキャッシュ）
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

// CacheWithoutHandler - キャッシュなしのAPI（毎回DBから最新100件を取得）
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

// ClearCacheHandler - キャッシュクリア用API
func (h *CacheHandler) ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if err := h.Redis.Del(ctx, cacheKey).Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to clear cache: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Cache cleared successfully",
	})
}

// InitTestDataHandler - テストデータ作成用API
func (h *CacheHandler) InitTestDataHandler(w http.ResponseWriter, r *http.Request) {
	// logsテーブルが存在しない場合は作成
	if !h.DB.Migrator().HasTable(&LogRecord{}) {
		if err := h.DB.AutoMigrate(&LogRecord{}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create table: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// 既存データをクリア
	h.DB.Exec("TRUNCATE TABLE log_records")

	// テストデータを挿入
	start := time.Now()
	logs := make([]LogRecord, testDataCount)
	levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}

	for i := 0; i < testDataCount; i++ {
		logs[i] = LogRecord{
			Message:   fmt.Sprintf("Log message #%d", i+1),
			Level:     levels[i%len(levels)],
			Timestamp: time.Now().Add(-time.Duration(i) * time.Second),
		}
	}

	// バッチインサート
	if err := h.DB.CreateInBatches(logs, 1000).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert test data: %v", err), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Test data created successfully",
		"record_count": testDataCount,
		"duration_ms":  duration.Milliseconds(),
	})
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
