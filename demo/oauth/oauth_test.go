package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	length := 10
	charSet := "ABCDEF"

	// 正常系: 指定した文字数と文字列で値が生成されること
	o := NewOauth(rand.Read)

	result, err := o.GenerateRandomString(length, charSet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 生成された文字列の文字数をチェック
	if len(result) != length {
		t.Errorf("expected length %d, got %d", length, len(result))
	}

	// 生成された文字列の各文字が、charSetに存在するかチェック
	for _, char := range result {
		if !strings.ContainsRune(charSet, char) {
			t.Errorf("generated character %c is not in the allowed charset %s", char, charSet)
		}
	}

	// 異常系: ReadRandomBytes がエラーを返す場合
	mockReadRandomBytes := func(_ []byte) (int, error) {
		return 0, errors.New("mocked error")
	}
	mockOauth := NewOauth(mockReadRandomBytes)

	_, err = mockOauth.GenerateRandomString(length, charSet)
	if err == nil {
		t.Error("expected error, but got nil")
	} else if err.Error() != "mocked error" {
		t.Errorf("expected mocked error, got %v", err)
	}
}

func TestGenerateState(t *testing.T) {
	// 指定した文字数と文字列でstateが生成されること
	o := NewOauth(rand.Read)

	state, err := o.GenerateState()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 生成されたstateの文字数チェック
	if len(state) != stateLength {
		t.Errorf("expected state length of 16, got %d", len(state))
	}

	// 生成されたstateの各文字が、charSetに存在するかチェック
	for _, char := range state {
		if !containsRune(charSet, char) {
			t.Errorf("character %c is not in the allowed charset", char)
		}
	}
}

func TestGenerateCodeVerifier(t *testing.T) {
	// 指定した文字数と文字列でverifierが生成されること
	o := NewOauth(rand.Read)

	verifier, err := o.GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 生成されたverifierの文字数チェック
	if len(verifier) != verifierLength {
		t.Errorf("expected verifier length of 128, got %d", len(verifier))
	}

	// 生成されたverifierの各文字が、charSetに存在するかチェック
	for _, char := range verifier {
		if !containsRune(charSet, char) {
			t.Errorf("character %c is not in the allowed charset", char)
		}
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	o := NewOauth(rand.Read)

	// 正常系: 生成されたcodeChallengeがexpectedChallengeと等しいこと
	verifier := "test-code-verifier"
	expectedHash := sha256.Sum256([]byte(verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(expectedHash[:])

	challenge, err := o.GenerateCodeChallenge(verifier)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if challenge != expectedChallenge {
		t.Errorf("expected code challenge %s, got %s", expectedChallenge, challenge)
	}

	// 異常系: verifier が空の場合
	_, err = o.GenerateCodeChallenge("")
	if err == nil {
		t.Error("expected error for empty code_verifier, but got nil")
	}

	expectedErrMsg := "code_verifier cannot be empty"
	if err.Error() != expectedErrMsg {
		t.Errorf("unexpected error message: %v", err)
	}
}

func containsRune(charSet string, char rune) bool {
	for _, c := range charSet {
		if c == char {
			return true
		}
	}
	return false
}
