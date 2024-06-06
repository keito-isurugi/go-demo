package singleton

import (
	"fmt"
	"sync"
)

type Singleton struct {
	value string
}

var singletonInstans *Singleton
var once sync.Once

func GetInstance(s string) *Singleton {
	once.Do(func ()  {
		singletonInstans = &Singleton{value: s}
	})
	return singletonInstans
}

func Exec() {
	instance1 := GetInstance("Hello World!")
	fmt.Println(instance1.value)

	instance2 := GetInstance("Hello World!!!!!!!!")
	fmt.Println(instance2.value)

	if instance1 == instance2 {
		fmt.Println("Both instances are the same")
	} else {
		fmt.Println("Instances are different")
	}
}