package domain

import "errors"

type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value < 0 {
		return Quantity{}, errors.New("quantity must be non-negative")
	}
	
	return Quantity{value: value}, nil
}

func (q Quantity) Add(value int) (Quantity, error) {
	return NewQuantity(q.value + value)
}

func (q Quantity) Subtract(value int) (Quantity, error) {
	return NewQuantity(q.value - value)
}

func (q Quantity) IsZero() bool {
	return q.value == 0
}

func (q Quantity) Value() int {
	return q.value
}