package accrual

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Adapter struct {
	addr string
}

func New(addr string) *Adapter {
	return &Adapter{
		addr: addr,
	}
}

func (a *Adapter) GetStateOfOrder(ctx context.Context, orderID entities.OrderID) (
	entities.OrderStatus, error,
) {
	var currentStatus entities.OrderStatus
	return currentStatus, errors.New("TODO")
}
