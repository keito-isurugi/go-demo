package main

import "fmt"

func main() {
	fmt.Println("goroutineとchannelの基礎")
	fmt.Println("================================")

	// 各例を順番に実行
	example1_BasicGoroutine()
	example2_BasicChannel()
	example3_ChannelSynchronization()
	example4_BufferedChannel()
	example5_SelectStatement()
	example6_Timeout()
	example7_DoneChannel()
	example8_CloseChannel()

	fmt.Println("\n================================")
	fmt.Println("すべての例が完了しました")
}
