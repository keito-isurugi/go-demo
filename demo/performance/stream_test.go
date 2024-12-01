package main

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
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

type CountingWriter struct {
	wrapped    *bufio.Writer
	flushCount uint64
}

func NewCountingWriter(writer *bufio.Writer) *CountingWriter {
	return &CountingWriter{wrapped: writer}
}

func (cw *CountingWriter) Write(p []byte) (nn int, err error) {
	nn, err = cw.wrapped.Write(p)
	if err != nil {
		return nn, err
	}

	// データサイズチェックと、バッファフラッシュ追加
	if len(p) >= cw.wrapped.Available() {
		err := cw.Flush()
		if err != nil {
			return nn, err
		}
	}
	return nn, nil
}

func (cw *CountingWriter) Flush() error {
	atomic.AddUint64(&cw.flushCount, 1)
	return cw.wrapped.Flush()
}

func (cw *CountingWriter) FlushCount() uint64 {
	return atomic.LoadUint64(&cw.flushCount)
}

func BenchmarkStreamProcessing2(b *testing.B) {
	db, _ := dbConn()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		flushCount, err := fileOutPutTodosWithStreamBench2(db, "file_stream.txt")
		if err != nil {
			b.Fatalf("failed to process: %v", err)
		}

		fmt.Printf("Flush Count: %d\n", flushCount)
		//os.Remove("file_stream.txt")
	}
}

func fileOutPutTodosWithStreamBench2(db *gorm.DB, fileName string) (uint64, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	//writer := bufio.NewWriterSize(file, 1024)
	writer := bufio.NewWriter(file)
	countingWriter := NewCountingWriter(writer)

	rows, err := db.Table("todos").Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := db.ScanRows(rows, &todo); err != nil {
			return 0, err
		}

		fields := []string{
			fmt.Sprintf("ID: %v", todo.ID),
			fmt.Sprintf("Title: %v", todo.Title),
			fmt.Sprintf("Note: %v", todo.Note),
		}

		line := fmt.Sprintf("{%s},\n", strings.Join(fields, ", "))
		_, err = countingWriter.Write([]byte(line))
		if err != nil {
			return 0, err
		}
	}

	err = countingWriter.Flush()
	if err != nil {
		return 0, err
	}

	return countingWriter.FlushCount(), nil
}
