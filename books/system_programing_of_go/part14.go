package main

import (
	"fmt"
	"sync"
	"time"
)

func Part14() {
	part1474()
}

func sub1(c int) {
	fmt.Println("share by arguments:", c*c)
}

func part14211() {
	// 引数渡し
	go sub1(10)

	// クロージャのキャプチャ渡し
	d := 20
	go func() {
		fmt.Println("share by capture", d*d)
	}()
	time.Sleep(time.Second)
}

func part14212() {
	tasks := []string{
		"cmake ..",
		"cmake . --build Release",
		"cpack",
	}
	for _, task := range tasks {
		go func() {
			// goroutineが起動するときにはループが回りきって
			// 全部のtaskが最後のタスクになってしまう
			fmt.Println(task)
		}()
	}
	time.Sleep(time.Second)
}

var id int

func generateId(mutex *sync.Mutex) int {
	// Lock() / Unlock()をペアで呼び出してロックする
	mutex.Lock()
	defer mutex.Unlock()
	id++
	return id
}

func part1471() {
	// sync.Mutex構造体の変数宣言
	// 次の宣言をしてもポインタ型になるだけで正常に動作する
	// mutex := new(sync.Mutex)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Printf("id: %d\n", generateId(&mutex))
			wg.Done()
		}()
	}
	wg.Wait()
}

func part1472() {
	var wg sync.WaitGroup

	// ジョブ数をあらかじめ登録
	wg.Add(2)

	go func() {
		// 非同期で仕事をする (1)
		fmt.Println("仕事1")
		// Doneで完了を通知
		wg.Done()
	}()
	go func() {
		// 非同期で仕事をする (2)
		fmt.Println("仕事2")
		// Doneで完了を通知
		wg.Done()
	}()

	// すべての処理が終わるのを待つ
	wg.Wait()
	fmt.Println("終了")
}

func initialize() {
	fmt.Println("初期化処理")
}

var once sync.Once

func part1473() {
	// 3回呼び出しても一度しか呼ばれない
	once.Do(initialize)
	once.Do(initialize)
	once.Do(initialize)
}

func part1474() {
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)

	for _, name := range []string{"A", "B", "C"} {
		go func(name string) {
			// ロックしてからWaitメソッドを呼ぶ
			mutex.Lock()
			defer mutex.Unlock()
			// Broadcast()が呼ばれるまで待つ
			cond.Wait()
			// 呼ばれた
			fmt.Println(name)
		}(name)
	}

	fmt.Println("よーい")
	time.Sleep(time.Second)
	fmt.Println("どん！")
	// 待っているgoroutineを一斉に起こす
	cond.Broadcast()
	time.Sleep(time.Second)
}