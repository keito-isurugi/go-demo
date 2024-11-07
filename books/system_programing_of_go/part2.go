package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func Part2() {
	byteArray := []byte("ASCII")
	fmt.Println(byteArray) // [65 83 67 73 73]

	str := string([]byte{65, 83, 67, 73, 73})
	fmt.Println(str) // ASCII


	// ファイルへの書き込み
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	file.Write([]byte("os.File example\n"))
	file.Close()

	// 画面出力(fmt.Print()と等価)
	os.Stdout.Write([]byte("os.File example\n"))

	// 書かれた内容を記憶しておく
	var buffer bytes.Buffer
	buffer.Write([]byte("bytes.Buffer  example1\n"))
	buffer.Write([]byte("bytes.Buffer  example2\n"))
	fmt.Println(buffer.String())

	// 書かれた内容を記憶しておく2
	var builder strings.Builder
	builder.Write([]byte("strings.Builder  example1\n"))
	builder.Write([]byte("strings.Builder  example2\n"))
	fmt.Println(builder.String())
}
