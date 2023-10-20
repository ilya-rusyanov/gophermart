package usecases

import (
	"context"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type BalanceIncreaseStorage interface {
}

type BalanceIncrease struct {
	logger   Logger
	ordersCh <-chan entities.Order
	storage  BalanceIncreaseStorage
}

func NewBalanceIncrease(
	logger Logger,
	ordersCh <-chan entities.Order,
	storage BalanceIncreaseStorage,
) *BalanceIncrease {
	return &BalanceIncrease{
		logger:   logger,
		ordersCh: ordersCh,
		storage:  storage,
	}
}

func (i *BalanceIncrease) Run(ctx context.Context) <-chan error {
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		select {
		case <-ctx.Done():
			return
		}
	}()

	return errors
}
