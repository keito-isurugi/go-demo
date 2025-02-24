package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 引数チェック
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . ls <path>")
		os.Exit(1)
	}

	// コマンドと引数を設定
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// コマンドを実行
	cmd := exec.Command(cmdName, cmdArgs...)

	// 標準出力と標準エラーを取得
	output, err := cmd.CombinedOutput()

	// エラーが発生した場合
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		fmt.Printf("Command output: %s\n", string(output))
		os.Exit(1)
	}

	// コマンドの出力を表示
	fmt.Println("Command executed successfully.")
	fmt.Println(string(output))
}
