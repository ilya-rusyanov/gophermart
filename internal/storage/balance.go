package storage

import (
	"context"
	"database/sql"
	"errors"

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
	return result, errors.New("TODO")
}
