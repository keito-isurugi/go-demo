package persistence

import (
	"context"

	"ddd/domain"
)

type OrderRepositoryImpl struct{}

func NewOrderRepositoryImpl() domain.OrderRepository {
	return &OrderRepositoryImpl{}
}

func (o *OrderRepositoryImpl) Save(ctx context.Context, order domain.Order) error {
	// Order保存処理
	return nil
}

func (o *OrderRepositoryImpl) FindByID(ctx context.Context, id domain.OrderID) (domain.Order, error) {
	// Order取得処理
	return domain.Order{}, nil
}

func (o *OrderRepositoryImpl) FindByCustomerID(ctx context.Context, customerID domain.CustomerID) ([]domain.Order, error) {
	// Order一覧取得処理
	return nil, nil
}

func (o *OrderRepositoryImpl) Delete(ctx context.Context, id domain.OrderID) error {
	// Order削除処理
	return nil
}
