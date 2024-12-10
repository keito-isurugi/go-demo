package main

import (
	"encoding/hex"
	"fmt"
)

func toBinaryString(data string) string {
	binaryStr := ""
	for _, c := range data {
		binaryStr += fmt.Sprintf("%08b ", c)
	}
	return binaryStr
}

func xorEncryptDecrypt(data, key string) string {
	dataBytes := []byte(data)
	keyBytes := []byte(key)
	output := make([]byte, len(dataBytes))

	keyLength := len(keyBytes)

	for i := 0; i < len(dataBytes); i++ {
		output[i] = dataBytes[i] ^ keyBytes[i%keyLength] // XOR演算を適用
	}

	return string(output)
}

func xorExec() {
	plaintext := "Hello"
	key := "aaaaa"

	fmt.Printf("平文        : %s\n", plaintext)
	fmt.Printf("平文 (バイナリ): %s\n", toBinaryString(plaintext))
	fmt.Printf("鍵 (バイナリ)  : %s\n", toBinaryString(key))

	encrypted := xorEncryptDecrypt(plaintext, key)
	fmt.Printf("暗号文(バイナリ): %s\n", toBinaryString(encrypted))
	fmt.Printf("暗号文         : %s\n", hex.EncodeToString([]byte(encrypted)))

	// 復号化は暗号化と同じ操作を行う（XORは可逆操作）
	decrypted := xorEncryptDecrypt(encrypted, key)
	fmt.Printf("復号文: %s\n", decrypted)
}
