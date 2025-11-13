package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("UUIDの有効性チェック")
	fmt.Println(strings.Repeat("=", 80))

	testCases := []string{
		"00000000-0000-0000-0000-00000000000a", // 末尾が'a'
		"00000000-0000-0000-0000-00000000000f", // 末尾が'f'
		"00000000-0000-0000-0000-00000000000g", // 末尾が'g' ← 無効なはず
		"ffffffff-ffff-ffff-ffff-ffffffffffff", // 全部'f'
		"gggggggg-gggg-gggg-gggg-gggggggggggg", // 全部'g' ← 無効なはず
		"12345678-1234-5678-1234-567812345678", // 普通の数字
		"abcdefab-cdef-abcd-efab-cdefabcdefab", // a-fのみ
		"0123456789abcdef-0123-4567-89ab-cdef", // 0-9とa-f混在
		"FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", // 大文字 ← どうなる？
	}

	fmt.Println("\n各文字列の有効性をチェック:")
	fmt.Println(strings.Repeat("-", 80))

	for _, testStr := range testCases {
		parsed, err := uuid.Parse(testStr)
		if err != nil {
			fmt.Printf("❌ 無効: %s\n", testStr)
			fmt.Printf("    エラー: %v\n", err)
		} else {
			fmt.Printf("✅ 有効: %s\n", testStr)
			fmt.Printf("    パース結果: %s\n", parsed.String())
		}
		fmt.Println()
	}

	// 実際のUUIDを生成して末尾が'a'になるケースを探す
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("実際に生成されたUUIDで'a'を含むものを探す")
	fmt.Println(strings.Repeat("=", 80))

	foundExamples := []string{}
	for i := 0; i < 100; i++ {
		id := uuid.New().String()
		// 'a'を含むか？
		if strings.Contains(id, "a") {
			foundExamples = append(foundExamples, id)
			if len(foundExamples) >= 5 {
				break
			}
		}
	}

	fmt.Println("\n'a'を含むUUIDの例:")
	for i, id := range foundExamples {
		fmt.Printf("%d: %s\n", i+1, id)
		// 'a'の位置を表示
		for j, char := range id {
			if char == 'a' {
				fmt.Printf("   %s^ ここに'a' (位置%d)\n", strings.Repeat(" ", j), j)
			}
		}
	}

	// 16進数の説明
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("16進数で有効な文字")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n✅ 有効な文字（16種類）:")
	fmt.Println("  0 1 2 3 4 5 6 7 8 9 a b c d e f")
	fmt.Println("  ↑                       ↑     ↑")
	fmt.Println("  最小値                 'a'は有効  最大値")

	fmt.Println("\n❌ 無効な文字:")
	fmt.Println("  g h i j k l m n o p q r s t u v w x y z")
	fmt.Println("  ↑ これらは16進数では使えない")

	fmt.Println("\n💡 重要な点:")
	fmt.Println("  - 'a' から 'f' までは有効（16進数の10〜15）")
	fmt.Println("  - 'g' から 'z' までは無効（16進数に存在しない）")
	fmt.Println("  - 'a' は「アルファベットの最初」ではなく「16進数の10」")
}
