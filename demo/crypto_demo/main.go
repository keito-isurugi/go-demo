package main

import (
	"fmt"
	rsaDemo "crypto_demo/rsa_demo"
)

func main() {
	fmt.Println("======crypto======")
	// privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// if err != nil {
	// 	log.Fatalf("鍵の生成に失敗しました: %v", err)
	// }

	// fmt.Println("privateKey: ", privateKey.PublicKey)
	rsaDemo.ExecRSA()
}
