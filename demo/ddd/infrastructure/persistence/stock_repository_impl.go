package persistence

import (
	"context"

	"ddd/domain"
)

type StockRepositoryImpl struct{}

func NewStockRepositoryImpl() domain.StockRepository {
	return &StockRepositoryImpl{}
}

func (s *StockRepositoryImpl) Save(ctx context.Context, stock domain.Stock) error {
	// Stock保存処理
	return nil
}

func (s *StockRepositoryImpl) FindByProductID(ctx context.Context, productID domain.ProductID) (domain.Stock, error) {
	// Stock取得処理
	return domain.Stock{}, nil
}

func (s *StockRepositoryImpl) FindByProductIDs(ctx context.Context, productIDs []domain.ProductID) ([]domain.Stock, error) {
	// Stock一覧取得処理
	return []domain.Stock{}, nil
}

func (s *StockRepositoryImpl) Delete(ctx context.Context, productID domain.ProductID) error {
	// Stock削除処理
	return nil
}