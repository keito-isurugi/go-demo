package main

import (
	"fmt"
	"math/rand"
)

// fakeOneWayHash は、入力文字列から疑似ハッシュ値を計算する自作ハッシュ関数です。
// セキュリティ目的には使えませんが、学習やテストには役立ちます。
func fakeOneWayHash(s string) int {
	hash := 0 // 初期値を0に変更

	// 文字列の各文字について処理
	for _, r := range s {
		// 単純な加算と乗算のみを使用
		hash = hash*31 + int(r)
	}

	// 正の32ビット整数に変換
	return hash & 0x7FFFFFFF
}

// verifyWeakCollision は、第二現像攻撃　によるfakeOneWayHashの弱衝突を検出する
func verifyWeakCollision(targetString string) {
	targetHash := fakeOneWayHash(targetString)
	fmt.Printf("対象文字列: %s\n", targetString)
	fmt.Printf("対象ハッシュ値: %d\n", targetHash)

	// 小文字、大文字、数字を含む文字セット
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	found := false

	// 3文字の全パターンを総当たり
	for _, c1 := range letters {
		for _, c2 := range letters {
			for _, c3 := range letters {
				s := string([]rune{c1, c2, c3})
				if fakeOneWayHash(s) == targetHash && s != targetString {
					fmt.Printf("衝突を発見: %s (ハッシュ値: %d)\n", s, fakeOneWayHash(s))
					found = true
				}
			}
		}
	}

	if !found {
		fmt.Println("衝突する文字列は見つかりませんでした")
	}
}

// verifyStrongCollision は、衝突攻撃によるfakeOneWayHashの強衝突を検出する
func verifyStrongCollision() {
	// ハッシュ値と文字列のマッピングを保持するマップ
	hashMap := make(map[int]string)
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	found := false

	// ランダムな文字列を生成して衝突を探す
	for i := 0; i < 1000000; i++ { // 最大100万回試行
		// 3文字のランダムな文字列を生成
		s := string([]rune{
			letters[rand.Intn(len(letters))],
			letters[rand.Intn(len(letters))],
			letters[rand.Intn(len(letters))],
		})

		hash := fakeOneWayHash(s)

		// 同じハッシュ値を持つ文字列が既に存在するかチェック
		if existing, exists := hashMap[hash]; exists && existing != s {
			fmt.Printf("強衝突を発見:\n")
			fmt.Printf("文字列1: %s (ハッシュ値: %d)\n", existing, hash)
			fmt.Printf("文字列2: %s (ハッシュ値: %d)\n", s, hash)
			found = true
			break
		}

		// 現在の文字列とハッシュ値をマップに保存
		hashMap[hash] = s
	}

	if !found {
		fmt.Println("強衝突は見つかりませんでした")
	}
}
