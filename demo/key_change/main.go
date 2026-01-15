package main

import (
	"fmt"
	"os"
)

// =============================================================================
//
//	鍵交換アルゴリズム デモプログラム
//	Key Exchange Algorithm Demonstration
//
// =============================================================================
//
// このプログラムは、TLS/SSLで使用される主要な鍵交換アルゴリズムを
// 実装・比較するためのデモンストレーションです。
//
// 【対象アルゴリズム】
// 1. RSA鍵交換（RSA Key Transport）
//   - クライアントがプリマスターシークレットを生成し、サーバーの公開鍵で暗号化
//   - TLS 1.2以前で使用、TLS 1.3では廃止
//   - 前方秘匿性なし
//
// 2. DHE鍵交換（Diffie-Hellman Ephemeral）
//   - 離散対数問題に基づく
//   - 一時鍵を使用するため前方秘匿性あり
//   - 計算コストが高い
//
// 3. ECDHE鍵交換（Elliptic Curve Diffie-Hellman Ephemeral）
//   - 楕円曲線暗号に基づく
//   - 小さな鍵サイズで高いセキュリティ
//   - TLS 1.3で推奨
//
// 【実行方法】
//
//	go run .                    # すべてのデモと比較を実行
//	go run . rsa                # RSA鍵交換のみ
//	go run . dhe                # DHE鍵交換のみ
//	go run . ecdhe              # ECDHE鍵交換のみ
//	go run . benchmark          # ベンチマーク比較
//	go run . comparison         # セキュリティ比較表
//
// =============================================================================
func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     鍵交換アルゴリズム デモプログラム                        ║")
	fmt.Println("║                     Key Exchange Algorithm Demonstration                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")

	args := os.Args[1:]

	if len(args) == 0 {
		// すべてのデモを実行
		runAllDemos()
		return
	}

	switch args[0] {
	case "rsa":
		runRSADemo()
	case "dhe":
		runDHEDemo()
	case "ecdhe":
		runECDHEDemo()
	case "benchmark":
		runBenchmark()
	case "comparison":
		PrintSecurityComparison()
	case "help":
		printHelp()
	default:
		fmt.Printf("不明なコマンド: %s\n", args[0])
		printHelp()
	}
}

func runAllDemos() {
	fmt.Println("\n全てのデモを実行します...")

	// RSA鍵交換
	runRSADemo()

	// DHE鍵交換
	runDHEDemo()

	// ECDHE鍵交換
	runECDHEDemo()

	// ベンチマーク
	runBenchmark()

	// セキュリティ比較
	PrintSecurityComparison()
}

func runRSADemo() {
	fmt.Println("\n┌────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    RSA鍵交換のデモンストレーション                 │")
	fmt.Println("└────────────────────────────────────────────────────────────────────┘")

	// 2048ビットRSA
	if err := DemonstrateRSAKeyExchange(2048); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 4096ビットRSA
	if err := DemonstrateRSAKeyExchange(4096); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
}

func runDHEDemo() {
	fmt.Println("\n┌────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    DHE鍵交換のデモンストレーション                 │")
	fmt.Println("└────────────────────────────────────────────────────────────────────┘")

	// 2048ビットDHE
	if err := DemonstrateDHEKeyExchange(2048); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
}

func runECDHEDemo() {
	fmt.Println("\n┌────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                   ECDHE鍵交換のデモンストレーション                │")
	fmt.Println("└────────────────────────────────────────────────────────────────────┘")

	// 各曲線でデモ
	curves := []string{"P-256", "P-384", "X25519"}
	for _, curve := range curves {
		if err := DemonstrateECDHEKeyExchange(curve); err != nil {
			fmt.Printf("エラー (%s): %v\n", curve, err)
		}
	}
}

func runBenchmark() {
	fmt.Println("\n┌────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                   パフォーマンスベンチマーク                       │")
	fmt.Println("└────────────────────────────────────────────────────────────────────┘")

	fmt.Println("\nベンチマークを実行中... (各アルゴリズム10回の平均)")
	results := RunBenchmark(10)
	PrintComparison(results)
}

func printHelp() {
	fmt.Println(`
使用方法:
  go run .                    全てのデモを実行
  go run . rsa                RSA鍵交換のデモのみ
  go run . dhe                DHE鍵交換のデモのみ
  go run . ecdhe              ECDHE鍵交換のデモのみ
  go run . benchmark          パフォーマンスベンチマーク
  go run . comparison         セキュリティ比較表
  go run . help               このヘルプを表示

各アルゴリズムについて:

  RSA鍵交換:
    - クライアントがプリマスターシークレットを生成
    - サーバーの公開鍵で暗号化して送信
    - 前方秘匿性なし（サーバー鍵漏洩で過去通信も危険）
    - TLS 1.3では廃止

  DHE鍵交換:
    - Diffie-Hellman鍵交換の一時鍵版
    - 離散対数問題に基づく
    - 前方秘匿性あり
    - 計算コストが高い

  ECDHE鍵交換:
    - 楕円曲線を使用したDiffie-Hellman
    - 小さな鍵で高いセキュリティ
    - 前方秘匿性あり
    - TLS 1.3で推奨
`)
}
