package iodemo

import (
	"errors"
	"fmt"
	"os"
)

func ReadDemo() {
	fmt.Println("Read demo")

	f, _ := os.Open("test.txt")
	defer f.Close()

	data := make([]byte, 1024)
	c, _ := f.Read(data)

	fmt.Println("length:", c)
	fmt.Println("string:", string(data[:c]))
	fmt.Println("bytes:", data[:c])
	fmt.Println("file name:", f.Name())

	linkErr := os.LinkError{
		Op: "Op",
		Old: "Old",
		New: "New",
		Err: errors.New("Err"),
	}

	fmt.Println(linkErr.Err)
	fmt.Println(linkErr.Unwrap())
}
