package main

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"reflect"
	"strings"
	"testing"
)

func BenchmarkBatchProcessing(b *testing.B) {
	db, _ := dbConn()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fileOutPutTodosBatchedBench(db, "file_batched.txt")
		os.Remove("file_batched.txt")
	}
}

func fileOutPutTodosBatchedBench(db *gorm.DB, fileName string) error {
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

func BenchmarkStreamProcessing(b *testing.B) {
	db, _ := dbConn()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fileOutPutTodosWithStreamBench(db, "file_stream.txt")
		os.Remove("file_stream.txt")
	}
}

func fileOutPutTodosWithStreamBench(db *gorm.DB, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	rows, err := db.Table("todos").Rows()
	if err != nil {
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
			return err
		}
	}

	return nil
}
