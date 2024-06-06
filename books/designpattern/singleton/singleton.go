package singleton

import (
	"fmt"
	"sync"
)

type singleton struct {
	value string
}

var singletonInstans *singleton
var once sync.Once

func GetInstance(s string) *singleton {
	once.Do(func ()  {
		singletonInstans = &singleton{value: s}
	})
	return singletonInstans
}

func Exec() {
	instance1 := GetInstance("Hello World!")
	fmt.Println(instance1.value)

	instance2 := GetInstance("Hello World!!!!!!!!")
	fmt.Println(instance2.value)

	instance3 := singleton{value: "hoge!"} // 同一パッケージ内だと複数生成できる。
	fmt.Println(instance3.value)

	if instance1 == instance2 {
		fmt.Println("Both instances are the same")
	} else {
		fmt.Println("Instances are different")
	}
}