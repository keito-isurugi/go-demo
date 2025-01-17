package main

import (
	"flag"
	"fmt"
)

// 実行コマンド: go run . -option1=<オプション1> -option2=<オプション2> -option3=<オプション3>
func main() {
    // オプション
    op1 := flag.String("option1", "オプション1のデフォルト値", "オプション1")
    op2 := flag.String("option2", "オプション2のデフォルト値", "オプション2")
    op3 := flag.String("option3", "オプション3のデフォルト値", "オプション3")

    flag.Parse()

	fmt.Println("option1: " + *op1)
	fmt.Println("option2: " + *op2)
	fmt.Println("option3: " + *op3)
}

// go run . -option1=Hello -option3=World!
// 出力
// option1: Hello
// option2: オプション2のデフォルト値
// option3: World!