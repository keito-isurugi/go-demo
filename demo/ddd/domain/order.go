package domain

import "errors"

type OrderID string
type CustomerID string
type OrderStatus string

const (
	Draft OrderStatus = "Draft"
	Confirmed OrderStatus = "Confirmed"
	Paid OrderStatus = "Paid"
	Shipped OrderStatus = "Shipped"
	Delivered OrderStatus = "Delivered"
	Cancelled OrderStatus = "Cancelled"
)

type Order struct {
	id OrderID
	customerID CustomerID
	lines []OrderLine
	status OrderStatus
}

func NewOrder(id OrderID, customerID CustomerID, lines []OrderLine, status OrderStatus) (Order, error) {
	if id == "" {
		return Order{}, errors.New("order id cannot be empty")
	}

	if customerID == "" {
		return Order{}, errors.New("customer id cannot be empty")
	}

	return Order{
		id: id,
		customerID: customerID,
		lines: lines,
		status: status,
	}, nil
}

func (o *Order) AddLine(line OrderLine) error {
	if o.status != Draft {
		return errors.New("can only add lines to draft orders")
	}

	o.lines = append(o.lines, line)
	return nil
}

func (o *Order) RemoveLine(productID ProductID) error {
	if o.status != Draft {
		return errors.New("can only add lines to draft orders")
	}

	for i, l := range o.lines {
		if l.ProductID() == productID {
			o.lines = append(o.lines[:i], o.lines[i+1:]...)
			return nil
		}
	}
	
	return errors.New("order line not found")
}

func (o *Order) Confirmed() error {
	if o.status != Draft {
		return errors.New("can only confirm draft orders")
	}

	if len(o.lines) == 0 {
		return errors.New("cannot confirm order with no lines")
	}

	o.status = Confirmed
	return nil
}

func (o *Order) Pay() error {
    if o.status != Confirmed {
        return errors.New("can only pay confirmed orders")
    }
    o.status = Paid
    return nil
}

func (o *Order) Ship() error {
    if o.status != Paid {
        return errors.New("can only ship paid orders")
    }
    o.status = Shipped
    return nil
}

func (o *Order) Cancel() error {
    if o.status == Shipped || o.status == Delivered {
        return errors.New("cannot cancel shipped or delivered orders")
    }
    o.status = Cancelled
    return nil
}

func (o Order) Total() int {
	var res int
	for _, l := range o.lines {
		res += l.Subtotal()
	}
	return res
}

func (o Order) ID() OrderID {
	return o.id
}

func (o Order) CustomerID() CustomerID {
	return o.customerID
}

func (o Order) Line() []OrderLine {
	return o.lines
}

func (o Order) Status() OrderStatus {
	return o.status
}