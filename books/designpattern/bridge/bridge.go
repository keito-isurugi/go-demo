package bridge

type DisplayImpl interface {
	rawOpen()
	rawPrint()
	rawClose()
}

type Display struct {
	impl DisplayImpl
}

// NewDisplay creates a new Display instance
func NewDisplay(di DisplayImpl) *Display {
	return &Display{impl: di}
}

// Open opens the display
func (d *Display) Open() {
	d.impl.rawOpen()
}

// Print prints to the display
func (d *Display) Print() {
	d.impl.rawPrint()
}

// Close closes the display
func (d *Display) Close() {
	d.impl.rawClose()
}

type CountDisplay struct {
	Display *Display
}

// NewCountDisplay creates a new CountDisplay instance
func NewCountDisplay(impl DisplayImpl) *CountDisplay {
	return &CountDisplay{
		Display: NewDisplay(impl),
	}
}

// MultiDisplay displays multiple times
func (cd *CountDisplay) MultiDisplay(times int) {
	cd.Display.Open()
	for i := 0; i < times; i++ {
		cd.Display.Print()
	}
	cd.Display.Close()
}