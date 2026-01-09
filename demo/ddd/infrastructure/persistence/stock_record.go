package persistence

import (
	"time"

	"ddd/domain"
)

type StockRecord struct {
	ProductID string `db:"product_id"`
	Available int `db:"available"`
	Reserved int `db:"reserved"`
	CreatedAt time.Time `db:"created_at"`
}

func toStockRecord(stock domain.Stock) StockRecord {
	return StockRecord{
		ProductID: stock.ProductID().Value(),
		Available: int(stock.Available()),
		Reserved: int(stock.Reserved()),
	}
}

func toStockDomain(record StockRecord) (domain.Stock, error) {
	// return domain.NewStock(), nil
	return domain.Stock{}, nil
}


