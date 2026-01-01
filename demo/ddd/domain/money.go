package domain

import "errors"

type Currensy string

const (
	JPY Currensy = "JPY"
	USD Currensy = "USD"
)

type Money struct {
	amount int
	currensy Currensy
}

func NewMoney(amount int, currensy Currensy) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount must be non-negative")
	}

	return Money{
		amount: amount,
		currensy: currensy,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currensy != other.currensy {
		return Money{}, errors.New("currensy mismatch")
	}

	return NewMoney(other.amount, other.currensy)
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currensy() Currensy {
	return m.currensy
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currensy == other.currensy
}