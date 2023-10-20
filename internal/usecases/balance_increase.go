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
	errors   chan error
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
		errors:   make(chan error),
	}
}

func (i *BalanceIncrease) Run(ctx context.Context) <-chan error {
	return i.errors
}
