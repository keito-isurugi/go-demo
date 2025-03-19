package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// INFO https://tex2e.github.io/rfc-translater/html/rfc7636.html#4-1--Client-Creates-a-Code-Verifier
const (
	verifierLength  = 128 // 安全性を考慮して最大値
	verifierCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
)

func GenerateCodeVerifier() (string, error) {
	verifier := make([]byte, verifierLength)

	for i := range verifier {
		b := make([]byte, 1)
		// 予測の難しいい乱数値がほしいのでmath/randではなくcrypto/randを使用
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}

		verifier[i] = verifierCharSet[b[0]%byte(len(verifierCharSet))]
	}

	return string(verifier), nil
}

func GenerateCodeChallenge(verifier string) (string, error) {
	if verifier == "" {
		return "", errors.New("code_verifier cannot be empty")
	}

	hash := sha256.Sum256([]byte(verifier))

	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return codeChallenge, nil
}

func execPKCE() {
	codeVerifier, _ := GenerateCodeVerifier()
	codeChallenge, _ := GenerateCodeChallenge(codeVerifier)

	fmt.Println("code_verifier: " + codeVerifier)
	fmt.Println("code_challenge: " + codeChallenge)
}
// 出力例
// code_verifier: LgLcHqsW3bSgP1p1Bc4VlrqnuuEiVluwyuECltiDAgGH25z4d6ZSZdi091Hk~DoVSVclVYf~d_phFdu_3TiS9mMu8eOCy-.273JX71x59YlfCKAjhF0H0RE53etRYlHM
// code_challenge: KPhLPYvrFhwtYPXg73QpdUcgpLYleVcwy-HxPrKu6Mc