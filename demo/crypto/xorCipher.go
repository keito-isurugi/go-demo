package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// XOR暗号化・復号化を行う関数
// data: 暗号化または復号化する平文または暗号文
// key: 使用する鍵
// XOR演算を適用し、暗号化または復号化した結果を返す
func xorEncryptDecrypt(data, key string) string {
	dataBytes := []byte(data)              // 平文または暗号文をバイト列に変換
	keyBytes := []byte(key)                // 鍵をバイト列に変換
	output := make([]byte, len(dataBytes)) // 出力用のバイト列を準備

	keyLength := len(keyBytes) // 鍵の長さを取得

	// XOR演算を各バイトに対して適用
	for i := 0; i < len(dataBytes); i++ {
		output[i] = dataBytes[i] ^ keyBytes[i%keyLength] // i番目のデータバイトと鍵バイトをXOR
	}

	// 結果を文字列として返す
	return string(output)
}

// 文字列をバイナリ形式の文字列に変換
func toBinaryString(data string) string {
	binaryStr := ""
	for _, c := range data {
		binaryStr += fmt.Sprintf("%08b ", c)
	}
	return binaryStr
}

// 平文と同じ長さの鍵を生成する関数
func generateRandomKey(length int) (string, error) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func xorExec() {
	plaintext := "HelloWorld"
	key, err := generateRandomKey(len(plaintext)) // 平文と同じ長さの鍵を生成
	if err != nil {
		fmt.Println("ランダム鍵の生成に失敗しました:", err)
		return
	}

	fmt.Printf("平文: %s\n", plaintext)
	fmt.Printf("平文(バイナリ): %s\n", toBinaryString(plaintext))
	fmt.Printf("鍵(バイナリ)  : %s\n", toBinaryString(key))

	encrypted := xorEncryptDecrypt(plaintext, key)
	fmt.Printf("暗号文(バイナリ): %s\n", toBinaryString(encrypted))
	fmt.Printf("暗号文: %s\n", hex.EncodeToString([]byte(encrypted)))

	// 復号化は暗号化と同じ操作を行う（XORは可逆操作）
	decrypted := xorEncryptDecrypt(encrypted, key)
	fmt.Printf("復号文: %s\n", decrypted)
}

// 出力
// 平文　: HelloWorld
// 暗号文: 294fd36b7622b728ab71
// 復号文: HelloWorld
