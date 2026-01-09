package persistence

import (
	"ddd/domain"
)

type OrderLineRecord struct {
	OrderID     string `db:"order_id"`
	ProductID   string `db:"product_id"`
	ProductName string `db:"product_name"`
	Amount      int64  `db:"amount"`
	Currency    string `db:"currency"`
	Quantity    int    `db:"quantity"`
}

func toOrderLineRecord(orderID domain.OrderID, line domain.OrderLine) OrderLineRecord {
	return OrderLineRecord{
		OrderID:     string(orderID),
		ProductID:   line.ProductID().Value(),
		ProductName: string(line.ProductName()),
		Amount:      int64(line.Price().Amount()),
		Currency:    string(line.Price().Currency()),
		Quantity:    line.Quantity().Value(),
	}
}

func toOrderLineDomain(record OrderLineRecord) (domain.OrderLine, error) {
	// return NewOrderLine(productID, productName, money, quantity), nil
	return domain.OrderLine{}, nil
}
