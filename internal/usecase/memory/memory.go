package memory

import (
	"L0_EVRONE/internal/aggregate"
	"L0_EVRONE/internal/usecase"
	"fmt"
	"sync"
)

type MemoryRepository struct {
	orders map[string]aggregate.Order
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		orders: make(map[string]aggregate.Order),
	}
}

func (r *MemoryRepository) GetRep(orderId string) (aggregate.Order, error) {
	if order, ok := r.orders[orderId]; ok {
		return order, nil
	}
	return aggregate.Order{}, usecase.ErrOrderNotFound
}

func (r *MemoryRepository) PutRep(o aggregate.Order) error {
	if r.orders == nil {
		return usecase.ErrNilOrder
	}
	if _, ok := r.orders[o.GetOrderID()]; ok {
		return fmt.Errorf("order already exists: %w", usecase.ErrFailedToAddOrder)
	}

	r.Lock()
	r.orders[o.GetOrderID()] = o
	r.Unlock()
	return nil
}
