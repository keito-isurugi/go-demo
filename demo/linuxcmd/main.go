package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 標準入力を受け取る
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter command: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// 改行を取り除く
	input = input[:len(input)-1]

	// コマンドと引数を分ける
	cmdArgs := []string{}
	cmdName := ""
	if len(input) > 0 {
		// 最初の単語をコマンド名として分ける
		cmdName = input[:len(input)]
		cmdArgs = append(cmdArgs, cmdName)
	}

	// コマンドを実行
	cmd := exec.Command("bash", "-c", input)
	output, err := cmd.CombinedOutput()

	// エラーが発生した場合
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		fmt.Printf("Command output: %s\n", string(output))
		return
	}

	// コマンドの出力を表示
	fmt.Println("Command executed successfully.")
	fmt.Println(string(output))
}
