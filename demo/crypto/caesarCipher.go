package main

import "fmt"

const (
	AlphabetSize = 26
	Shift        = 1
)

func caesarExec() {
	text := "Hello, world!"
	fmt.Printf("Original: %s\n", text)

	encrypted := caesarCipher(text, Shift)
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted := caesarCipher(encrypted, -Shift)
	fmt.Printf("Decrypted: %s\n", decrypted)
}

// シーザー暗号の実装
func caesarCipher(input string, shift int) string {
	// 「shift % AlphabetSize」・・・-25~25の範囲内の数値に変換
	// 「+ AlphabetSize」・・・-25~-1の値を0~25の範囲内の数値に変換。0~25だった数値は26~51となってしまう
	// 「% AlphabetSize」・・・26~51となってしまった数値を0~25の範囲内の数値に変換
	shift = (shift%AlphabetSize + AlphabetSize) % AlphabetSize // シフトを正の範囲に正規化

	var result []rune
	for _, char := range input {
		if char >= 'A' && char <= 'Z' { // 大文字の場合
			shifted := 'A' + (char-'A'+rune(shift))%AlphabetSize
			result = append(result, shifted)
		} else if char >= 'a' && char <= 'z' { // 小文字の場合
			shifted := 'a' + (char-'a'+rune(shift))%AlphabetSize
			result = append(result, shifted)
		} else {
			// アルファベット以外の文字はそのまま
			result = append(result, char)
		}
	}
	return string(result)
}
