package main

import "fmt"

// 固定長配列
func FixedArray() {
	// goは配列のサイズも含め型
	// 要素を増やしたり減らしたりはできない

	// 宣言(初期化) [サイズ]型
	var arr1 [4]int
	fmt.Println(arr1)
	
	// 要素を変更
	arr1[0] = 123
	arr1[3] = 789
	fmt.Println(arr1)
	// 参照
	fmt.Println(arr1[3])
	
	// 初期値をセットして宣言
	arr2 := [2]int{1, 2}
	fmt.Println(arr2)

	// サイズ
	fmt.Println(len(arr2))
}

// 可変長配列(スライス)
func VariableArray() {
	// 要素の追加や削除をできる配列
	// 基本的にこっちを使う
	
	// 宣言(初期化) []型
	var arr1 []int
	fmt.Println(arr1)

	// 要素を追加
	arr1 = append(arr1, 1)
	fmt.Println(arr1)
	// 参照
	fmt.Println(arr1[0])

	// makeを使用しサイズ指定で初期値0のスライス
	arr2 := make([]int, 3)
	fmt.Println(arr2)
	fmt.Println(arr2[2])
	fmt.Println(len(arr2))

	// [0, 2, 3, .... , 10]のスライスを作成
	var arr3 []int
	for i := 0; i <= 10; i++ {
		arr3 = append(arr3, i)
	}
	fmt.Println(arr3)
	// 範囲指定して取得
	// 左のインデックスは含まれる、右のインデックスは含まれない
	// インデックス先頭~4まで取得
	fmt.Println(arr3[:5])
	// インデックス5~末尾まで取得
	fmt.Println(arr3[5:])
	// インデックス2~7まで取得
	fmt.Println(arr3[2:8])

	// 要素を削除(インデックス3の要素を削除)
	arr3 = append(arr3[:3], arr3[4:]...)
	fmt.Println(arr3)
}

// マップ
func Map() {
	// 宣言(初期化) 長さ0のnilのマップが作成される
	var map1 map[int]int
	fmt.Println(map1)

	// 代入
	map1[1] = 100
	// 参照
	fmt.Println(map1[1])
}
