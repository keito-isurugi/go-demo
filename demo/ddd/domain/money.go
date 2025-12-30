package domain

import "errors"

type Currensy string

const (
	JPY Currensy = "JPY"
	USD Currensy = "USD"
)

type Money struct {
	ammount int
	currensy Currensy
}

func NewMoney(ammount int, currensy Currensy) (Money, error) {
	if ammount < 0 {
		return Money{}, errors.New("ammount must be non-negative")
	}

	return Money{
		ammount: ammount,
		currensy: currensy,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currensy != other.currensy {
		return Money{}, errors.New("currensy mismatch")
	}

	return NewMoney(other.ammount, other.currensy)
}

func (m Money) Ammount() int {
	return m.ammount
}

func (m Money) Currensy() Currensy {
	return m.currensy
}

func (m Money) Equals(other Money) bool {
	return m.ammount == other.ammount && m.currensy == other.currensy
}