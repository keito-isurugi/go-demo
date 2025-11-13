package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("UUIDv4 と UUIDv7 のソート順序比較")
	fmt.Println(strings.Repeat("=", 80))

	// UUIDv4（ランダム）を生成
	fmt.Println("\n【UUIDv4】ランダムに生成される")
	uuidv4List := make([]string, 10)
	for i := 0; i < 10; i++ {
		uuidv4List[i] = uuid.New().String()
	}

	fmt.Println("\n生成順:")
	for i, id := range uuidv4List {
		fmt.Printf("%2d: %s\n", i+1, id)
	}

	fmt.Println("\nソート後（辞書順）:")
	sorted := make([]string, len(uuidv4List))
	copy(sorted, uuidv4List)
	sort.Strings(sorted)
	for i, id := range sorted {
		fmt.Printf("%2d: %s\n", i+1, id)
	}

	fmt.Println("\n→ 生成順とソート順が全く違う！（バラバラに配置される）")

	// UUIDv7（タイムスタンプベース）を生成
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("【UUIDv7】タイムスタンプベースで生成される")
	fmt.Println(strings.Repeat("=", 80))

	uuidv7List := make([]string, 10)
	for i := 0; i < 10; i++ {
		uuidv7List[i] = uuid.Must(uuid.NewV7()).String()
	}

	fmt.Println("\n生成順:")
	for i, id := range uuidv7List {
		fmt.Printf("%2d: %s\n", i+1, id)
		// タイムスタンプ部分を表示
		if i == 0 {
			fmt.Printf("     ^^^^^^^^^^ タイムスタンプ部分（時系列順）\n")
		}
	}

	fmt.Println("\nソート後（辞書順）:")
	sortedv7 := make([]string, len(uuidv7List))
	copy(sortedv7, uuidv7List)
	sort.Strings(sortedv7)
	for i, id := range sortedv7 {
		fmt.Printf("%2d: %s\n", i+1, id)
	}

	fmt.Println("\n→ 生成順とソート順が一致！（連続して配置される）")

	// 16進数としての比較を可視化
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("16進数としての比較の仕組み")
	fmt.Println(strings.Repeat("=", 80))

	examples := []string{
		"00055ab0-f3ee-47ab-becf-ffb1487d1b4a",
		"0007ba31-4a08-4937-b737-f6b197d338b6",
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
		"80000000-0000-0000-0000-000000000000",
	}

	fmt.Println("\n辞書順でソート:")
	sort.Strings(examples)
	for i, id := range examples {
		// 先頭8文字を16進数として解釈
		fmt.Printf("%2d: %s (先頭: 0x%s)\n", i+1, id, strings.ReplaceAll(id[:8], "-", ""))
	}

	fmt.Println("\n説明:")
	fmt.Println("- 文字列として1文字ずつ比較")
	fmt.Println("- '0' < '8' < 'f' という順序")
	fmt.Println("- 実質的に16進数の大小比較と同じ結果")
}
