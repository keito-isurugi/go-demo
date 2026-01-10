package domain

import "time"

// StockDepletedThreshold 在庫枯渇の閾値
const StockDepletedThreshold = 10

// StockDepleted 在庫枯渇イベント
// 在庫が閾値を下回ったときに発行され、発注処理をトリガーする
type StockDepleted struct {
	productID  ProductID
	available  Available
	occurredAt time.Time
}

// NewStockDepleted StockDepletedイベントを生成する
func NewStockDepleted(productID ProductID, available Available) StockDepleted {
	return StockDepleted{
		productID:  productID,
		available:  available,
		occurredAt: time.Now(),
	}
}

// OccurredAt イベント発生日時を返す
func (e StockDepleted) OccurredAt() time.Time {
	return e.occurredAt
}

// AggregateID 集約のIDを返す
func (e StockDepleted) AggregateID() string {
	return e.productID.Value()
}

// EventType イベントの種類を返す
func (e StockDepleted) EventType() string {
	return "StockDepleted"
}

// ProductID 商品IDを返す
func (e StockDepleted) ProductID() ProductID {
	return e.productID
}

// Available 残りの利用可能在庫数を返す
func (e StockDepleted) Available() Available {
	return e.available
}
