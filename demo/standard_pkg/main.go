package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("test.txt") // 読み込み専用
	// f, err := os.Create("test.txt") // 書き込み権限あり
	defer func () { // 後続にWriteを使用する場合はerrをキャッチしたほうがいいかも
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	f2, _ := os.Create("test2.txt")
	defer f2.Close()

	// ファイルディスクリプタ
	fmt.Println(f.Fd())
	fmt.Println(f2.Fd())

	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data[:count]))
	fmt.Println(data[:count])
}
