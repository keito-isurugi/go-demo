package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func Part2_2() {
	// conn, err := net.Dial("tcp", "example.com:80")
	// if err != nil {
	// 	panic(err)
	// }
	// io.WriteString(conn, "GET / HTTP/1.0\r\nHost: example.com\r\n\r\n")
	// io.Copy(os.Stdout, conn)

	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)

	// file, err := os.Create("multiwriter.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// writer := io.MultiWriter(file, os.Stdout)
	// io.WriteString(writer, "io.MultiWriter example\n")

	// file, err := os.Create("test.txt.zip")
	// if err != nil {
	// 	panic(err)
	// }
	// writer := gzip.NewWriter(file)
	// writer.Header.Name = "test.txt"
	// io.WriteString(writer, "gzip.Writer example\n")
	// writer.Close()

	// buffer := bufio.NewWriter(os.Stdout)
	// buffer.WriteString("bufio.Writer ")
	// buffer.Flush()
	// buffer.WriteString("example\n")
	// buffer.Flush()

	// fmt.Fprintf(os.Stdout, "Write with os.Stdout at %v", time.Now())

	// encoder := json.NewEncoder(os.Stdout)
	// encoder.SetIndent("", "   ")
	// encoder.Encode(map[string]string{
	// 	"example": "encoding/'json",
	// 	"hello":   "world",
	// })

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-TEST", "able add header")
	req.Write(os.Stdout)
}

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "http.ResponseWriter sample")
}
