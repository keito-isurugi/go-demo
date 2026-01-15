package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// =============================================================================
// DHE鍵交換 (Diffie-Hellman Ephemeral)
// =============================================================================
//
// 【概要】
// Diffie-Hellman鍵交換は、1976年にWhitfield DiffieとMartin Hellmanによって
// 発明された、安全でない通信路上で共通鍵を確立するための画期的な方式です。
// "Ephemeral"（一時的）とは、セッションごとに新しい鍵ペアを生成することを意味します。
//
// 【数学的基礎】
// - 大きな素数 p と生成元 g を公開パラメータとして使用
// - 離散対数問題の困難性に基づく
//   → g^a mod p = A が与えられても、a を求めることは計算量的に困難
//
// 【仕組み】
// 1. 公開パラメータ（素数p、生成元g）を決定
// 2. アリス: 秘密値 a を選び、A = g^a mod p を計算して送信
// 3. ボブ:   秘密値 b を選び、B = g^b mod p を計算して送信
// 4. アリス: 共有鍵 K = B^a mod p = g^(ab) mod p を計算
// 5. ボブ:   共有鍵 K = A^b mod p = g^(ab) mod p を計算
// → 両者が同じ K を得る
//
// 【セキュリティ特性】
// - 前方秘匿性（Forward Secrecy）あり
//   → セッションごとに一時鍵を使用するため、長期鍵が漏洩しても過去の通信は安全
// - パラメータサイズによる強度
//   → 2048ビット: 現在の推奨最小サイズ
//   → 3072ビット: 128ビットセキュリティ相当
//   → 4096ビット: 長期的なセキュリティ向け
//
// 【注意点】
// - 中間者攻撃（MITM）に対して脆弱（認証が必要）
// - ECDHEと比較して計算コストが高い
// - パラメータの品質が重要（弱いパラメータは攻撃される可能性）
//
// =============================================================================

// DHEKeyExchange はDiffie-Hellman Ephemeral鍵交換を実装します
type DHEKeyExchange struct {
	P          *big.Int // 大きな素数
	G          *big.Int // 生成元
	PrivateKey *big.Int // 秘密鍵 (a または b)
	PublicKey  *big.Int // 公開鍵 (g^a mod p または g^b mod p)
	BitSize    int
}

// 事前定義されたDHパラメータ (RFC 3526 MODP Groups)
// 実際の実装では、より大きなグループを使用すべきです

// GetDHParameters は指定されたビットサイズのDHパラメータを返します
func GetDHParameters(bits int) (*big.Int, *big.Int) {
	var pHex string

	switch bits {
	case 1024:
		// 1024ビット MODP Group (RFC 2409) - 現在は非推奨
		pHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
			"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
			"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
			"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
			"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381" +
			"FFFFFFFFFFFFFFFF"
	case 2048:
		// 2048ビット MODP Group (RFC 3526 Group 14)
		pHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
			"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
			"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
			"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
			"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
			"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
			"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
			"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
			"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
			"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
			"15728E5A8AACAA68FFFFFFFFFFFFFFFF"
	default:
		// デフォルトは2048ビット
		return GetDHParameters(2048)
	}

	p := new(big.Int)
	p.SetString(pHex, 16)

	// 生成元は通常 2
	g := big.NewInt(2)

	return p, g
}

// NewDHEKeyExchange は新しいDHE鍵交換インスタンスを作成します
func NewDHEKeyExchange(bits int) (*DHEKeyExchange, error) {
	p, g := GetDHParameters(bits)

	// 秘密鍵を生成 (1 < privateKey < p-1)
	// 秘密鍵のビット長はpの半分程度で十分
	privateKeyBits := bits / 2
	privateKey, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
	if err != nil {
		return nil, fmt.Errorf("秘密鍵生成に失敗: %w", err)
	}
	privateKey.Add(privateKey, big.NewInt(1)) // 1以上にする

	// 公開鍵を計算: publicKey = g^privateKey mod p
	publicKey := new(big.Int).Exp(g, privateKey, p)

	return &DHEKeyExchange{
		P:          p,
		G:          g,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		BitSize:    privateKeyBits,
	}, nil
}

// ComputeSharedSecret は相手の公開鍵から共有シークレットを計算します
func (d *DHEKeyExchange) ComputeSharedSecret(otherPublicKey *big.Int) *big.Int {
	// sharedSecret = otherPublicKey^privateKey mod p
	return new(big.Int).Exp(otherPublicKey, d.PrivateKey, d.P)
}

// DemonstrateDHEKeyExchange はDHE鍵交換の完全なフローを示します
func DemonstrateDHEKeyExchange(bits int) error {
	fmt.Printf("\n%s\n", "==========================================================")
	fmt.Printf("DHE鍵交換 (Diffie-Hellman Ephemeral, %dビット)\n", bits)
	fmt.Printf("%s\n\n", "==========================================================")

	// 公開パラメータを表示
	p, g := GetDHParameters(bits)
	fmt.Println("【公開パラメータ】")
	fmt.Printf("  素数 p のビット長: %d ビット\n", p.BitLen())
	fmt.Printf("  生成元 g: %s\n", g.String())

	// 1. アリス（クライアント）の鍵生成
	fmt.Println("\n【ステップ1】アリス（クライアント）が一時鍵ペアを生成...")
	alice, err := NewDHEKeyExchange(bits)
	if err != nil {
		return err
	}
	fmt.Printf("  アリスの秘密鍵 a (最初の64ビット): %x...\n", alice.PrivateKey.Bytes()[:8])
	fmt.Printf("  アリスの公開鍵 A = g^a mod p (最初の64ビット): %x...\n", alice.PublicKey.Bytes()[:8])

	// 2. ボブ（サーバー）の鍵生成
	fmt.Println("\n【ステップ2】ボブ（サーバー）が一時鍵ペアを生成...")
	bob, err := NewDHEKeyExchange(bits)
	if err != nil {
		return err
	}
	fmt.Printf("  ボブの秘密鍵 b (最初の64ビット): %x...\n", bob.PrivateKey.Bytes()[:8])
	fmt.Printf("  ボブの公開鍵 B = g^b mod p (最初の64ビット): %x...\n", bob.PublicKey.Bytes()[:8])

	// 3. 公開鍵を交換（実際のTLSでは署名付きで送信）
	fmt.Println("\n【ステップ3】公開鍵を交換...")
	fmt.Println("  アリス → ボブ: 公開鍵 A を送信")
	fmt.Println("  ボブ → アリス: 公開鍵 B を送信")

	// 4. 共有シークレットを計算
	fmt.Println("\n【ステップ4】それぞれが共有シークレットを計算...")
	aliceSharedSecret := alice.ComputeSharedSecret(bob.PublicKey)
	bobSharedSecret := bob.ComputeSharedSecret(alice.PublicKey)
	fmt.Printf("  アリスが計算: K = B^a mod p (最初の64ビット): %x...\n", aliceSharedSecret.Bytes()[:8])
	fmt.Printf("  ボブが計算:   K = A^b mod p (最初の64ビット): %x...\n", bobSharedSecret.Bytes()[:8])

	// 5. 検証
	fmt.Println("\n【検証】共有シークレットが一致するか確認...")
	if aliceSharedSecret.Cmp(bobSharedSecret) == 0 {
		fmt.Println("  ✓ 成功: アリスとボブで同じ共有シークレットを確立")
		fmt.Printf("  共有シークレットのサイズ: %d バイト (%d ビット)\n",
			len(aliceSharedSecret.Bytes()), aliceSharedSecret.BitLen())
	} else {
		fmt.Println("  ✗ 失敗: 共有シークレットが一致しません")
	}

	// 前方秘匿性の説明
	fmt.Println("\n【前方秘匿性について】")
	fmt.Println("  DHEでは各セッションで新しい一時鍵ペアを生成するため、")
	fmt.Println("  サーバーの長期秘密鍵が漏洩しても、過去のセッションで")
	fmt.Println("  使用された一時鍵は復元できず、通信内容は保護されます。")

	return nil
}
