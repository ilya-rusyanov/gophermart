package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/ports"
)

type Order struct {
	db *sql.DB
}

type CreateOrderTransaction struct {
	tx *sql.Tx
}

func NewOrder(db *sql.DB) *Order {
	return &Order{
		db: db,
	}
}

func (o *Order) Begin(ctx context.Context) (ports.CreateOrderTransaction, error) {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &CreateOrderTransaction{tx: tx}, nil
}

func (t *CreateOrderTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *CreateOrderTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *CreateOrderTransaction) FindUserForOrder(
	ctx context.Context, order entities.OrderID,
) (entities.Login, error) {
	var login entities.Login
	return login, errors.New("TODO")
}

func (t *CreateOrderTransaction) CreateOrder(
	ctx context.Context, req entities.CreateOrderRequest,
) error {
	return errors.New("TODO")
}
