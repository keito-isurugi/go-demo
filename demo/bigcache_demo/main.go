package main

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/allegro/bigcache/v3"
)

func main() {
	// bigcacheの設定
	cfg := bigcache.DefaultConfig(10 * time.Minute)
	bigCa, err := bigcache.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	// bigcacheの操作を抽象化したインスタンスを生成
	cache, err := NewCache(bigCa)
	if err != nil {
		log.Fatal(err)
	}

	// bigint型のデータをキャッシュに保存
	var bigIntValue big.Int
	bigIntValue.SetString("1234567890", 10)

	cache.Set("test-key", &bigIntValue)
	if err != nil {
		log.Fatal(err)
	}

	// キャッシュからデータを取得
	var result big.Int
	err = cache.Get("test-key", &result)
	if err != nil {
		log.Fatal(err)
	}
	println("result:", result.String())

	// キャッシュからデータを削除
	err = cache.Delete("test-key")
	if err != nil {
		log.Fatal(err)
	}

	// Getでキャッシュが空かどうかを確認
	var result2 big.Int
	err = cache.Get("test-key", &result2)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			println(err.Error())
		} else {
			println(err.Error())
		}
	}
	println("result2:", result2.String())
}
// 出力
// result: 1234567890
// Entry not found
// result2: 0
