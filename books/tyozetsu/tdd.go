package tyozetsu

import (
	"fmt"
)

type Math struct{}
func (math *Math) min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
func (math *Math) max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

type MathUtil struct{}
func (mathUtil *MathUtil)saturate(value, minValue, maxValue int) int{
	math := &Math{}
	return math.min(math.max(value, minValue), maxValue)
}

func TddExec() {
	mathUtinl := &MathUtil{}
	fmt.Println(mathUtinl.saturate(1, 2, 3))
}
