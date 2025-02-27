package main

import (
	"fmt"
	"os"
	"sync"
)

func Part16() {
	part163()
}

func part161() {
	fmt.Printf("Page Size: %d\n", os.Getpagesize())
}

func part163() {
	// Poolを作成。Newで新規作成時のコードを実装
	var count int
	pool := sync.Pool{
		New: func() interface{} {
			count++
			return fmt.Sprintf("created: %d", count)
		},
	}

	// 追加した要素から受け取れる
	// プールが空だと新規作成
	pool.Put("manually added: 1")
	pool.Put("manually added: 2")
	fmt.Println(pool.Get())
	fmt.Println(pool.Get())
	fmt.Println(pool.Get())

}
