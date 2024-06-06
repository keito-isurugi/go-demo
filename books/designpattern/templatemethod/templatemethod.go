package templatemethod

import (
	"fmt"
)

type IDisplay interface {
	open()
	print()
	close()
}

type AbstractDisplay struct {
	Display IDisplay
}
func NewAbstractDisplay(id IDisplay) *AbstractDisplay{
	return &AbstractDisplay{Display: id} 
}
func (d *AbstractDisplay) display() {
	d.Display.open()
	for i := 0; i < 5; i++ {
		d.Display.print()
	}
	d.Display.close()
}

type CharDisplay struct {
	number int
}
func NewCharDisplay(n int) CharDisplay{
	return CharDisplay{number: n}
}
func (cd CharDisplay) open() {
	fmt.Print("<<")
}
func (cd CharDisplay) print() {
	fmt.Print(cd.number)
}
func (cd CharDisplay) close() {
	fmt.Println(">>")
}


type StringDisplay struct {
	string string
	width int
}
func NewStringDisplay(s string) StringDisplay{
	return StringDisplay{
		string: s,
		width: len(s) + 1,
	}
}
func (cd StringDisplay) open() {
	cd.printLine()
}
func (cd StringDisplay) print() {
	fmt.Println("|", cd.string , "|")
}
func (cd StringDisplay) close() {
	cd.printLine()
}
func (cd StringDisplay) printLine() {
	fmt.Print("+")
	for i := 0; i <= cd.width; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func Exec() {
	cd := NewCharDisplay(9)
	d1 := NewAbstractDisplay(cd)
	d1.display()

	sd := NewStringDisplay("Hello, world")
	d2 := NewAbstractDisplay(sd)
	d2.display()
}