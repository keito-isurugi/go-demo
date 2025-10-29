package main

import (
	"context"
	"fmt"

	"github.com/keito-isurugi/go-demo/books"
	"github.com/keito-isurugi/go-demo/handler"
	"github.com/keito-isurugi/go-demo/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"net/http"
	"time"
)

func main() {
	// DB接続
	dsn := "host=db user=postgres password=postgres dbname=go_demo port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Redis接続
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

	// CacheHandlerの初期化
	cacheHandler := &handler.CacheHandler{
		DB:    dbConn,
		Redis: rdb,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!!!")
    })

	// 時関
	http.HandleFunc("/demo/time", handler.TimeDemoHandler)

	// アルゴリズム
	http.HandleFunc("/demo/algorithm", handler.AlgorithmDemoHandler)

	// 書籍
	http.HandleFunc("/demo/books", books.BooksDemoHandler)

	// キャッシュデモ
	// テストデータ(10,000件)を作成
	http.HandleFunc("/demo/cache/init", cacheHandler.InitTestDataHandler)
	// Redisキャッシュを使用してログ取得
	http.HandleFunc("/demo/cache/with", cacheHandler.CacheWithHandler)    
	// キャッシュなしでDBから直接ログ取得
	http.HandleFunc("/demo/cache/without", cacheHandler.CacheWithoutHandler) 
	// Redisキャッシュをクリア
	http.HandleFunc("/demo/cache/clear", cacheHandler.ClearCacheHandler)

	// レート制限付きAPI (1分間に10回まで)
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	http.Handle("/api/limited", rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Success! This endpoint is rate-limited to 10 requests per minute."}`))
	})))

	// セキュリティデモAPI
	securityHandler := &handler.SecurityDemoHandler{DB: dbConn}
	// 情報エンドポイント
	http.HandleFunc("/api/security", securityHandler.SecurityInfoHandler)

	// CORS デモ
	http.HandleFunc("/api/security/cors/vulnerable", securityHandler.CORSVulnerableHandler)
	http.HandleFunc("/api/security/cors/secure", securityHandler.CORSSecureHandler)

	// CSRF デモ
	http.HandleFunc("/api/security/csrf/token", securityHandler.CSRFTokenHandler)
	http.HandleFunc("/api/security/csrf/vulnerable", securityHandler.CSRFVulnerableHandler)
	http.HandleFunc("/api/security/csrf/secure", securityHandler.CSRFSecureHandler)

	// XSS デモ
	http.HandleFunc("/api/security/xss/vulnerable", securityHandler.XSSVulnerableHandler)
	http.HandleFunc("/api/security/xss/secure", securityHandler.XSSSecureHandler)

	// SQL Injection デモ
	http.HandleFunc("/api/security/sql-injection/vulnerable", securityHandler.SQLInjectionVulnerableHandler)
	http.HandleFunc("/api/security/sql-injection/secure", securityHandler.SQLInjectionSecureHandler)

	// 並列URL取得API
	parallelFetchHandler := &handler.ParallelFetchHandler{}
	// 最速レスポンスを返す
	http.HandleFunc("/api/fetch/fastest", parallelFetchHandler.FetchFastestHandler)
	// 全結果を返す
	http.HandleFunc("/api/fetch/all", parallelFetchHandler.FetchAllHandler)

	// 複数API並列実行・集約API
	aggregateHandler := &handler.AggregateAPIHandler{}
	// カスタムAPIを並列実行して集約
	http.HandleFunc("/api/aggregate", aggregateHandler.AggregateHandler)
	// プリセットAPIを並列実行して集約（デモ用）
	http.HandleFunc("/api/aggregate/preset", aggregateHandler.PresetAggregateHandler)

	fmt.Println("localhost:8080 server runnig ...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func db() (Todo, error) {
	dsn := "host=db user=postgres password=postgres dbname=go_demo port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return Todo{}, err
	}

	var todo Todo
	if result := db.Preload("Users").First(&todo); result.Error != nil {
		fmt.Println(result.Error)
		return Todo{}, err
	}

	return todo, nil
}

type Todo struct {
	ID        int `gorm:"primaryKey"`
	UserID		int
	User      User
	Title     string
	DoneFlag  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Category struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TodoCategory struct {
	ID         int `gorm:"primaryKey"`
	TodoID     int
	CategoryID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type Todos []Todo
