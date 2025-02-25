package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// コマンドの入力を受け取る
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// 改行を取り除く
		input = input[:len(input)-1]

		// 終了コマンドを入力されたらループを終了
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// コマンドを実行
		cmd := exec.Command("bash", "-c", input)
		output, err := cmd.CombinedOutput()

		// エラーが発生した場合
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			fmt.Printf("Command output: %s\n", string(output))
			continue
		}

		// コマンドの出力を表示
		fmt.Println("Command executed successfully.")
		fmt.Println(string(output))
	}
}
