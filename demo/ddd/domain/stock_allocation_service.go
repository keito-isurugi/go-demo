package domain

// import "errors"

type AllocatedItem struct {
    ProductID ProductID
    Quantity  int
}

type AllocationResult struct {
	Success bool
	AllocatedItems []AllocatedItem
	FailedItems []ProductID
	Message string
}

func NewAllocationResult(order Order, stocs map[ProductID]*Stock) (AllocationResult, error) {
	var success bool
	var ai []AllocatedItem
	var fi []ProductID
	var message string


	for _, o := range order.lines {
		if s, ok := stocs[o.ProductID()]; ok {
			if !s.CanReserve(o.quantity.value) {
				fi = append(fi, o.ProductID())

				continue
			}
			s.Release(o.quantity.value)
			// ai = append(ai, AllocatedItem{ProductID: o.ProductID(), Quantity: o.quantity.value})
		}
	}


	return AllocationResult{
		Success: success,
		AllocatedItems: ai,
		FailedItems: fi,
		Message: message,
	}, nil
}