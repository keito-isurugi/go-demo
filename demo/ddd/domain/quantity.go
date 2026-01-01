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

func (q Quantity) Add(other Quantity) (Quantity, error) {
	return NewQuantity(q.value + other.value)
}

func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	return NewQuantity(q.value - other.value)
}

func (q Quantity) Multiply(multiplier int) (Quantity, error) {
	return NewQuantity(q.value * multiplier)
}

func (q Quantity) IsZero() bool {
	return q.value == 0
}

func (q Quantity) Value() int {
	return q.value
}
