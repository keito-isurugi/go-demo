package main

import (
	"fmt"
	"net"
	"os"
)

func Part7() {
	args := os.Args[1:] // プログラム名を除く引数を取得

	if len(args) == 0 {
		fmt.Println("Usage: go run . [part721|part722]")
		return
	}

	switch args[0] {
	case "part721":
		part721()
	case "part722":
		part722()
	default:
		fmt.Printf("Unknown argument for part7: %s\n", args[0])
	}

}

func part721() {
	fmt.Println("Server is running at localhost:8888")
	conn, err := net.ListenPacket("udp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		length, remoteAddress, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf(
			"Received from %v: %v\n",
			remoteAddress, string(buffer[:length]),
		)
		_, err = conn.WriteTo([]byte("Hello from Server"), remoteAddress)
		if err != nil {
			panic(err)
		}
	}
}

func part722() {
	conn, err := net.Dial("udp4", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Sending to server")
	_, err = conn.Write([]byte("Hello from Client"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Receiving from server")
	buffer := make([]byte, 1500)
	length, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %s\n", string(buffer[:length]))
}
