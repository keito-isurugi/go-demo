package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/edsrzf/mmap-go"
	"github.com/fsnotify/fsnotify"
)

func Part10() {
	part103()
}

func part101() {
	counter := 0
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event: ", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file: ", event.Name)
					counter++
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
					counter++
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file: ", event.Name)
					counter++
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("renamed file: ", event.Name)
					counter++
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("chmod file: ", event.Name)
					counter++
				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
			if counter > 3 {
				done <- true
			}
		}

	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

type FileLock struct {
	l  sync.Mutex
	fd int
}

func NewFileLock(filename string) *FileLock {
	if filename == "" {
		panic("filename needed")
	}
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, 0750)
	if err != nil {
		panic(err)
	}
	return &FileLock{fd: fd}
}

func (m *FileLock) Lock() {
	m.l.Lock()
	if err := syscall.Flock(m.fd, syscall.LOCK_EX); err != nil {
		panic(err)
	}
}

func (m *FileLock) UnLock() {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		panic(err)
	}
	m.l.Unlock()
}

func part1022() {
	file := NewFileLock("file.txt")
	fmt.Println("try locking...")
	file.Lock()
	fmt.Println("locked!")
	time.Sleep(15 * time.Second)
	file.UnLock()
	fmt.Println("unlock")
}

func part103() {
	// テストデータを書き込み
	testData := []byte("0123456789ABCDEF")
	testPath := filepath.Join(os.TempDir(), "testdata")
	err := ioutil.WriteFile(testPath, testData, 0644)
	if err != nil {
		panic(err)
	}

	// メモリにマッピング
	f, err := os.OpenFile(testPath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// mはmmap.MMap型で[]byteのエイリアス
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer m.Unmap()

	// メモリ上のデータを修正して書き込む
	m[9] = 'X'
	m.Flush()

	// 読み込んでみる
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("original: %s\n", testData)
	fmt.Printf("mmap:     %s\n", m)
	fmt.Printf("file:     %s\n", fileData)
}
