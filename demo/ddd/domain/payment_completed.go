package domain

import "time"

// PaymentCompleted 支払完了イベント
// 支払いが完了したときに発行され、出荷準備処理をトリガーする
type PaymentCompleted struct {
	orderID    OrderID
	customerID CustomerID
	amount     Money
	occurredAt time.Time
}

// NewPaymentCompleted PaymentCompletedイベントを生成する
func NewPaymentCompleted(orderID OrderID, customerID CustomerID, amount Money) PaymentCompleted {
	return PaymentCompleted{
		orderID:    orderID,
		customerID: customerID,
		amount:     amount,
		occurredAt: time.Now(),
	}
}

// OccurredAt イベント発生日時を返す
func (e PaymentCompleted) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID 集約のIDを返す
func (e PaymentCompleted) AggregateID() string {
	return string(e.orderID)
}

// EventType イベントの種類を返す
func (e PaymentCompleted) EventType() string {
	return "PaymentCompleted"
}

// OrderID 注文IDを返す
func (e PaymentCompleted) OrderID() OrderID {
	return e.orderID
}

// CustomerID 顧客IDを返す
func (e PaymentCompleted) CustomerID() CustomerID {
	return e.customerID
}

// Amount 支払金額を返す
func (e PaymentCompleted) Amount() Money {
	return e.amount
}
