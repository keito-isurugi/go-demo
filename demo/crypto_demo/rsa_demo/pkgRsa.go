package rsademo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

// RSAによる暗号化
func encryptRSA(pub *rsa.PublicKey, message string) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte(message), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSAによる復号化
func decryptRSA(priv *rsa.PrivateKey, ciphertext string) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, decodedCiphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func ExecPkgRSA() {
	// RSA鍵ペアを生成
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("鍵ペアの生成に失敗しました: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// 平文
	message := "Hello, World!"

	// 暗号化
	ciphertext, err := encryptRSA(publicKey, message)
	if err != nil {
		log.Fatalf("暗号化に失敗しました: %v", err)
	}
	fmt.Println("暗号化されたデータ:", ciphertext)

	// 復号化
	decryptedMessage, err := decryptRSA(privateKey, ciphertext)
	if err != nil {
		log.Fatalf("復号化に失敗しました: %v", err)
	}
	fmt.Println("復号化されたデータ:", decryptedMessage)
}

// 出力
// 暗号化されたデータ: RAAYXaNj4lLeGpoBq3hjEmVQwVyQJdREK9jA8X56jwgGit6NB4ygvHtog7coZ1e72TPfkRDjbjUrL4gXKnVi6JjcfBwnBRIzxDnT5V4baceIfE4rlJFFv0g8Ut6vFLDPiKWnsmfqAu8Vmt2iejP6r8cADdwkK8HMs4/bH+VSzzDgabQJOeT7p4/i5xTynNZifrH2hYK7lwrh6FwojDMex/gb1XxQrtpwseZdYxZqB/YgYQNZY/DHnQYZ3QL3Ev6e68qPUvUqAsg7qWZ3SlyfEu7nWQhUslB6C6mAVohpNbQ+Ks5rjG1u8yE5Q7PJVDmDJutpXerL1dd0z4I1vb1F3A==
// 復号化されたデータ: Hello, RSA Encryption!