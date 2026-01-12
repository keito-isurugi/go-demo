package application

import (
	"context"
	"ddd/domain"
	"errors"
)

// ProcessPaymentUseCase 支払処理ユースケース
type ProcessPaymentUseCase struct {
	orderRepo      domain.OrderRepository
	paymentGateway PaymentGateway
	eventPublisher EventPublisher
}

// NewProcessPaymentUseCase ProcessPaymentUseCaseを生成する
func NewProcessPaymentUseCase(
	orderRepo domain.OrderRepository,
	paymentGateway PaymentGateway,
	eventPublisher EventPublisher,
) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		orderRepo:      orderRepo,
		paymentGateway: paymentGateway,
		eventPublisher: eventPublisher,
	}
}

// ProcessPaymentInput 支払処理の入力DTO
type ProcessPaymentInput struct {
	OrderID string
}

// ProcessPaymentOutput 支払処理の出力DTO
type ProcessPaymentOutput struct {
	Success       bool
	TransactionID string
	Message       string
}

// Execute 支払処理を実行する
func (uc *ProcessPaymentUseCase) Execute(ctx context.Context, input ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	// 1. 注文を取得
	order, err := uc.orderRepo.FindByID(ctx, domain.OrderID(input.OrderID))
	if err != nil {
		return nil, err
	}

	// 2. 注文がConfirmed状態か確認
	if order.Status() != domain.Confirmed {
		return nil, errors.New("order must be in confirmed status to process payment")
	}

	// 3. 合計金額を取得
	total, err := order.Total()
	if err != nil {
		return nil, err
	}

	// 4. 外部決済サービスで決済を実行（腐敗防止層経由）
	paymentResult, err := uc.paymentGateway.Charge(
		ctx,
		input.OrderID,
		total.Amount(),
		string(total.Currency()),
	)
	if err != nil {
		return nil, err
	}

	if !paymentResult.Success {
		return &ProcessPaymentOutput{
			Success: false,
			Message: paymentResult.ErrorMessage,
		}, errors.New("payment failed: " + paymentResult.ErrorMessage)
	}

	// 5. 注文を支払済みに変更
	if err := order.Pay(); err != nil {
		// 支払い成功後にステータス変更に失敗した場合は返金が必要
		// TODO: 補償トランザクションの実装
		return nil, err
	}

	// 6. 永続化
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	// 7. ドメインイベントを発行
	for _, event := range order.DomainEvents() {
		if err := uc.eventPublisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	order.ClearDomainEvents()

	return &ProcessPaymentOutput{
		Success:       true,
		TransactionID: paymentResult.TransactionID,
		Message:       "Payment processed successfully",
	}, nil
}
