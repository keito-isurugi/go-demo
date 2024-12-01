package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func BenchmarkFileOutPutTodos(b *testing.B) {
	db, _ := dbConn()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FileOutPutTodos(db, "file_output.txt")
		os.Remove("file_output.txt")
	}
}

func FileOutPutTodos(db *gorm.DB, fileName string) error {
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

func BenchmarkFileOutPutTodosWithRefrect(b *testing.B) {
	db, _ := dbConn()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FileOutPutTodosWithRefrect(db, "file_output.txt")
		os.Remove("file_output.txt")
	}
}

func FileOutPutTodosWithRefrect(db *gorm.DB, fileName string) error {
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

	for _, todo := range todos {
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
			return err
		}
	}

	return nil
}
