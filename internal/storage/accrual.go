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

func (a *Accrual) GetUnfinishedOrders(ctx context.Context) (
	entities.OrderList, error,
) {
	var result entities.OrderList

	rows, err := a.db.QueryContext(ctx,
		`SELECT id, state, upload_time, username FROM orders
WHERE state != 'INVALID' AND state != 'PROCESSED'`)
	if err != nil {
		return result, fmt.Errorf("failed to select order states: %w", err)
	}
	defer rows.Close()

	var (
		order entities.Order
	)
	for rows.Next() {
		err := rows.Scan(
			&order.ID, &order.Status, &order.UploadedAt, &order.User)
		if err != nil {
			return result, fmt.Errorf("failed to scan order state: %w", err)
		}
		result = append(result, order)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("failed to finalize rows: %w", err)
	}

	return result, nil
}

func (a *Accrual) UpdateOrder(
	ctx context.Context,
	order entities.Order,
) error {
	_, err := a.db.ExecContext(ctx,
		`UPDATE orders SET state = $1, value = $2 WHERE id = $3`,
		order.Status, order.Accrual, order.ID)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}
