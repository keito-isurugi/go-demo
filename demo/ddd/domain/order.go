package domain

import "errors"

type OrderID string
type CustomerID string
type OrderStatus string

const (
	Draft     OrderStatus = "Draft"
	Confirmed OrderStatus = "Confirmed"
	Paid      OrderStatus = "Paid"
	Shipped   OrderStatus = "Shipped"
	Delivered OrderStatus = "Delivered"
	Cancelled OrderStatus = "Cancelled"
)

const MaxOrderAmount = 1000000

type Order struct {
	id           OrderID
	customerID   CustomerID
	lines        []OrderLine
	status       OrderStatus
	domainEvents []DomainEvent
}

func NewOrder(id OrderID, customerID CustomerID, lines []OrderLine, status OrderStatus) (Order, error) {
	if id == "" {
		return Order{}, errors.New("order id cannot be empty")
	}

	if customerID == "" {
		return Order{}, errors.New("customer id cannot be empty")
	}

	return Order{
		id:         id,
		customerID: customerID,
		lines:      lines,
		status:     status,
	}, nil
}

func (o *Order) AddLine(line OrderLine) error {
	if o.status != Draft {
		return errors.New("can only add lines to draft orders")
	}

	// 重複商品チェック
	for _, existingLine := range o.lines {
		if existingLine.ProductID() == line.ProductID() {
			return errors.New("product already exists in order")
		}
	}

	// 現在の合計を取得
	currentTotal, err := o.Total()
	if err != nil {
		return err
	}

	// 新しい明細の小計を取得
	lineSubtotal, err := line.Subtotal()
	if err != nil {
		return err
	}

	// 新しい合計を計算
	newTotal, err := currentTotal.Add(lineSubtotal)
	if err != nil {
		return err
	}

	// 上限チェック（100万円）
	maxAmount, _ := NewMoney(MaxOrderAmount, JPY)
	isGreater, err := newTotal.IsGreaterThan(maxAmount)
	if err != nil {
		return err
	}
	if isGreater {
		return errors.New("order total exceeds maximum allowed amount")
	}

	o.lines = append(o.lines, line)
	return nil
}

func (o *Order) RemoveLine(productID ProductID) error {
	if o.status != Draft {
		return errors.New("can only add lines to draft orders")
	}

	for i, l := range o.lines {
		if l.ProductID() == productID {
			o.lines = append(o.lines[:i], o.lines[i+1:]...)
			return nil
		}
	}

	return errors.New("order line not found")
}

func (o *Order) Confirm() error {
	if o.status != Draft {
		return errors.New("can only confirm draft orders")
	}

	if len(o.lines) == 0 {
		return errors.New("cannot confirm order with no lines")
	}

	o.status = Confirmed

	// OrderConfirmedイベントを生成
	total, err := o.Total()
	if err != nil {
		return err
	}
	event := NewOrderConfirmed(o.id, o.customerID, total)
	o.domainEvents = append(o.domainEvents, event)

	return nil
}

func (o *Order) Pay() error {
	if o.status != Confirmed {
		return errors.New("can only pay confirmed orders")
	}
	o.status = Paid

	// PaymentCompletedイベントを生成
	total, err := o.Total()
	if err != nil {
		return err
	}
	event := NewPaymentCompleted(o.id, o.customerID, total)
	o.domainEvents = append(o.domainEvents, event)

	return nil
}

func (o *Order) Ship() error {
	if o.status != Paid {
		return errors.New("can only ship paid orders")
	}
	o.status = Shipped
	return nil
}

func (o *Order) Cancel() error {
	if o.status == Shipped || o.status == Delivered {
		return errors.New("cannot cancel shipped or delivered orders")
	}
	o.status = Cancelled
	return nil
}

func (o Order) Total() (Money, error) {
	if len(o.lines) == 0 {
		// デフォルト通貨でゼロを返す（または適切な通貨を決定）
		return NewMoney(0, JPY)
	}

	// 最初の明細の通貨を基準とする
	firstSubtotal, err := o.lines[0].Subtotal()
	if err != nil {
		return Money{}, err
	}

	total := firstSubtotal
	for i := 1; i < len(o.lines); i++ {
		subtotal, err := o.lines[i].Subtotal()
		if err != nil {
			return Money{}, err
		}
		total, err = total.Add(subtotal)
		if err != nil {
			return Money{}, err
		}
	}

	return total, nil
}

func (o Order) ID() OrderID {
	return o.id
}

func (o Order) CustomerID() CustomerID {
	return o.customerID
}

func (o Order) Line() []OrderLine {
	return o.lines
}

func (o Order) Status() OrderStatus {
	return o.status
}

// DomainEvents 保持しているドメインイベントを返す
func (o Order) DomainEvents() []DomainEvent {
	return o.domainEvents
}

// ClearDomainEvents ドメインイベントをクリアする
// イベント発行後にアプリケーション層から呼ばれる
func (o *Order) ClearDomainEvents() {
	o.domainEvents = nil
}
