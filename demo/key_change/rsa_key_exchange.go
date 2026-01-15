package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// =============================================================================
// RSA鍵交換 (RSA Key Transport)
// =============================================================================
//
// 【概要】
// RSA鍵交換は、クライアントが共通鍵（プリマスターシークレット）を生成し、
// サーバーの公開鍵で暗号化して送信する方式です。
// TLS 1.2以前で広く使われていましたが、TLS 1.3では廃止されました。
//
// 【仕組み】
// 1. サーバーがRSA鍵ペア（公開鍵・秘密鍵）を生成
// 2. サーバーが公開鍵をクライアントに送信（証明書に含まれる）
// 3. クライアントがプリマスターシークレット（48バイトの乱数）を生成
// 4. クライアントがサーバーの公開鍵でプリマスターシークレットを暗号化
// 5. サーバーが秘密鍵で復号し、共通鍵を取得
//
// 【セキュリティ特性】
// - 前方秘匿性（Forward Secrecy）がない
//   → サーバーの秘密鍵が漏洩すると、過去の通信も全て復号される
// - 鍵サイズによる強度の違い
//   → 2048ビット: 現在の最低推奨サイズ
//   → 3072ビット: 2030年以降も安全とされる
//   → 4096ビット: 長期的なセキュリティ向け
//
// =============================================================================

// RSAKeyExchange はRSA鍵交換をシミュレートします
type RSAKeyExchange struct {
	ServerPrivateKey *rsa.PrivateKey
	ServerPublicKey  *rsa.PublicKey
	KeyBits          int
}

// NewRSAKeyExchange は指定されたビット数でRSA鍵交換を初期化します
func NewRSAKeyExchange(bits int) (*RSAKeyExchange, error) {
	// サーバー側: RSA鍵ペアを生成
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("RSA鍵生成に失敗: %w", err)
	}

	return &RSAKeyExchange{
		ServerPrivateKey: privateKey,
		ServerPublicKey:  &privateKey.PublicKey,
		KeyBits:          bits,
	}, nil
}

// ClientGenerateAndEncryptSecret はクライアント側でプリマスターシークレットを
// 生成し、サーバーの公開鍵で暗号化します
func (r *RSAKeyExchange) ClientGenerateAndEncryptSecret() ([]byte, []byte, error) {
	// プリマスターシークレット（48バイト）を生成
	// TLSでは最初の2バイトがプロトコルバージョン、残り46バイトが乱数
	preMasterSecret := make([]byte, 48)
	if _, err := rand.Read(preMasterSecret); err != nil {
		return nil, nil, fmt.Errorf("プリマスターシークレット生成に失敗: %w", err)
	}

	// サーバーの公開鍵でプリマスターシークレットを暗号化
	// OAEP（Optimal Asymmetric Encryption Padding）を使用
	encryptedSecret, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		r.ServerPublicKey,
		preMasterSecret,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("暗号化に失敗: %w", err)
	}

	return preMasterSecret, encryptedSecret, nil
}

// ServerDecryptSecret はサーバー側で暗号化されたプリマスターシークレットを復号します
func (r *RSAKeyExchange) ServerDecryptSecret(encryptedSecret []byte) ([]byte, error) {
	// サーバーの秘密鍵で復号
	decryptedSecret, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		r.ServerPrivateKey,
		encryptedSecret,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("復号に失敗: %w", err)
	}

	return decryptedSecret, nil
}

// DemonstrateRSAKeyExchange はRSA鍵交換の完全なフローを示します
func DemonstrateRSAKeyExchange(bits int) error {
	fmt.Printf("\n%s\n", "==========================================================")
	fmt.Printf("RSA鍵交換 (%dビット)\n", bits)
	fmt.Printf("%s\n\n", "==========================================================")

	// 1. サーバー: 鍵ペアを生成
	fmt.Println("【ステップ1】サーバーがRSA鍵ペアを生成...")
	kex, err := NewRSAKeyExchange(bits)
	if err != nil {
		return err
	}
	fmt.Printf("  公開鍵サイズ: %d ビット\n", kex.ServerPublicKey.N.BitLen())
	fmt.Printf("  公開指数(e): %d\n", kex.ServerPublicKey.E)

	// 2. クライアント: プリマスターシークレットを生成・暗号化
	fmt.Println("\n【ステップ2】クライアントがプリマスターシークレットを生成・暗号化...")
	preMasterSecret, encryptedSecret, err := kex.ClientGenerateAndEncryptSecret()
	if err != nil {
		return err
	}
	fmt.Printf("  プリマスターシークレット (最初の16バイト): %x...\n", preMasterSecret[:16])
	fmt.Printf("  暗号化後のサイズ: %d バイト\n", len(encryptedSecret))

	// 3. サーバー: 復号
	fmt.Println("\n【ステップ3】サーバーが秘密鍵で復号...")
	decryptedSecret, err := kex.ServerDecryptSecret(encryptedSecret)
	if err != nil {
		return err
	}
	fmt.Printf("  復号されたシークレット (最初の16バイト): %x...\n", decryptedSecret[:16])

	// 4. 検証
	fmt.Println("\n【検証】共有シークレットが一致するか確認...")
	if string(preMasterSecret) == string(decryptedSecret) {
		fmt.Println("  ✓ 成功: クライアントとサーバーで同じシークレットを共有")
	} else {
		fmt.Println("  ✗ 失敗: シークレットが一致しません")
	}

	return nil
}
