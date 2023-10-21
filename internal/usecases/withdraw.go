package usecases

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Withdraw struct {
}

func NewWithdraw() *Withdraw {
	return &Withdraw{}
}

func (w *Withdraw) Withdraw(
	ctx context.Context, request entities.WithdrawalRequest,
) error {
	return errors.New("TODO")
}
