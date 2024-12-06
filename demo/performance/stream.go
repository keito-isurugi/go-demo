package main

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync/atomic"
)

func Stream(db *gorm.DB) {
	fmt.Println("ストリーム処理")
	fileName := "todos.txt"
	FileOutPutTodosWithFprintf(db, fileName)
	os.Remove(fileName)
	FileOutPutTodosWithBufioFlushCount(db, fileName)
	os.Remove(fileName)
}

// fmt.Fprintfを使用
func FileOutPutTodosWithFprintf(db *gorm.DB, fileName string) error {
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

// bufio.Writerを使用
func FileOutPutTodosWithBufio(db *gorm.DB, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

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
		line := fmt.Sprintf("{%s},\n", strings.Join(fields, ", "))
		_, err = writer.Write([]byte(line))
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

// bufio.Writerを使用、書き込み回数を計測
func FileOutPutTodosWithBufioFlushCount(db *gorm.DB, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//writer := bufio.NewWriterSize(file, 1024)
	writer := bufio.NewWriter(file)
	countingWriter := NewCountingWriter(writer)

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
		line := fmt.Sprintf("{%s},\n", strings.Join(fields, ", "))
		_, err = countingWriter.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	err = countingWriter.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("Flush Count: %d\n", countingWriter.FlushCount())

	return nil
}
