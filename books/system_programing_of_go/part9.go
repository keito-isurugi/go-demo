package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Part9() {
	part957()
}

func part921() {
	f, _ := os.Create("file.txt")
	a := time.Now()

	f.Write([]byte("緑の怪獣"))
	b := time.Now()

	f.Sync()
	c := time.Now()
	f.Close()

	d := time.Now()
	fmt.Printf("Write: %v\n", b.Sub(a))
	fmt.Printf("Sync: %v\n", c.Sub(b))
	fmt.Printf("Close: %v\n", d.Sub(c))
}

func part925() {
	if len(os.Args) == 1 {
		fmt.Printf("%s [exec file name]", os.Args[0])
		os.Exit(1)
	}

	info, err := os.Stat(os.Args[1])
	if err == os.ErrNotExist {
		fmt.Printf("file not found: %s\n", os.Args[1])
	} else if err != nil {
		panic(err)
	}
	fmt.Println("FileInfo")
	fmt.Printf("  ファイル名: %v\n", info.Name())
	fmt.Printf("  サイズ: %v\n", info.Size())
	fmt.Printf("  変更日時: %v\n", info.ModTime())
	fmt.Println("Mode()")
	fmt.Printf("  ディレクトリ？ %v\n", info.Mode().IsDir())
	fmt.Printf("  読み書き可能な通常ファイル？ %v\n", info.Mode().IsRegular())
	fmt.Printf("  Unixのファイルアクセス権限ビット %o\n", info.Mode().Perm())
	fmt.Printf("  モードのテキスト表現 %v\n", info.Mode().String())
}

func part9211() {
	dir, err := os.Open("./")
	if err != nil {
		panic(err)
	}
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			fmt.Printf("[Dir]  %s\n", fileInfo.Name())
		} else {
			fmt.Printf("[File] %s\n", fileInfo.Name())
		}
	}
}

func part951() {
	fmt.Printf("Temp File Path: %s\n", filepath.Join(os.TempDir(), "./file.txt"))
}

func part952() {
	dir, name := filepath.Split(os.Getenv("GOPATH"))
	fmt.Printf("Dir: %s, Name: %s\n", dir, name)
}

func part953() {
	if len(os.Args) == 1 {
		fmt.Printf("%s [exec file name]", os.Args[0])
		os.Exit(1)
	}

	for _, path := range filepath.SplitList(os.Getenv("GOPATH")) {
		execpath := filepath.Join(path, os.Args[1])
		_, err := os.Stat(execpath)
		if !os.IsNotExist(err) {
			fmt.Println(execpath)
			return
		}
	}
	os.Exit(1)
}

func part954() {
	// パスをそのままクリーンにする
	fmt.Println(filepath.Clean("./path/filepath/../path.go"))

	// パスを絶対パスに整形
	abspath, _ := filepath.Abs("path/filepath/path_unix.go")
	fmt.Println(abspath)

	// パスを相対パスに整形
	relpath, _ := filepath.Rel("/usr/local/go/src", "/usr/local/go/src/path/filepath/path.go")
	fmt.Println(relpath)
}

var imageSuffix = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
	".tiff": true,
	".eps":  true,
}

func part957() {
	if len(os.Args) == 1 {
		fmt.Printf(`Find images
Usage:
	%s [path to find]
`, os.Args[0])
		return
	}
	root := os.Args[1]
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if info.Name() == "_build" {
					return filepath.SkipDir
				}
				return nil
			}
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if imageSuffix[ext] {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return nil
				}
				fmt.Printf("%s\n", rel)
			}
			return nil
		})
	if err != nil {
		fmt.Println(1, err)
	}
}
