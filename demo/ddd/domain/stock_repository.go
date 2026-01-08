package domain

import "context"

type StockRepository interface {
	Save(ctx context.Context, stock Stock) error
	FindByProductID(ctx context.Context, productID ProductID) (Stock, error)
	FindByProductIDs(ctx context.Context, productIDs []ProductID) ([]Stock, error)
	Delete(ctx context.Context, productID ProductID) error
}