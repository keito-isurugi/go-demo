package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 実行時間比較
	//for i := 0; i < 10; i++ {
	//	fmt.Printf("-----%v回目-----\n", i+1)
	//	PerformanceHnadler()
	//}

	db, err := dbConn()
	if err != nil {
		fmt.Println(err)
	}

	err = fileOutPutTodos(db, "file.txt")
	if err != nil {
		fmt.Println(err)
	}
	// registerDummyTodos(db)
}

func dbConn() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=go_demo port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

type Todo struct {
	ID        int `gorm:"primaryKey"`
	Title     string
	Note      string
	DoneFlag  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func getTodo(db *gorm.DB, id int) (Todo, error) {
	var todo Todo
	if result := db.Table("todos").Where("id", id).First(&todo); result.Error != nil {
		return todo, result.Error
	}

	return todo, nil
}
func listTodos(db *gorm.DB) ([]Todo, error) {
	var todos []Todo
	if result := db.Table("todos").Find(&todos); result.Error != nil {
		return todos, result.Error
	}
	return todos, nil
}

func registerDummyTodos(db *gorm.DB) {
	if err := db.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE").Error; err != nil {
		fmt.Errorf("failed to truncate todos table: %v", err)
	}

	count := 100000
	for i := 0; i <= count; i++ {
		t := Todo{
			Title: fmt.Sprintf("ダミーTODO_%v", i+1),
			Note:  fmt.Sprintf("ダミーで登録したTODO_%vです", i+1),
		}
		db.Create(&t)
	}
}

func fileOutPutTodos(db *gorm.DB, fileName string) error {
	todos, err := listTodos(db)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for _, todo := range todos {
		val := reflect.ValueOf(todo)
		typ := val.Type()

		var fields []string
		for i := 0; i < typ.NumField(); i++ {
			key := typ.Field(i).Name
			value := fmt.Sprintf("%v", val.Field(i).Interface())
			fields = append(fields, fmt.Sprintf("%v: %v", key, value))
		}
		// カンマで区切り、波括弧で囲んだフィールドのリストをテキストファイルに書き込む
		_, err := fmt.Fprintf(file, "{%s},\n", strings.Join(fields, ", "))
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("file output done.")
	return nil
}

func fileOutPutTodosWithStream(db *gorm.DB, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	rows, err := db.Table("todos").Rows()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		db.ScanRows(rows, &todo)
		val := reflect.ValueOf(todo)
		typ := val.Type()

		var fields []string
		for i := 0; i < typ.NumField(); i++ {
			key := typ.Field(i).Name
			value := fmt.Sprintf("%v", val.Field(i).Interface())
			fields = append(fields, fmt.Sprintf("%v: %v", key, value))
		}
		_, err := fmt.Fprintf(file, "{%s},\n", strings.Join(fields, ", "))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("file output done.")
	return nil
}
