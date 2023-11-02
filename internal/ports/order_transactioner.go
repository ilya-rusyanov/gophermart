package ports

import (
	"context"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type CreateOrderTransactioner interface {
	Begin(context.Context) (CreateOrderTransaction, error)
}

type CreateOrderTransaction interface {
	Commit() error
	Rollback() error
	FindUserForOrder(context.Context, entities.OrderID) (entities.Login, error)
	CreateOrder(context.Context, entities.CreateOrderRequest) error
}
