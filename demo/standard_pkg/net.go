package main

import (
	"fmt"
	"net"
)

func netDemo() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
	}

	str := "hello, net pkg!!"
	data := []byte(str)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}