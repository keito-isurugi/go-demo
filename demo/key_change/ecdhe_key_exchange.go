package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

// =============================================================================
// ECDHE鍵交換 (Elliptic Curve Diffie-Hellman Ephemeral)
// =============================================================================
//
// 【概要】
// ECDHEは、楕円曲線暗号（ECC）を用いたDiffie-Hellman鍵交換です。
// 従来のDHEと同等のセキュリティを、より小さな鍵サイズで実現できます。
// TLS 1.3ではECDHEが推奨される鍵交換方式となっています。
//
// 【楕円曲線暗号の基礎】
// - 楕円曲線上の点の加算を基本演算として使用
// - 離散対数問題の楕円曲線版（ECDLP）の困難性に基づく
// - y² = x³ + ax + b (mod p) の形式の曲線を使用
//
// 【よく使われる曲線】
// - P-256 (secp256r1): 128ビットセキュリティ、最も広く使用
// - P-384 (secp384r1): 192ビットセキュリティ
// - P-521 (secp521r1): 256ビットセキュリティ
// - X25519: 128ビットセキュリティ、高速で安全、TLS 1.3推奨
//
// 【仕組み】
// 1. 使用する楕円曲線とベースポイント G を決定
// 2. アリス: 秘密値 a を選び、A = aG（スカラー倍）を計算して送信
// 3. ボブ:   秘密値 b を選び、B = bG を計算して送信
// 4. アリス: 共有点 K = aB = abG を計算
// 5. ボブ:   共有点 K = bA = abG を計算
// → 両者が同じ点 K を得る（そのx座標が共有シークレット）
//
// 【セキュリティ特性】
// - 前方秘匿性（Forward Secrecy）あり
// - DHEより小さな鍵サイズで同等のセキュリティ
//   → P-256 (256ビット) ≈ RSA/DH 3072ビット
//   → P-384 (384ビット) ≈ RSA/DH 7680ビット
// - 計算効率が高い
//
// 【鍵サイズ比較】
// +------------------+-------------------+------------------+
// | セキュリティ強度  | RSA/DH 鍵サイズ    | ECC 鍵サイズ      |
// +------------------+-------------------+------------------+
// | 80ビット         | 1024ビット        | 160ビット        |
// | 112ビット        | 2048ビット        | 224ビット        |
// | 128ビット        | 3072ビット        | 256ビット        |
// | 192ビット        | 7680ビット        | 384ビット        |
// | 256ビット        | 15360ビット       | 521ビット        |
// +------------------+-------------------+------------------+
//
// =============================================================================

// ECDHEKeyExchange はElliptic Curve Diffie-Hellman Ephemeral鍵交換を実装します
type ECDHEKeyExchange struct {
	Curve      ecdh.Curve
	CurveName  string
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

// CurveInfo は楕円曲線の情報を保持します
type CurveInfo struct {
	Name           string
	Curve          ecdh.Curve
	KeySizeBits    int
	SecurityBits   int
	RSAEquivalent  int
	Description    string
}

// GetSupportedCurves はサポートされている楕円曲線のリストを返します
func GetSupportedCurves() []CurveInfo {
	return []CurveInfo{
		{
			Name:          "P-256",
			Curve:         ecdh.P256(),
			KeySizeBits:   256,
			SecurityBits:  128,
			RSAEquivalent: 3072,
			Description:   "NIST P-256。最も広く使用される曲線。TLS 1.2/1.3対応。",
		},
		{
			Name:          "P-384",
			Curve:         ecdh.P384(),
			KeySizeBits:   384,
			SecurityBits:  192,
			RSAEquivalent: 7680,
			Description:   "NIST P-384。高セキュリティ要件向け。政府系システムで使用。",
		},
		{
			Name:          "P-521",
			Curve:         ecdh.P521(),
			KeySizeBits:   521,
			SecurityBits:  256,
			RSAEquivalent: 15360,
			Description:   "NIST P-521。最高レベルのセキュリティ。計算コスト高。",
		},
		{
			Name:          "X25519",
			Curve:         ecdh.X25519(),
			KeySizeBits:   256,
			SecurityBits:  128,
			RSAEquivalent: 3072,
			Description:   "Curve25519ベース。高速で安全。TLS 1.3推奨。サイドチャネル攻撃耐性。",
		},
	}
}

// NewECDHEKeyExchange は指定された曲線でECDHE鍵交換を初期化します
func NewECDHEKeyExchange(curveName string) (*ECDHEKeyExchange, error) {
	var curve ecdh.Curve

	switch curveName {
	case "P-256":
		curve = ecdh.P256()
	case "P-384":
		curve = ecdh.P384()
	case "P-521":
		curve = ecdh.P521()
	case "X25519":
		curve = ecdh.X25519()
	default:
		return nil, fmt.Errorf("未対応の曲線: %s", curveName)
	}

	// 秘密鍵を生成
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("秘密鍵生成に失敗: %w", err)
	}

	return &ECDHEKeyExchange{
		Curve:      curve,
		CurveName:  curveName,
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}, nil
}

// ComputeSharedSecret は相手の公開鍵から共有シークレットを計算します
func (e *ECDHEKeyExchange) ComputeSharedSecret(otherPublicKey *ecdh.PublicKey) ([]byte, error) {
	return e.PrivateKey.ECDH(otherPublicKey)
}

// DemonstrateECDHEKeyExchange はECDHE鍵交換の完全なフローを示します
func DemonstrateECDHEKeyExchange(curveName string) error {
	fmt.Printf("\n%s\n", "==========================================================")
	fmt.Printf("ECDHE鍵交換 (Elliptic Curve Diffie-Hellman Ephemeral)\n")
	fmt.Printf("使用曲線: %s\n", curveName)
	fmt.Printf("%s\n\n", "==========================================================")

	// 曲線情報を表示
	curves := GetSupportedCurves()
	for _, c := range curves {
		if c.Name == curveName {
			fmt.Println("【曲線パラメータ】")
			fmt.Printf("  曲線名: %s\n", c.Name)
			fmt.Printf("  鍵サイズ: %d ビット\n", c.KeySizeBits)
			fmt.Printf("  セキュリティ強度: %d ビット\n", c.SecurityBits)
			fmt.Printf("  RSA等価鍵長: %d ビット\n", c.RSAEquivalent)
			fmt.Printf("  説明: %s\n", c.Description)
			break
		}
	}

	// 1. アリス（クライアント）の鍵生成
	fmt.Println("\n【ステップ1】アリス（クライアント）が一時鍵ペアを生成...")
	alice, err := NewECDHEKeyExchange(curveName)
	if err != nil {
		return err
	}
	alicePubBytes := alice.PublicKey.Bytes()
	fmt.Printf("  アリスの公開鍵サイズ: %d バイト\n", len(alicePubBytes))
	fmt.Printf("  アリスの公開鍵 (先頭16バイト): %x...\n", alicePubBytes[:min(16, len(alicePubBytes))])

	// 2. ボブ（サーバー）の鍵生成
	fmt.Println("\n【ステップ2】ボブ（サーバー）が一時鍵ペアを生成...")
	bob, err := NewECDHEKeyExchange(curveName)
	if err != nil {
		return err
	}
	bobPubBytes := bob.PublicKey.Bytes()
	fmt.Printf("  ボブの公開鍵サイズ: %d バイト\n", len(bobPubBytes))
	fmt.Printf("  ボブの公開鍵 (先頭16バイト): %x...\n", bobPubBytes[:min(16, len(bobPubBytes))])

	// 3. 公開鍵を交換
	fmt.Println("\n【ステップ3】公開鍵を交換...")
	fmt.Println("  アリス → ボブ: 公開鍵 A を送信")
	fmt.Println("  ボブ → アリス: 公開鍵 B を送信")

	// 4. 共有シークレットを計算
	fmt.Println("\n【ステップ4】それぞれが共有シークレットを計算...")
	aliceSharedSecret, err := alice.ComputeSharedSecret(bob.PublicKey)
	if err != nil {
		return fmt.Errorf("アリスの共有シークレット計算に失敗: %w", err)
	}
	bobSharedSecret, err := bob.ComputeSharedSecret(alice.PublicKey)
	if err != nil {
		return fmt.Errorf("ボブの共有シークレット計算に失敗: %w", err)
	}
	fmt.Printf("  アリスが計算した共有シークレット (先頭16バイト): %x...\n", aliceSharedSecret[:min(16, len(aliceSharedSecret))])
	fmt.Printf("  ボブが計算した共有シークレット (先頭16バイト): %x...\n", bobSharedSecret[:min(16, len(bobSharedSecret))])

	// 5. 検証
	fmt.Println("\n【検証】共有シークレットが一致するか確認...")
	if string(aliceSharedSecret) == string(bobSharedSecret) {
		fmt.Println("  ✓ 成功: アリスとボブで同じ共有シークレットを確立")
		fmt.Printf("  共有シークレットのサイズ: %d バイト (%d ビット)\n",
			len(aliceSharedSecret), len(aliceSharedSecret)*8)
	} else {
		fmt.Println("  ✗ 失敗: 共有シークレットが一致しません")
	}

	return nil
}

// min は2つの整数の小さい方を返します
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
