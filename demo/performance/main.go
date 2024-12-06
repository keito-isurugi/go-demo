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
	// err = FileOutPutTodosWithRefrect1(db, "hoge.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//registerDummyTodos(db)
	Stream(db)
}

func FileOutPutTodosWithRefrect1(db *gorm.DB, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var todos []Todo
	result := db.Find(&todos)
	if result.Error != nil {
		return result.Error
	}

	fmt.Println(todos)
	for _, todo := range todos {
		fields := []string{
			fmt.Sprintf("ID: %v", todo.ID),
			fmt.Sprintf("Title: %v", todo.Title),
			fmt.Sprintf("Note: %v", todo.Note),
		}
		_, err := fmt.Fprintf(file, "{%s},\n", strings.Join(fields, ", "))
		if err != nil {
			return err
		}
	}

	return nil
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
	UserID    int
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

	count := 10000
	for i := 0; i < count; i++ {
		t := Todo{
			UserID: 1,
			Title:  fmt.Sprintf("ダミーTODO_%v", i+1),
			Note:   fmt.Sprintf("ダミーで登録したTODO_%vです", i+1),
		}
		db.Create(&t)
	}
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
