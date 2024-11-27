package main

import (
	"fmt"
	"time"
)

func Part18() {
	part185()
}

func part181() {
	t := time.Now()
	fmt.Println(t.String())
}


func part1842() {
	fmt.Println("waiting 5 seconds")
	time.Sleep(5 * time.Second)
	fmt.Println("done")
}

func part1843() {
	fmt.Println("waiting 5 seconds")
	after := time.After(5 * time.Second)
	<-after
	fmt.Println("done")

	fmt.Println("waiting 5 seconds")
	for now := range time.Tick(5 * time.Second) {
		fmt.Println("now: ", now)
	}
}

func part185() {
	now := time.Now()
	fmt.Println(now.String())
	fmt.Println(now.Format(time.RFC822))
	fmt.Println(now.Format(time.RFC850))
	fmt.Println(now.Format(time.RFC1123))

	fmt.Println(now.Format("2006/01/02 03:04:05 MST"))
}