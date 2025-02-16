package main

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type singleton struct{}

var singleInstance *singleton

func getInstance() *singleton {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating Single Instance Now")
			singleInstance = &singleton{}
		} else {
			fmt.Println("Single Instance already created")
		}
	} else {
		fmt.Println("Single Instance already created")
	}

	return singleInstance
}

func singletonExec() {
	for i := 0; i < 5; i++ {
		go getInstance()
	}
	fmt.Scanln()
}