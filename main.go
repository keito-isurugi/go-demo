package main

import (
	"context"
	"fmt"

	"github.com/keito-isurugi/go-demo/books"
	"github.com/keito-isurugi/go-demo/handler"
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
