package application

import "context"

// PaymentResult 支払い結果（ドメインの言葉で表現）
type PaymentResult struct {
	Success       bool
	TransactionID string
	ErrorMessage  string
}

// PaymentGateway 外部決済サービスとの連携インターフェース（腐敗防止層）
// 外部サービスのモデルを自ドメインのモデルに変換する役割を持つ
type PaymentGateway interface {
	// Charge 決済を実行する
	// orderID: 注文ID
	// amount: 支払金額
	// currency: 通貨
	Charge(ctx context.Context, orderID string, amount int, currency string) (*PaymentResult, error)

	// Refund 返金を実行する
	Refund(ctx context.Context, transactionID string, amount int, currency string) (*PaymentResult, error)
}
