package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// encryptAES 関数は、指定された平文データを AES-GCM を使用して暗号化し、結果を Base64 エンコードした文字列として返します。
//
// Parameters:
// - plainText: 暗号化するデータ (平文)。
// - key: 暗号化に使用する秘密鍵。長さは 16, 24, 32 バイト（AES-128, AES-192, AES-256 に対応）である必要があります。
//
// Returns:
// - 暗号化されたデータを Base64 エンコードした文字列。
// - エラーが発生した場合にはエラーメッセージ。
func encryptAES(plainText, key string) (string, error) {
	// AES 暗号化モジュールを作成（指定された鍵を使用）。
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// 鍵が無効な長さの場合、エラーを返す。
		return "", err
	}

	// GCM (Galois/Counter Mode) インスタンスを作成。
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		// GCM がサポートされていない、もしくは初期化に失敗した場合にエラーを返す。
		return "", err
	}

	// 暗号化のために Nonce（初期化ベクトル、IV）を生成。Nonce のサイズは GCM 仕様に依存。
	nonce := make([]byte, gcm.NonceSize())
	// 暗号学的に安全な乱数生成器を使用して Nonce を埋める。
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		// Nonce 生成に失敗した場合にエラーを返す。
		return "", err
	}

	// 平文（plainText）を暗号化し、Nonce をプレフィックスとして付与。
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	// 暗号文を Base64 形式の文字列に変換して返す。
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// decryptAES 関数は、AES-GCM を使用して暗号文を復号化し、平文データを返します。
//
// Parameters:
// - cipherText: Base64 エンコードされた暗号文。
// - key: 復号化に使用する秘密鍵。暗号化時に使用したものと同じ鍵を指定する必要があります。
//
// Returns:
// - 復号化された平文データ。
// - エラーが発生した場合にはエラーメッセージ。
func decryptAES(cipherText, key string) (string, error) {
	// AES 復号化モジュールを作成（指定された鍵を使用）。
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// 鍵が無効な長さの場合、エラーを返す。
		return "", err
	}

	// GCM (Galois/Counter Mode) インスタンスを作成。
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		// GCM がサポートされていない、もしくは初期化に失敗した場合にエラーを返す。
		return "", err
	}

	// Base64 形式の暗号文をデコードしてバイトデータに変換する。
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		// Base64 デコードに失敗した場合、エラーを返す。
		return "", err
	}

	// 暗号文の先頭に付加された Nonce を分離。Nonce サイズを確認。
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		// データが Nonce サイズより小さい場合は不正な暗号文としてエラーを返す。
		return "", errors.New("不正な暗号文です")
	}

	// Nonce（初期化ベクトル）と暗号化されたデータそれぞれに分割。
	nonce, cipherTextBytes := data[:nonceSize], data[nonceSize:]

	// 暗号化データを復号化して平文を取得。
	plainText, err := gcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		// 復号化に失敗した場合にエラーを返す。
		return "", err
	}

	// 復号化された平文を文字列として返す。
	return string(plainText), nil
}

func execAES() {
	// keyを32バイトに設定
	aesKey := "aaaaaaaaaabbbbbbbbbbccccccccccdd" 
	plainText := "Hello, World!"

	cipherText, err := encryptAES(plainText, aesKey)
	if err != nil {
		panic(err)
	}

	decryptText, err := decryptAES(cipherText, aesKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("plainText: " + plainText)
	fmt.Println("cipherText: " + cipherText)
	fmt.Println("decryptText: " + decryptText)
}
