package main

import (
	"fmt"

	"github.com/keito-isurugi/go-demo/books"
	"github.com/keito-isurugi/go-demo/handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!!!")
    })

	// 時関
	http.HandleFunc("/demo/time", handler.TimeDemoHandler)

	// アルゴリズム
	http.HandleFunc("/demo/algorithm", handler.AlgorithmDemoHandler)

	// 書籍
	http.HandleFunc("/demo/books", books.BooksDemoHandler)
	
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
