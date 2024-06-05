package adapter

import  "fmt"

type Print interface{
	printWeak()
	printStrong()
}

type Banner struct {
	String string
}
func NewBanner(s string) *Banner {
	return &Banner{String: s}
}
func (b *Banner) showWithParen() {
	fmt.Println("(", b.String ,")")
}
func (b *Banner) showWithAster() {
	fmt.Println("*", b.String ,"*")
}

type PrintBanner struct {
	Banner *Banner
}
func NewPrintBanner(s string) *PrintBanner{
	return &PrintBanner{Banner: NewBanner(s)}
}
func (pb *PrintBanner) printWeak() {
	pb.Banner.showWithParen()
}
func (pb *PrintBanner) printStrong() {
	pb.Banner.showWithAster()
}

func Exec() {
	p := NewPrintBanner("Hello World!")
	p.printWeak()
	p.printStrong()
}