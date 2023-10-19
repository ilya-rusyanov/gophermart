package storage

import (
	"context"
	"database/sql"
	"fmt"

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
	result := entities.NewUnfinishedOrders()

	rows, err := a.db.QueryContext(ctx,
		`SELECT id, state FROM orders
WHERE state != "INVALID" AND state != "PROCESSED"`)
	if err != nil {
		return result, fmt.Errorf("failed to select order states: %w", err)
	}
	defer rows.Close()

	var (
		id    entities.OrderID
		state entities.OrderStatus
	)
	for rows.Next() {
		err := rows.Scan(&id, &state)
		if err != nil {
			return result, fmt.Errorf("failed to scan order state: %w", err)
		}
		result[id] = state
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("failed to finalize rows: %w", err)
	}

	return result, nil
}

func (a *Accrual) UpdateOrderState(
	ctx context.Context,
	orderID entities.OrderID,
	nextStatus entities.OrderStatus,
	value *float64,
) error {
	_, err := a.db.ExecContext(ctx,
		`UPDATE orders SET state = $1, value = $2 WHERE id = $3`,
		nextStatus, value, orderID)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}
