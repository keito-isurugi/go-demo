package domain

import "errors"

type Available int
type Reserved int

type Stock struct {
	productID ProductID
	available Available
	reserved Reserved
}

func NewStock(productID ProductID, available Available, reserved Reserved) (Stock, error) {
	if productID.value == "" {
		return Stock{}, errors.New("product id cannot be empty")
	}

	if available < 0 {
		return Stock{}, errors.New("available must be non-negative")
	}

	if reserved < 0 {
		return Stock{}, errors.New("reserved must be non-negative")
	}

	return Stock{
		productID: productID,
		available: available,
		reserved: reserved,
	}, nil
}

func (s Stock) Reserve(quantity int) (Stock, error) {
	if s.available < Available(quantity) {
		return Stock{}, errors.New("insufficient available stock") 
	}

	newAvailable := s.available - Available(quantity)
	newReserved := s.reserved + Reserved(quantity)

	return NewStock(s.productID, newAvailable, newReserved)
}

func (s Stock) Release(quantity int) (Stock, error) {
	if s.reserved < Reserved(quantity) {
		return Stock{}, errors.New("insufficient reserved stock") 
	}

	newAvailable := s.available + Available(quantity)
	newReserved := s.reserved - Reserved(quantity)

	return NewStock(s.productID, newAvailable, newReserved)
}

func (s Stock) CanReserve(quantity int) bool {
	return s.available >= Available(quantity)
}

func (s Stock) ProductID() ProductID {
	return s.productID
}

func (s Stock) Available() Available {
	return s.available
}

func (s Stock) Reserved() Reserved {
	return s.reserved
}