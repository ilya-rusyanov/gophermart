package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/ports"
)

type Order struct {
	db     *sql.DB
	logger Logger
}

type CreateOrderTransaction struct {
	tx     *sql.Tx
	logger Logger
}

func NewOrder(db *sql.DB, logger Logger) *Order {
	return &Order{
		db:     db,
		logger: logger,
	}
}

func (o *Order) Begin(ctx context.Context) (ports.CreateOrderTransaction, error) {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &CreateOrderTransaction{
		tx:     tx,
		logger: o.logger,
	}, nil
}

func (o *Order) ListOrders(
	ctx context.Context, request entities.ListOrdersRequest,
) (entities.OrderList, error) {
	var result entities.OrderList

	rows, err := o.db.QueryContext(ctx,
		`SELECT id, state, value, upload_time FROM orders
WHERE username = $1`, request.Login)
	if err != nil {
		return result, fmt.Errorf(
			"failed to select orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entities.Order
		err := rows.Scan(
			&order.ID, &order.Status, &order.Accrual, &order.UploadedAt,
		)
		if err != nil {
			return result, fmt.Errorf(
				"failed to scan single order: %w", err)
		}

		result = append(result, order)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf(
			"failure to finalize orders request: %w", err)
	}

	return result, nil
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

	err := t.tx.QueryRowContext(ctx,
		`SELECT username FROM orders WHERE id = $1`, order).Scan(&login)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return login, entities.ErrNotFound
	case err != nil:
		return login, fmt.Errorf("unexpected SQL error: %w", err)
	}

	return login, nil
}

func (t *CreateOrderTransaction) CreateOrder(
	ctx context.Context, req entities.CreateOrderRequest,
) error {
	t.logger.Debug("going to insert into orders table")
	_, err := t.tx.ExecContext(ctx,
		`INSERT INTO orders (id, username, upload_time, state)
VALUES ($1, $2, $3, $4)`,
		req.ID, req.User, req.Time, entities.OrderStatusNew)
	if err != nil {
		t.logger.Debug("insert failure")
		return fmt.Errorf("failed to insert order: %w", err)
	}
	t.logger.Debug("insert success")

	return nil
}
