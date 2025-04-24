package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// グローバル変数として巨大なスライスへの参照を保持
var leakedData = make([][]byte, 0)

func memoryLeak() {
	for {
		// 1MBのデータを作成
		data := make([]byte, 1024*1024)

		// グローバル変数にデータの参照を保持し続ける
		leakedData = append(leakedData, data)

		// 少し待つ
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// pprof用にHTTPサーバを起動（:6060でアクセス可能）
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// メモリリークを発生させる関数を呼び出し
	memoryLeak()
}