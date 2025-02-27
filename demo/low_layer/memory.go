package main

import (
	"fmt"
)

// グローバル変数（.dataセグメント）
var globalVar int = 42

// 未初期化グローバル変数（.bssセグメント）
var uninitVar int

func memoryLayout() {
	localVar := 10  // スタック
	ptr := new(int) // ヒープ
	*ptr = 20

	fmt.Printf("コード領域 (関数のアドレス): %p\n", main)
	fmt.Printf("データ領域 (.data): %p\n", &globalVar)
	fmt.Printf("データ領域 (.bss): %p\n", &uninitVar)
	fmt.Printf("スタック領域: %p\n", &localVar)
	fmt.Printf("ヒープ領域: %p\n", ptr)
}
