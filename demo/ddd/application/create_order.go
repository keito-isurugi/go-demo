package application

import (
	"context"
	"ddd/domain"
	"errors"

	"github.com/google/uuid"
)

// CreateOrderUseCase 注文作成ユースケース
type CreateOrderUseCase struct {
	orderRepo domain.OrderRepository
}

// NewCreateOrderUseCase CreateOrderUseCaseを生成する
func NewCreateOrderUseCase(orderRepo domain.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo: orderRepo,
	}
}

// CreateOrderInput 注文作成の入力DTO
type CreateOrderInput struct {
	CustomerID string
	Items      []CreateOrderItemInput
}

// CreateOrderItemInput 注文明細の入力DTO
type CreateOrderItemInput struct {
	ProductID   string
	ProductName string
	Price       int
	Currency    string
	Quantity    int
}

// CreateOrderOutput 注文作成の出力DTO
type CreateOrderOutput struct {
	OrderID string
	Success bool
	Message string
}

// Execute 注文作成を実行する
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// 入力バリデーション
	if input.CustomerID == "" {
		return nil, errors.New("customer id is required")
	}

	if len(input.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// OrderIDを生成
	orderID := domain.OrderID(uuid.New().String())
	customerID := domain.CustomerID(input.CustomerID)

	// 注文を作成（Draft状態）
	order, err := domain.NewOrder(orderID, customerID, nil, domain.Draft)
	if err != nil {
		return nil, err
	}

	// 注文明細を追加
	for _, item := range input.Items {
		// 値オブジェクトを生成
		productID, err := domain.NewProductID(item.ProductID)
		if err != nil {
			return nil, err
		}

		currency := domain.Currency(item.Currency)
		price, err := domain.NewMoney(item.Price, currency)
		if err != nil {
			return nil, err
		}

		quantity, err := domain.NewQuantity(item.Quantity)
		if err != nil {
			return nil, err
		}

		orderLine, err := domain.NewOrderLine(
			productID,
			domain.ProductName(item.ProductName),
			price,
			quantity,
		)
		if err != nil {
			return nil, err
		}

		// 注文に明細を追加（不変条件のチェックはOrder内で行われる）
		if err := order.AddLine(orderLine); err != nil {
			return nil, err
		}
	}

	// 永続化
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return &CreateOrderOutput{
		OrderID: string(orderID),
		Success: true,
		Message: "Order created successfully",
	}, nil
}
