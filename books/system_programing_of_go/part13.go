package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/lestrrat/go-server-starter/listener"
)

func Part13() {
	part1350()
}

func part1340() {
	// サイズが1より大きいチャネルを作成
	signals := make(chan os.Signal, 1)
	// 最初のチャネル以降は、可変長引数で任意の数のシグナルを設定可能
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	s := <-signals
	switch s {
	case syscall.SIGINT:
		fmt.Println("SIGINT")
	case syscall.SIGTERM:
		fmt.Println("SIGTERM")
	}
}

func part1341() {
	// 最初の10秒はCtrl+Cで止まる
	fmt.Println("Accept Ctrl + C for 10 seconds")
	time.Sleep(time.Second * 10)

	// 可変長引数で任意の数のシグナルを設定可能
	signal.Ignore(syscall.SIGINT, syscall.SIGHUP)

	// 次の10秒はCtrl+Cを無視する
	fmt.Println("Ignore Ctrl + C for 10 seconds")
	time.Sleep(time.Second * 10)
}

func part1344() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [pid[\n", os.Args[0])
		return
	}
	// 第一引数で指定されたプロセスIDを数値に変換
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	// シグナルを送る
	process.Signal(os.Kill)
	// Killの場合は次のショートカットも利用可能
	process.Kill()
}

func part1350() {
	// シグナル初期化
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	// Server::Starterからもらったソケットを確認
	listeners, err := listener.ListenAll()
	if err != nil {
		panic(err)
	}
	// ウェブサーバーをgoroutineで起動
	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "server pid: %d %v\n", os.Getpid(), os.Environ())
		}),
	}
	go server.Serve(listeners[0])

	// SIGTERMを受け取ったら終了させる
	<-signals
	server.Shutdown(context.Background())

}
