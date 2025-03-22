package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

type StringWithHash struct {
	Count  int
	String string
	Hash   int
}

// breakthroughHash: simpleHashの衝突をランダムに探し、衝突した文字列とそのハッシュ値を返す
func breakthroughHash(inputString string) []StringWithHash {
	// simpleHashのハッシュ値を計算
	hashValue := simpleHash(inputString)

	result := make([]StringWithHash, 0)
	num := 17000
	fmt.Println("探索回数:", convert(num))

	// 衝突を探す
	for i := 0; i < num; i++ {
		// ランダムな文字列を生成
		randomString, _ := generateRandomString(uint32(len(inputString)))
		// 生成した文字列のハッシュ値を計算
		if simpleHash(randomString) == hashValue {
			result = append(result, StringWithHash{Count: i, String: randomString, Hash: hashValue})
		}
	}
	return result
}

// generateRandomString: 指定された長さのランダムな文字列を生成
func generateRandomString(digit uint32) (string, error) {
	// const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const letters = "abcdefghijklmnopqrstuvwxyz"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func simpleHash(inputString string) int {
	hashValue := 0
	// 文字列の各文字のASCIIコードを足し合わせる
	for _, char := range inputString {
		hashValue += int(char)
	}
	return hashValue
}

func execHash() {
	input := "hog"
	fmt.Printf("Input: %s\n", input)
	sh := simpleHash(input)
	fmt.Printf("Hash: %d\n", sh)

	fmt.Printf("文字数：%d\n", len(input))
	patern := 26
	for i := 1; i < len(input); i++ {
		patern *= 26
	}
	fmt.Printf("パターン数：%s\n", convert(patern))

	// 衝突を突破する文字列とそのハッシュ値を生成
	stringWithHash := breakthroughHash(input)

	hoge := make([]StringWithHash, 0)
	for _, sh := range stringWithHash {
		if input == sh.String {
			hoge = append(hoge, sh)
		}
	}

	if len(hoge) > 0 {
		fmt.Println("Collision found.")
		fmt.Println("見つかった衝突数:", len(hoge))
		for _, sh := range hoge {
			fmt.Printf("Count: %d, String: %s, Hash: %d\n", sh.Count, sh.String, sh.Hash)
		}
		return
	}
	fmt.Println("No collision found.")
}

func convert(integer int) string {
	arr := strings.Split(fmt.Sprintf("%d", integer), "")
	cnt := len(arr) - 1
	res := ""
	i2 := 0
	for i := cnt; i >= 0; i-- {
		if i2 > 2 && i2%3 == 0 {
			res = fmt.Sprintf(",%s", res)
		}
		res = fmt.Sprintf("%s%s", arr[i], res)
		i2++
	}
	return res
}
