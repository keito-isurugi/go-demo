package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func Part6() {
	part662()
}

func part662() {
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	current := 0
	var conn net.Conn = nil
	// リトライ用にループで全体を囲う
	for {
		var err error
		// まだコネクションを張ってない / エラーでリトライ
		if conn == nil {
			// Dialから行ってconnを初期化
			conn, err = net.Dial("tcp", "localhost:8888")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Access: %d\n", current)
		}

		// POSTで文字列を送るリクエストを作成
		request, err := http.NewRequest(
			"POST", 
			"http://localhost:8888", 
			strings.NewReader(sendMessages[current]),
		)
		if err != nil {
			panic(err)
		}
		err = request.Write(conn)
		if err != nil {
			panic(err)
		}

		// サーバーから読み込む。タイムアウトはここでエラーになるのでリトライ
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}
		// 結果を表示
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// 全部送信完了していれば終了
		current++
		if current == len(sendMessages) {
			break
		}
	}
	conn.Close()
}

func part661() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("server is running at localhost:8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// Accept後のソケットで何度も応答を返すためにループ
			for {
				// タイムアウトを設定
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				// リクエストを読み込む
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// タイムアウトもしくはソケットクローズ時は終了
					// それ以外はエラーにする
					neterr, ok := err.(net.Error) // ダウンキャスト
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						break
					}
					panic(err)
				}

				// リクエストを表示
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))
				content := "Hello World!\n"
				
				//レスポンスを書き込む
				// HTTP/1.1かつ、ContentLengthの設定が必要
				response := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 1,
					ContentLength: int64(len(content)),
					Body: io.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
		}()
	}
}

func part652() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("GET", "http://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	request.Write(conn)
	response, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}


func part651() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("server is running at localhost:8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// リクエストを読み込む
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(dump))
			//レスポンスを書き込む
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body: io.NopCloser(strings.NewReader("Hello World!\n")),
			}
			response.Write(conn)
			conn.Close()
		}()
	}
}