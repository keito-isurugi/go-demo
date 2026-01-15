package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"
)

// =============================================================================
// 鍵交換アルゴリズムの比較と分析
// =============================================================================
//
// 【総合比較表】
// +-------------+---------------+------------------+------------------+----------+
// | アルゴリズム | 前方秘匿性    | 推奨鍵サイズ      | 計算速度         | TLS 1.3  |
// +-------------+---------------+------------------+------------------+----------+
// | RSA         | なし          | 2048-4096ビット  | 中（暗号化高速） | 非対応   |
// | DHE         | あり          | 2048-4096ビット  | 低（べき乗計算） | 非推奨   |
// | ECDHE       | あり          | 256-521ビット    | 高               | 推奨     |
// +-------------+---------------+------------------+------------------+----------+
//
// 【セキュリティ強度の等価性】
// +--------------+--------+----------+---------+
// | セキュリティ | RSA/DH | ECDHE    | 対称鍵  |
// +--------------+--------+----------+---------+
// | 80ビット    | 1024   | 160      | 80      |
// | 112ビット   | 2048   | 224      | 112     |
// | 128ビット   | 3072   | 256      | 128     |
// | 192ビット   | 7680   | 384      | 192     |
// | 256ビット   | 15360  | 521      | 256     |
// +--------------+--------+----------+---------+
//
// 【推奨事項】
// 1. 新規実装: ECDHE (X25519またはP-256) を使用
// 2. TLS設定: TLS 1.3 + ECDHE を優先
// 3. レガシー: DHE 2048ビット以上（RSA鍵交換は避ける）
//
// =============================================================================

// ComparisonResult はベンチマーク結果を保持します
type ComparisonResult struct {
	Algorithm    string
	KeySize      string
	KeyGenTime   time.Duration
	ExchangeTime time.Duration
	TotalTime    time.Duration
	SharedSecret int // バイト数
}

// RunBenchmark は各アルゴリズムのベンチマークを実行します
func RunBenchmark(iterations int) []ComparisonResult {
	results := []ComparisonResult{}

	// RSA 2048ビット
	results = append(results, benchmarkRSA(2048, iterations))

	// RSA 4096ビット
	results = append(results, benchmarkRSA(4096, iterations))

	// DHE 2048ビット
	results = append(results, benchmarkDHE(2048, iterations))

	// ECDHE P-256
	results = append(results, benchmarkECDHE("P-256", iterations))

	// ECDHE P-384
	results = append(results, benchmarkECDHE("P-384", iterations))

	// ECDHE X25519
	results = append(results, benchmarkECDHE("X25519", iterations))

	return results
}

func benchmarkRSA(bits int, iterations int) ComparisonResult {
	var totalKeyGen, totalExchange time.Duration

	for i := 0; i < iterations; i++ {
		// 鍵生成
		start := time.Now()
		privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
		totalKeyGen += time.Since(start)

		// 鍵交換（暗号化・復号）
		start = time.Now()
		secret := make([]byte, 48)
		rand.Read(secret)
		encrypted, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, secret, nil)
		rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, nil)
		totalExchange += time.Since(start)
	}

	return ComparisonResult{
		Algorithm:    "RSA",
		KeySize:      fmt.Sprintf("%dビット", bits),
		KeyGenTime:   totalKeyGen / time.Duration(iterations),
		ExchangeTime: totalExchange / time.Duration(iterations),
		TotalTime:    (totalKeyGen + totalExchange) / time.Duration(iterations),
		SharedSecret: 48,
	}
}

func benchmarkDHE(bits int, iterations int) ComparisonResult {
	var totalKeyGen, totalExchange time.Duration
	p, g := GetDHParameters(bits)

	for i := 0; i < iterations; i++ {
		// アリスの鍵生成
		start := time.Now()
		alicePrivate, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
		alicePublic := new(big.Int).Exp(g, alicePrivate, p)
		// ボブの鍵生成
		bobPrivate, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
		bobPublic := new(big.Int).Exp(g, bobPrivate, p)
		totalKeyGen += time.Since(start)

		// 鍵交換
		start = time.Now()
		new(big.Int).Exp(bobPublic, alicePrivate, p)
		new(big.Int).Exp(alicePublic, bobPrivate, p)
		totalExchange += time.Since(start)
	}

	return ComparisonResult{
		Algorithm:    "DHE",
		KeySize:      fmt.Sprintf("%dビット", bits),
		KeyGenTime:   totalKeyGen / time.Duration(iterations),
		ExchangeTime: totalExchange / time.Duration(iterations),
		TotalTime:    (totalKeyGen + totalExchange) / time.Duration(iterations),
		SharedSecret: bits / 8,
	}
}

func benchmarkECDHE(curveName string, iterations int) ComparisonResult {
	var totalKeyGen, totalExchange time.Duration
	var curve ecdh.Curve

	switch curveName {
	case "P-256":
		curve = ecdh.P256()
	case "P-384":
		curve = ecdh.P384()
	case "X25519":
		curve = ecdh.X25519()
	}

	for i := 0; i < iterations; i++ {
		// 鍵生成
		start := time.Now()
		alicePrivate, _ := curve.GenerateKey(rand.Reader)
		bobPrivate, _ := curve.GenerateKey(rand.Reader)
		totalKeyGen += time.Since(start)

		// 鍵交換
		start = time.Now()
		alicePrivate.ECDH(bobPrivate.PublicKey())
		bobPrivate.ECDH(alicePrivate.PublicKey())
		totalExchange += time.Since(start)
	}

	return ComparisonResult{
		Algorithm:    "ECDHE",
		KeySize:      curveName,
		KeyGenTime:   totalKeyGen / time.Duration(iterations),
		ExchangeTime: totalExchange / time.Duration(iterations),
		TotalTime:    (totalKeyGen + totalExchange) / time.Duration(iterations),
		SharedSecret: 32,
	}
}

// PrintComparison は比較結果を表示します
func PrintComparison(results []ComparisonResult) {
	fmt.Printf("\n%s\n", "══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                      鍵交換アルゴリズム パフォーマンス比較")
	fmt.Printf("%s\n\n", "══════════════════════════════════════════════════════════════════════════════")

	// ヘッダー
	fmt.Printf("%-12s %-12s %-18s %-18s %-18s\n",
		"アルゴリズム", "鍵サイズ", "鍵生成時間", "交換時間", "合計時間")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────")

	// 結果
	for _, r := range results {
		fmt.Printf("%-12s %-12s %-18s %-18s %-18s\n",
			r.Algorithm, r.KeySize,
			r.KeyGenTime.Round(time.Microsecond),
			r.ExchangeTime.Round(time.Microsecond),
			r.TotalTime.Round(time.Microsecond))
	}
}

// PrintSecurityComparison はセキュリティ強度の比較を表示します
func PrintSecurityComparison() {
	fmt.Printf("\n%s\n", "══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                      セキュリティ強度の比較")
	fmt.Printf("%s\n\n", "══════════════════════════════════════════════════════════════════════════════")

	fmt.Println("【等価セキュリティ強度】")
	fmt.Println("┌──────────────────┬─────────────────┬─────────────────┬────────────────┐")
	fmt.Println("│ セキュリティ強度 │ RSA/DH 鍵サイズ │ ECDHE 鍵サイズ  │ 対称鍵 (AES)   │")
	fmt.Println("├──────────────────┼─────────────────┼─────────────────┼────────────────┤")
	fmt.Println("│ 80ビット (非推奨)│ 1024ビット      │ 160ビット       │ 80ビット       │")
	fmt.Println("│ 112ビット        │ 2048ビット      │ 224ビット       │ 112ビット      │")
	fmt.Println("│ 128ビット (推奨) │ 3072ビット      │ 256ビット       │ 128ビット      │")
	fmt.Println("│ 192ビット        │ 7680ビット      │ 384ビット       │ 192ビット      │")
	fmt.Println("│ 256ビット        │ 15360ビット     │ 521ビット       │ 256ビット      │")
	fmt.Println("└──────────────────┴─────────────────┴─────────────────┴────────────────┘")

	fmt.Println("\n【特性比較】")
	fmt.Println("┌─────────────┬─────────────┬─────────────────────┬────────────────┬───────────┐")
	fmt.Println("│ アルゴリズム │ 前方秘匿性  │ 数学的基盤          │ TLS 1.3対応    │ 計算効率  │")
	fmt.Println("├─────────────┼─────────────┼─────────────────────┼────────────────┼───────────┤")
	fmt.Println("│ RSA         │ ✗ なし      │ 素因数分解問題      │ ✗ 非対応      │ 中        │")
	fmt.Println("│ DHE         │ ✓ あり      │ 離散対数問題        │ △ 非推奨      │ 低        │")
	fmt.Println("│ ECDHE       │ ✓ あり      │ 楕円曲線離散対数問題│ ✓ 推奨        │ 高        │")
	fmt.Println("└─────────────┴─────────────┴─────────────────────┴────────────────┴───────────┘")

	fmt.Println("\n【前方秘匿性（Forward Secrecy）について】")
	fmt.Println("  前方秘匿性とは、長期秘密鍵が将来漏洩しても、過去の通信を復号できない")
	fmt.Println("  という性質です。")
	fmt.Println("")
	fmt.Println("  RSA鍵交換: 前方秘匿性なし")
	fmt.Println("    → サーバーの秘密鍵が漏洩すると、記録されていた過去の通信が")
	fmt.Println("      全て復号される可能性があります。")
	fmt.Println("")
	fmt.Println("  DHE/ECDHE: 前方秘匿性あり")
	fmt.Println("    → 各セッションで一時的な鍵ペアを生成するため、長期鍵が漏洩しても")
	fmt.Println("      過去のセッション鍵は復元できません。")

	fmt.Println("\n【量子コンピュータへの耐性】")
	fmt.Println("  現在の公開鍵暗号（RSA、DH、ECDH）は、十分に強力な量子コンピュータが")
	fmt.Println("  実現した場合、Shorのアルゴリズムにより破られる可能性があります。")
	fmt.Println("")
	fmt.Println("  ┌─────────────┬────────────────────────────────────────────┐")
	fmt.Println("  │ アルゴリズム │ 量子コンピュータへの耐性                   │")
	fmt.Println("  ├─────────────┼────────────────────────────────────────────┤")
	fmt.Println("  │ RSA         │ ✗ 脆弱（素因数分解が容易になる）          │")
	fmt.Println("  │ DHE         │ ✗ 脆弱（離散対数が容易になる）            │")
	fmt.Println("  │ ECDHE       │ ✗ 脆弱（楕円曲線離散対数が容易になる）    │")
	fmt.Println("  └─────────────┴────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("  対策: 耐量子暗号（PQC）への移行が進行中")
	fmt.Println("    - NIST PQC標準化: Kyber (ML-KEM)、Dilithium など")
	fmt.Println("    - ハイブリッド方式: ECDHE + Kyber の組み合わせ")

	fmt.Println("\n【推奨設定】")
	fmt.Println("  1. TLS 1.3 を使用（RSA鍵交換が廃止されている）")
	fmt.Println("  2. 鍵交換: ECDHE (X25519 または P-256)")
	fmt.Println("  3. 認証: RSA-2048以上 または ECDSA P-256")
	fmt.Println("  4. 暗号スイート例:")
	fmt.Println("     - TLS_AES_256_GCM_SHA384")
	fmt.Println("     - TLS_CHACHA20_POLY1305_SHA256")
	fmt.Println("     - TLS_AES_128_GCM_SHA256")
}
