package domain

import "errors"

type ProductID struct {
	value string
}

func NewProductID(value string) (ProductID, error) {
	if value == "" {
		return ProductID{}, errors.New("product id cannot be empty")
	}

	return ProductID{value: value}, nil
}

func (pi ProductID) Value() string {
	return pi.value
}
