package main

import (
	"fmt"
	"designpattern/factory"
	"designpattern/adapter"
	"designpattern/bridge"
)

func main() {
	fmt.Println("Design Pattern")

	fmt.Println("==========Strategy Pattern==========")
	storategyExec()

	fmt.Println("==========Abstract Factory Pattern==========")
	abstractFactoryExec()

	fmt.Println("==========Singleton Pattern==========")
	singletonExec()

	fmt.Println("==========Factory Pattern==========")
	factory.FactoryExec()

	fmt.Println("==========Adaptor Pattern==========")
	adapter.AdapterExec()

	fmt.Println("==========Bridge Pattern==========")
	bridge.BridgeExec()
}
