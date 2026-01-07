package domain

type AllocatedItem struct {
	ProductID ProductID
	Quantity  int
}

type AllocationResult struct {
	Success        bool
	AllocatedItems []AllocatedItem
	FailedItems    []ProductID
	Message        string
}

type StockAllocationService struct{}

func NewStockAllocationService() *StockAllocationService {
	return &StockAllocationService{}
}

func (s *StockAllocationService) Allocate(order Order, stocks map[ProductID]*Stock) (AllocationResult, error) {
	var allocatedItems []AllocatedItem
	var failedItems []ProductID
	var message string

	for _, line := range order.lines {
		stock, ok := stocks[line.ProductID()]
		if !ok {
			failedItems = append(failedItems, line.ProductID())
			continue
		}

		if !stock.CanReserve(line.quantity.value) {
			failedItems = append(failedItems, line.ProductID())
			continue
		}

		newStock, err := stock.Reserve(line.quantity.value)
		if err != nil {
			failedItems = append(failedItems, line.ProductID())
			continue
		}

		*stocks[line.ProductID()] = newStock
		allocatedItems = append(allocatedItems, AllocatedItem{
			ProductID: line.ProductID(),
			Quantity:  line.quantity.value,
		})
	}

	success := len(failedItems) == 0
	message = ""
	if !success {
		message = "some items could not be allocated"
	}

	return AllocationResult{
		Success:        success,
		AllocatedItems: allocatedItems,
		FailedItems:    failedItems,
		Message:        message,
	}, nil
}
