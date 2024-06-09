package builder

import "fmt"

type House struct {
	windowType string
	doorType string
	floor int
}

type IBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() House
}

func getBilder(builderType string) IBuilder {
	if builderType == "normal" {
		return newNormalBuider()
	}

	if builderType == "igloo" {
		return newIglooBuider()
	}

	return nil
}

type NormalBuidler struct {
	windowType string
	doorType string
	floor int
}
func newNormalBuider() *NormalBuidler {
	return &NormalBuidler{}
}
func (b *NormalBuidler) setWindowType() {
	b.windowType = "wooden window"
}
func (b *NormalBuidler) setDoorType() {
	b.doorType = "wooden window"
}
func (b *NormalBuidler) setNumFloor() {
	b.floor = 2
}
func (b *NormalBuidler) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType: b.doorType,
		floor: b.floor,
	}
}

type IglooBuidler struct {
	windowType string
	doorType string
	floor int
}
func newIglooBuider() *IglooBuidler {
	return &IglooBuidler{}
}
func (b *IglooBuidler) setWindowType() {
	b.windowType = "igloo window"
}
func (b *IglooBuidler) setDoorType() {
	b.doorType = "igloo window"
}
func (b *IglooBuidler) setNumFloor() {
	b.floor = 10
}
func (b *IglooBuidler) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType: b.doorType,
		floor: b.floor,
	}
}

type Director struct {
	builder IBuilder
}
func newDirector(b IBuilder) *Director{
	return &Director{
		builder: b,
	}
}
func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}
func (d *Director) buildHouse() House {
	d.builder.setDoorType()
	d.builder.setWindowType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}

func Exec() {
	normalBuidler := getBilder("normal")
	iglooBuidler := getBilder("igloo")

	director := newDirector(normalBuidler)
	normalHouse := director.buildHouse()

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.doorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.windowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.floor)

	director.setBuilder(iglooBuidler)
	iglooHouse := director.buildHouse()

	fmt.Printf("\nIgloo House Door Type: %s\n", iglooHouse.doorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.windowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.floor)
}