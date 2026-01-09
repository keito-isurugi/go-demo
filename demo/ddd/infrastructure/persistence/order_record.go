package persistence

import (
	"time"

	"ddd/domain"
)

type OrderRecord struct {
	ID string `db:"id"`
	CustomerID string `db:"customer_id"`
	Status string `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func toOrderRecord(order domain.Order) OrderRecord {
	orderID := string(order.ID())
	customerID := string(order.CustomerID())
	status := string(order.Status())

	return OrderRecord{
		ID: orderID,
		CustomerID: customerID,
		Status: status,
	}
}

func toOrderDomain(record OrderRecord) (domain.Order, error) {
	// return domain.NewOrder(), nil
	return domain.Order{}, nil
}