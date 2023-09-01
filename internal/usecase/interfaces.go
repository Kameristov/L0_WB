package usecase

import (
	"L0_WB/internal/aggregate"
	"context"
	"errors"
)

var (
	ErrOrderNotFound    = errors.New("the order was not found in the repository")
	ErrNilOrder         = errors.New("failed to add the order is nil")
	ErrFailedToAddOrder = errors.New("failed to add the order to the repository")
)

type OrderNats interface {
	Set(context.Context, aggregate.Order)  error
}
type OrderApi interface {
	Get(context.Context, string) (aggregate.Order, error)
}
type OrderRepository interface {
	PutRep(aggregate.Order) error
	GetRep(string) (aggregate.Order, error)
}

type OrderPostgresRepository interface {
	PutPgRep(aggregate.Order) error
	GetPgRep(string) (aggregate.Order, error)
	GetPgAllRep() ([]aggregate.Order, error)
}