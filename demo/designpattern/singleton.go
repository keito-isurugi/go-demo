package main

import (
	"fmt"
	"sync"
)

// Singleton構造体：シングルトンとして扱う対象の構造体
type Singleton struct {
	Value string
}

var (
	once     sync.Once     // sync.Onceは一度だけ処理を実行するための同期プリミティブ
	instance *Singleton    // 実際のシングルトンインスタンス
)

// GetInstanceはシングルトンインスタンスを返す関数
// 最初の呼び出し時に一度だけインスタンスが生成される
func getInstance() *Singleton {
	// 初回のみこの関数内の処理が実行される
	once.Do(func() {
		fmt.Println("Creating Singleton instance")
		instance = &Singleton{
			Value: "初期値",
		}
	})
	return instance
}

func singletonExec() {
	// 複数回呼び出しても同じインスタンスが返る
	s1 := getInstance()
	s2 := getInstance()

	fmt.Println("s1.Value:", s1.Value)
	s2.Value = "変更後の値"
	fmt.Println("s1.Value after s2 update:", s1.Value)

	// 同じインスタンスであることを確認
	fmt.Println("s1とs2は同一インスタンスか？", s1 == s2)
}

// 出力
// Creating Singleton instance
// s1.Value: 初期値
// s1.Value after s2 update: 変更後の値
// s1とs2は同一インスタンスか？ true