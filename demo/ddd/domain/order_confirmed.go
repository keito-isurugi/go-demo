package domain

import "time"

// OrderConfirmed 注文確定イベント
// 注文が確定されたときに発行され、在庫引当処理をトリガーする
type OrderConfirmed struct {
	orderID     OrderID
	customerID  CustomerID
	totalAmount Money
	occurredAt  time.Time
}

// NewOrderConfirmed OrderConfirmedイベントを生成する
func NewOrderConfirmed(orderID OrderID, customerID CustomerID, totalAmount Money) OrderConfirmed {
	return OrderConfirmed{
		orderID:     orderID,
		customerID:  customerID,
		totalAmount: totalAmount,
		occurredAt:  time.Now(),
	}
}

// OccurredAt イベント発生日時を返す
func (e OrderConfirmed) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID 集約のIDを返す
func (e OrderConfirmed) AggregateID() string {
	return string(e.orderID)
}

// EventType イベントの種類を返す
func (e OrderConfirmed) EventType() string {
	return "OrderConfirmed"
}

// OrderID 注文IDを返す
func (e OrderConfirmed) OrderID() OrderID {
	return e.orderID
}

// CustomerID 顧客IDを返す
func (e OrderConfirmed) CustomerID() CustomerID {
	return e.customerID
}

// TotalAmount 合計金額を返す
func (e OrderConfirmed) TotalAmount() Money {
	return e.totalAmount
}
