package main

import (
	"fmt"

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
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Bar!")
    })
	http.HandleFunc("/demo/time", handler.TimeDemoHandler)
	http.HandleFunc("/demo/algorithm", handler.AlgorithmDemoHandler)
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
