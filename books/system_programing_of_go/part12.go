package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kr/pty"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/unix"
)

func Part12() {
	part12622()
}

func part1211() {
	path, _ := os.Executable()
	fmt.Printf("実行ファイル名: %s\n", os.Args[0])
	fmt.Printf("実行ファイルパス: %s\n", path)
}

func part1212() {
	fmt.Printf("プロセスID: %d\n", os.Getpid())
	fmt.Printf("親プロセスID: %d\n", os.Getppid())
}

func part1213() {
	sid, _ := unix.Getsid(os.Getpid())
	fmt.Fprintf(os.Stderr, "グループID: %d セッションID: %d\n", unix.Getpgrp(), sid)
}

func part1214() {
	fmt.Printf("ユーザーID: %d\n", os.Getuid())
	fmt.Printf("グループID: %d\n", os.Getgid())
	groups, _ := os.Getgroups()
	fmt.Printf("サブグループID: %v\n", groups)
}

func part1215() {
	fmt.Printf("ユーザーID: %d\n", os.Getuid())
	fmt.Printf("グループID: %d\n", os.Getgid())
	fmt.Printf("実効ユーザーID: %d\n", os.Geteuid())
	fmt.Printf("実効グループID: %d\n", os.Getegid())
}

func part1216() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}

func part1222() {
	fmt.Println(os.ExpandEnv("${HOME}/gobin"))
}

func part1223() {
	os.Exit(1)
}

func part123() {
	p, _ := process.NewProcess(int32(os.Getppid()))
	name, _ := p.Name()
	cmd, _ := p.Cmdline()
	fmt.Printf("parent pid: %d name: '%s' cmd: '%s'\n", p.Pid, name, cmd)
}

func part12511() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}

func part12512() {
	count := exec.Command("./count")
	stdout, _ := count.StdoutPipe()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("(stdout) %s\n", scanner.Text())
		}
	}()
	err := count.Run()
	if err != nil {
		panic(err)
	}
}

func part125() {
	if len(os.Args) == 1 {
		return
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	state := cmd.ProcessState
	fmt.Printf("%s\n", state.String())
	fmt.Printf(" Pid: %d\n", state.Pid())
	fmt.Printf(" Exited: %v\n", state.Exited())
	fmt.Printf(" Success: %v\n", state.Success())
	fmt.Printf(" System: %v\n", state.SystemTime())
	fmt.Printf(" User: %v\n", state.UserTime())
}

var data = "\033[34m\033[47m\033[4mB\033[31me\n\033[24m\033[30mOS\033[49m\033[m\n"

func part126() {
	var stdOut io.Writer
	if isatty.IsTerminal(os.Stdout.Fd()) {
		stdOut = colorable.NewColorableStdout()
	} else {
		stdOut = colorable.NewColorable(os.Stdout)
	}
	fmt.Println(stdOut, data)
}

func part12621() {
	var out io.Writer
	if isatty.IsTerminal(os.Stdout.Fd()) {
		out = colorable.NewColorableStdout()
	} else {
		out = colorable.NewColorable(os.Stdout)
	}
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprintln(out, "stdin: terminal")
	} else {
		fmt.Println("stdin: pipe")
	}
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Fprintln(out, "stdout: terminal")
	} else {
		fmt.Println("stdout: pipe")
	}
	if isatty.IsTerminal(os.Stderr.Fd()) {
		fmt.Fprintln(out, "stderr: terminal")
	} else {
		fmt.Println("stderr: pipe")
	}
}

func part12622() {
	cmd := exec.Command("./check")
	stdpty, stdtty, _ := pty.Open()
	defer stdtty.Close()
	cmd.Stdin = stdpty
	cmd.Stdout = stdpty
	errpty, errtty, _ := pty.Open()
	defer errtty.Close()
	cmd.Stderr = errtty
	go func() {
		io.Copy(os.Stdout, stdpty)
	}()
	go func() {
		io.Copy(os.Stderr, errpty)
	}()
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
