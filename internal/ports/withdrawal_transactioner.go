package ports

import (
	"context"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type WithdrawalTransactioner interface {
	Begin(context.Context) (WithdrawalTx, error)
}

type WithdrawalTx interface {
	Commit() error
	Rollback() error
	GetCurrentBalance(context.Context, entities.Login) (entities.Currency, error)
	DecreaseBalance(context.Context, entities.Login, entities.Currency) error
	IncreaseWithdrawn(context.Context, entities.Login, entities.Currency) error
}
