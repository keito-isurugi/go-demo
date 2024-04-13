package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Todo struct {
	ID        int `gorm:"primaryKey"`
	Title     string
	DoneFlag  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Todos []Todo

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		todo, err := db()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Todo: %+v", todo)
	})
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
	if result := db.First(&todo); result.Error != nil {
		fmt.Println(result.Error)
		return Todo{}, err
	}

	return todo, nil
}
