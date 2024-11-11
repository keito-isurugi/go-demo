package main

import (
	"context"
	"fmt"
	"math"
	"time"
)

func Part4() {
	fmt.Println("start sub()")

	go sub()

	go func() {
		fmt.Println("sub() is runnnig")
		time.Sleep(time.Second)
		fmt.Println("sub() is finished")
	}()

	time.Sleep(2 * time.Second)

	done := make(chan bool)
	go func() {
		fmt.Println("sub() is finished")
		done <- true
	}()
	<-done
	fmt.Println("all tasks are finished")

	pn := primeNumber()
	for n := range pn {
		fmt.Println(n)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("sub() is finished")
		cancel()
	}()
	<-ctx.Done()
	fmt.Println("all tasks are finished")
}

func sub() {
	fmt.Println("sub() is runnnig")
	time.Sleep(time.Second)
	fmt.Println("sub() is finished")
}

func primeNumber() chan int {
	result := make(chan int)
	go func() {
		result <- 2
		for i := 3; i < 100000; i += 2 {
			l := int(math.Sqrt(float64(i)))
			found := false
			for j := 3; j < l+1; j += 2 {
				if i%j == 0 {
					found = true
					break
				}
			}
			if !found {
				result <- i
			}
		}
		close(result)
	}()
	return result
}
