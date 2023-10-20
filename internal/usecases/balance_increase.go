package usecases

import (
	"context"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type BalanceIncreaseStorage interface {
	IncreaseBalance(context.Context, entities.Login, entities.Currency) error
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

		for {
			select {
			case <-ctx.Done():
				return
			case order := <-i.ordersCh:
				err := i.increase(ctx, order)
				if err != nil {
					errors <- fmt.Errorf(
						"storage failed to increase balance: %w",
						err)
				}
			}
		}
	}()

	return errors
}

func (i *BalanceIncrease) increase(
	ctx context.Context, order entities.Order,
) error {
	if order.Accrual == nil {
		return nil
	}

	err := i.storage.IncreaseBalance(
		ctx, order.User, *order.Accrual,
	)
	if err != nil {
		return fmt.Errorf("storage failed to increase balance: %w", err)
	}
	i.logger.Infof(
		"user %q balance increased by %v",
		order.User, order.Accrual)
	return nil
}
