package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Accrual struct {
	db *sql.DB
}

func NewAccrual(db *sql.DB) *Accrual {
	return &Accrual{
		db: db,
	}
}

func (a *Accrual) GetUnfinishedOrdersStates(ctx context.Context) (
	entities.UnfinishedOrders, error,
) {
	var result entities.UnfinishedOrders
	return result, errors.New("TODO")
}

func (a *Accrual) UpdateOrderState(
	ctx context.Context, orderID entities.OrderID, nextStatus entities.OrderStatus,
) error {
	return errors.New("TODO")
}
