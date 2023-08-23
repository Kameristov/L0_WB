package usecase

import (
	"L0_EVRONE/internal/aggregate"
	"context"
	"fmt"
)

type OrderUseCase struct {
	repo OrderRepository
}

func New(r OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: r}
}

func (s *OrderUseCase) Set(ctx context.Context, order aggregate.Order) error {
	err := s.repo.PutRep( order)
	if err != nil {
		return fmt.Errorf("OrderUseCase - Set - s.repo.Put: %w", err)
	}
	return nil
}

func (s *OrderUseCase) Get(ctx context.Context, id string) (aggregate.Order, error) {
	ord, err := s.repo.GetRep( id)
	if err != nil {
		return aggregate.Order{}, fmt.Errorf("OrderUseCase - Get - s.repo.Get: %w", err)
	}
	return ord, nil
}
