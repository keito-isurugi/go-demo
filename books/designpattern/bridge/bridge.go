package bridge

type DisplayImpl interface {
	rawOpen()
	rawPrint()
	rawClose()
}

type Display struct {
	impl DisplayImpl
}
func newDisplay(di DisplayImpl) *Display{
	return &Display{impl: di}
}
func (d *Display) open() {
	d.impl.rawOpen()
}
func (d *Display) print() {
	d.impl.rawPrint()
}
func (d *Display) close() {
	d.impl.rawClose()
}

type CountDisplay struct {
	Display *Display
}
func newCountDisplay(impl DisplayImpl) *CountDisplay {
	return &CountDisplay{
		Display: newDisplay(impl),
	}
}
func (cd *CountDisplay) multiDisplay(times int) {
	cd.Display.open()
	for i := 0; i < times; i++ {
		cd.Display.print()
	}
	cd.Display.close()
}