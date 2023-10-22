package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/ports"
)

type Withdrawal struct {
	db *sql.DB
}

type WithdrawalTx struct {
	tx *sql.Tx
}

func NewWithdrawal(db *sql.DB) *Withdrawal {
	return &Withdrawal{
		db: db,
	}
}

func (w *Withdrawal) Begin(ctx context.Context) (ports.WithdrawalTx, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &WithdrawalTx{tx: tx}, nil
}

func (w *Withdrawal) ListWithdrawals(
	ctx context.Context, user entities.Login,
) (entities.Withdrawals, error) {
	var result entities.Withdrawals

	rows, err := w.db.QueryContext(ctx,
		`SELECT id, upload_time, value FROM withdrawals WHERE username = $1`,
		user)
	if err != nil {
		return result, fmt.Errorf("failed to select %w", err)
	}

	var withdrawal entities.Withdrawal
	for rows.Next() {
		err := rows.Scan(
			&withdrawal.Order,
			&withdrawal.ProcessedAt,
			&withdrawal.Sum,
		)
		if err != nil {
			return result, fmt.Errorf("failed to scan: %w", err)
		}

		result = append(result, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("failed to finalize rows: %w", err)
	}

	return result, nil
}

func (w *WithdrawalTx) Commit() error {
	return w.tx.Commit()
}

func (w *WithdrawalTx) Rollback() error {
	return w.tx.Rollback()
}

func (w *WithdrawalTx) GetCurrentBalance(
	ctx context.Context,
	user entities.Login,
) (entities.Currency, error) {
	var result entities.Currency

	row := w.tx.QueryRowContext(ctx,
		`SELECT balance FROM users WHERE username = $1`,
		user)

	if err := row.Scan(&result); err != nil {
		return result, fmt.Errorf("failed to scan balance: %w", err)
	}

	return result, nil
}

func (w *WithdrawalTx) DecreaseBalance(
	ctx context.Context,
	user entities.Login,
	amount entities.Currency,
) error {
	_, err := w.tx.ExecContext(ctx,
		`UPDATE users SET balance = balance - $1 WHERE username = $2`,
		amount, user)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func (w *WithdrawalTx) IncreaseWithdrawn(
	ctx context.Context,
	user entities.Login,
	amount entities.Currency,
) error {
	_, err := w.tx.ExecContext(ctx,
		`UPDATE users SET withdrawn = withdrawn + $1 WHERE username = $2`,
		amount, user)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
