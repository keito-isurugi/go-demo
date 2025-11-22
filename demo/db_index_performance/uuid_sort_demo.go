package main

import (
	"fmt"
	"sort"
)

func main() {
	// 色々な先頭文字のUUID風の文字列
	examples := []string{
		"00000000-0000-0000-0000-000000000000", // 最小
		"0fffffff-ffff-ffff-ffff-ffffffffffff",
		"10000000-0000-0000-0000-000000000000",
		"7fffffff-ffff-ffff-ffff-ffffffffffff",
		"80000000-0000-0000-0000-000000000000",
		"9fffffff-ffff-ffff-ffff-ffffffffffff",
		"a0000000-0000-0000-0000-000000000000",
		"efffffff-ffff-ffff-ffff-ffffffffffff",
		"f0000000-0000-0000-0000-000000000000",
		"ffffffff-ffff-ffff-ffff-ffffffffffff", // 最大
	}

	fmt.Println("ソート前:")
	for i, id := range examples {
		fmt.Printf("%2d: %s\n", i+1, id)
	}

	sort.Strings(examples)

	fmt.Println("\nソート後（辞書順）:")
	for i, id := range examples {
		fmt.Printf("%2d: %s (先頭: '%c')\n", i+1, id, id[0])
	}

	fmt.Println("\n文字の順序: '0' < '9' < 'a' < 'f'")
	fmt.Println("'f'の次は'g'ではなく、存在しない（16進数の限界）")
}
