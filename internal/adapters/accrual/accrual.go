package accrual

import (
	"context"
	"errors"
	"net/http"

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
	resp, err := http.Get(a.addr)
	return currentStatus, errors.New("TODO")
}
