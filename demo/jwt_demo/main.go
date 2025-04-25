package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// カスタムクレームを定義
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// トークンを生成する関数
func generateToken(secretKey []byte, userID string) (string, error) {
	// クレームの作成
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24時間後に期限切れ
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jwt-demo",
		},
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名を付与
	return token.SignedString(secretKey)
}

// トークンを検証する関数
func verifyToken(tokenString string, secretKey []byte) (*CustomClaims, error) {
	// トークンの検証
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("トークンの検証に失敗しました: %v", err)
	}

	if claims, ok := parsedToken.Claims.(*CustomClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("トークンが無効です")
}

func main() {
	// シークレットキー
	secretKey := []byte("test-secret-key")
	userID := "12345"

	// トークンの生成
	tokenString, err := generateToken(secretKey, userID)
	if err != nil {
		fmt.Printf("トークンの生成に失敗しました: %v\n", err)
		return
	}

	fmt.Printf("生成されたJWTトークン: %s\n", tokenString)

	// トークンの検証
	claims, err := verifyToken(tokenString, secretKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("トークンの検証に成功しました。UserID: %s\n", claims.UserID)
}
