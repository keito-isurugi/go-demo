package domain

import "errors"

type ProductName string

type OrderLine struct {
	id ProductID
	productName ProductName
	price Money
	quantity Quantity
}

func NewOrderLine(id ProductID, productName ProductName, money Money, quantity Quantity) (OrderLine, error) {
	if productName == "" {
		return OrderLine{}, errors.New("product name cannot be empty")
	}

	if quantity.IsZero() {
		return OrderLine{}, errors.New("quantity must be greater than zero")
	}

	return OrderLine{
		id: id,
		productName: productName,
		price: money,
		quantity: quantity,
	}, nil
}

func (ol OrderLine) Subtotal() (Money, error) {
	amount := ol.price.Amount() * ol.quantity.Value()
	return NewMoney(amount, ol.price.currency)
}

func (ol OrderLine) ProductID() ProductID {
	return ol.id
}

func (ol OrderLine) ProductName() ProductName {
	return ol.productName
}

func (ol OrderLine) Price() Money {
	return ol.price
}

func (ol OrderLine) Quantity() Quantity {
	return ol.quantity
}
