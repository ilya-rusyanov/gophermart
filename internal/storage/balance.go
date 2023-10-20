package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Balance struct {
	db *sql.DB
}

func NewBalance(db *sql.DB) *Balance {
	return &Balance{
		db: db,
	}
}

func (b *Balance) Show(ctx context.Context, user entities.Login) (
	entities.Balance, error,
) {
	var result entities.Balance

	row := b.db.QueryRowContext(ctx,
		`SELECT balance, withdrawn FROM users
WHERE username = $1`, user)

	err := row.Scan(&result.Current, &result.Withdrawn)
	if err != nil {
		return result, fmt.Errorf("failed to scan row: %w", err)
	}

	return result, nil
}
