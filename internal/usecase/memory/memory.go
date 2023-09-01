package memory

import (
	"L0_WB/internal/aggregate"
	"L0_WB/internal/usecase"
	"L0_WB/internal/usecase/postgresdb"
	"fmt"
	"sync"
)

type MemoryRepository struct {
	orders map[string]aggregate.Order
	db     *postgresdb.OrderRepo
	sync.Mutex
}

func New(db *postgresdb.OrderRepo) *MemoryRepository {
	memRep := MemoryRepository{
		orders: make(map[string]aggregate.Order),
		db:     db,
	}

	ordersArray, err := memRep.db.GetPgAllRepos()
	if err != nil {
		fmt.Printf("Error get all data from DB\n")
		return &memRep
	}
	
	for _, order := range ordersArray {
		memRep.Lock()
		memRep.orders[order.GetOrderID()] = order
		memRep.Unlock()
	}
	return &memRep
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

	fmt.Printf("Save data to memory Id: %s\n", o.GetOrderID())
	err := r.db.PutPgRep(o)
	if err != nil {
		return fmt.Errorf("Error save data to Postgres err: %w", err)
	}
	fmt.Printf("Save data to Postgres\n")
	return nil
}
