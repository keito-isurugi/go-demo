package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "log"
	// "net/http"
	"time"
	"github.com/keito-isurugi/go-demo/goodbaddev"
)

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

func main() {
	// http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	// 	todo, err := db()
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "Todo: %+v", todo)
	// })
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// switchを一箇所にした実装
	magicGouka := &goodbaddev.MagicMnager{}
	magicGouka.NewMagic(goodbaddev.MagicHellFire)
	magicGouka.GetName()
	magicGouka.CostMagicPoint()
	magicGouka.AttackPower()
	magicGouka.CostTechnicalPoint()
	fmt.Println("=================")
	magicShiden := &goodbaddev.MagicMnager{}
	magicShiden.NewMagic(goodbaddev.MagicShiden)
	magicShiden.GetName()
	magicShiden.CostMagicPoint()
	magicShiden.AttackPower()
	magicShiden.CostTechnicalPoint()
	fmt.Println("=================")

	// strategyパターンを使用した実装
	fire := &goodbaddev.Fire{}
	shiden := &goodbaddev.Shiden{}
	
	manager := goodbaddev.NewMagicManager(fire)
	fmt.Println("ファイヤの情報：")
	manager.Execute()
	fmt.Println()
	manager.SetStrategy(shiden)
	fmt.Println("紫電の情報：")
	manager.Execute()
	
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
