package application

import (
	"context"
	"ddd/domain"
	"errors"
)

// ConfirmOrderUseCase 注文確定ユースケース
type ConfirmOrderUseCase struct {
	orderRepo       domain.OrderRepository
	stockRepo       domain.StockRepository
	allocService    *domain.StockAllocationService
	eventPublisher  EventPublisher
}

// NewConfirmOrderUseCase ConfirmOrderUseCaseを生成する
func NewConfirmOrderUseCase(
	orderRepo domain.OrderRepository,
	stockRepo domain.StockRepository,
	allocService *domain.StockAllocationService,
	eventPublisher EventPublisher,
) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{
		orderRepo:      orderRepo,
		stockRepo:      stockRepo,
		allocService:   allocService,
		eventPublisher: eventPublisher,
	}
}

// ConfirmOrderInput 注文確定の入力DTO
type ConfirmOrderInput struct {
	OrderID string
}

// ConfirmOrderOutput 注文確定の出力DTO
type ConfirmOrderOutput struct {
	Success bool
	Message string
}

// Execute 注文確定を実行する
func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, input ConfirmOrderInput) (*ConfirmOrderOutput, error) {
	// 1. 注文を取得
	order, err := uc.orderRepo.FindByID(ctx, domain.OrderID(input.OrderID))
	if err != nil {
		return nil, err
	}

	// 2. 注文の商品IDリストを取得
	lines := order.Line()
	productIDs := make([]domain.ProductID, len(lines))
	for i, line := range lines {
		productIDs[i] = line.ProductID()
	}

	// 3. 在庫を取得
	stockList, err := uc.stockRepo.FindByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// mapに変換
	stocks := make(map[domain.ProductID]*domain.Stock)
	for i := range stockList {
		stocks[stockList[i].ProductID()] = &stockList[i]
	}

	// 4. 在庫を引き当て
	result, err := uc.allocService.Allocate(order, stocks)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return &ConfirmOrderOutput{
			Success: false,
			Message: result.Message,
		}, errors.New("stock allocation failed")
	}

	// 5. 注文を確定
	if err := order.Confirm(); err != nil {
		return nil, err
	}

	// 6. 永続化（注文と在庫の両方をSave）
	// TODO: 本来はトランザクションで囲む必要がある
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	for _, stock := range stocks {
		if err := uc.stockRepo.Save(ctx, *stock); err != nil {
			return nil, err
		}
	}

	// 7. ドメインイベントを発行
	for _, event := range order.DomainEvents() {
		if err := uc.eventPublisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	order.ClearDomainEvents()

	// 在庫のイベントも発行
	for _, stock := range stocks {
		for _, event := range stock.DomainEvents() {
			if err := uc.eventPublisher.Publish(ctx, event); err != nil {
				return nil, err
			}
		}
		stock.ClearDomainEvents()
	}

	return &ConfirmOrderOutput{
		Success: true,
		Message: "Order confirmed successfully",
	}, nil
}
