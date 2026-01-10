package domain

import "errors"

type Available int
type Reserved int

type Stock struct {
	productID    ProductID
	available    Available
	reserved     Reserved
	domainEvents []DomainEvent
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
		reserved:  reserved,
	}, nil
}

func (s Stock) Reserve(quantity int) (Stock, error) {
	if s.available < Available(quantity) {
		return Stock{}, errors.New("insufficient available stock")
	}

	newAvailable := s.available - Available(quantity)
	newReserved := s.reserved + Reserved(quantity)

	newStock, err := NewStock(s.productID, newAvailable, newReserved)
	if err != nil {
		return Stock{}, err
	}

	// 在庫が閾値を下回った場合、StockDepletedイベントを生成
	if newAvailable < Available(StockDepletedThreshold) {
		event := NewStockDepleted(s.productID, newAvailable)
		newStock.domainEvents = append(newStock.domainEvents, event)
	}

	return newStock, nil
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

// DomainEvents 保持しているドメインイベントを返す
func (s Stock) DomainEvents() []DomainEvent {
	return s.domainEvents
}

// ClearDomainEvents ドメインイベントをクリアする
func (s *Stock) ClearDomainEvents() {
	s.domainEvents = nil
}
