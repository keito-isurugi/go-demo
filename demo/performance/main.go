package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("-----%v回目-----\n", i+1)
		PerformanceHnadler()
	}
}
