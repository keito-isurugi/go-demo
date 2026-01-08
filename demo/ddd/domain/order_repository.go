package domain

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id OrderID) (Order, error)
	FindByCustomerID(ctx context.Context, customerID CustomerID) ([]Order, error)
	Delete(ctx context.Context, id OrderID) error
}