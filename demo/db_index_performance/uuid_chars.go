package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("UUIDで使われる文字を調査")
	fmt.Println(strings.Repeat("=", 80))

	// 文字の出現頻度をカウント
	charCount := make(map[rune]int)

	// 10000個のUUIDを生成して文字を集計
	count := 10000
	fmt.Printf("\n%d個のUUIDv4を生成して文字を集計中...\n", count)

	for i := 0; i < count; i++ {
		id := uuid.New().String()
		for _, char := range id {
			if char != '-' { // ハイフンは除外
				charCount[char]++
			}
		}
	}

	fmt.Println("\n使われた文字とその出現回数:")
	fmt.Println(strings.Repeat("-", 80))

	// 0-9
	fmt.Println("\n数字:")
	for ch := '0'; ch <= '9'; ch++ {
		if count, ok := charCount[ch]; ok {
			fmt.Printf("  '%c': %6d回\n", ch, count)
		}
	}

	// a-f
	fmt.Println("\nアルファベット (小文字):")
	for ch := 'a'; ch <= 'f'; ch++ {
		if count, ok := charCount[ch]; ok {
			fmt.Printf("  '%c': %6d回\n", ch, count)
		}
	}

	// g以降が使われているかチェック
	fmt.Println("\ng以降のアルファベット:")
	foundBeyondF := false
	for ch := 'g'; ch <= 'z'; ch++ {
		if count, ok := charCount[ch]; ok {
			fmt.Printf("  '%c': %6d回\n", ch, count)
			foundBeyondF = true
		}
	}
	if !foundBeyondF {
		fmt.Println("  (なし)")
	}

	// 大文字が使われているかチェック
	fmt.Println("\n大文字:")
	foundUpperCase := false
	for ch := 'A'; ch <= 'Z'; ch++ {
		if count, ok := charCount[ch]; ok {
			fmt.Printf("  '%c': %6d回\n", ch, count)
			foundUpperCase = true
		}
	}
	if !foundUpperCase {
		fmt.Println("  (なし)")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("結論:")
	fmt.Println("UUIDで使われるのは 0-9 と a-f のみ（16進数の16文字）")
	fmt.Println(strings.Repeat("=", 80))

	// 実際の例を表示
	fmt.Println("\n実際のUUIDv4の例:")
	for i := 0; i < 5; i++ {
		id := uuid.New().String()
		fmt.Printf("  %s\n", id)
	}

	// 16進数の範囲を説明
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("16進数とは？")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\n10進数と16進数の対応:")
	fmt.Println("  10進数:  0  1  2  3  4  5  6  7  8  9  10  11  12  13  14  15")
	fmt.Println("  16進数:  0  1  2  3  4  5  6  7  8  9   a   b   c   d   e   f")
	fmt.Println("\n16進数の'f'は10進数の15を表す")
	fmt.Println("16進数では1桁で0〜15を表現できる")
}
