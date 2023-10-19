package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type AccrualStorage interface {
	GetUnfinishedOrdersStates(context.Context) (
		entities.UnfinishedOrders, error)
	UpdateOrderState(
		ctx context.Context,
		orderID entities.OrderID,
		status entities.OrderStatus,
		accrual *float64,
	) error
}

type AccrualService interface {
	GetStateOfOrder(context.Context, entities.OrderID) (
		status entities.OrderStatus, value float64, err error,
	)
}

type FeedAccrual struct {
	close   chan struct{}
	ticker  *time.Ticker
	errors  chan error
	storage AccrualStorage
	service AccrualService
	logger  Logger
}

func NewFeedAccrual(
	logger Logger, storage AccrualStorage, service AccrualService,
) *FeedAccrual {
	return &FeedAccrual{
		close:   make(chan struct{}),
		errors:  make(chan error, 1),
		storage: storage,
		service: service,
		logger:  logger,
	}
}

func (f *FeedAccrual) Run(ctx context.Context, basePeriod time.Duration) <-chan error {
	f.ticker = time.NewTicker(basePeriod)

	go func() {
		defer f.ticker.Stop()
		defer close(f.errors)

	mainLoop:
		for {
			select {
			case <-f.close:
				break mainLoop
			case <-f.ticker.C:
				err := f.reviseOrders(ctx)
				if err != nil {
					f.errors <- err
				}
			}
		}
	}()

	return f.errors
}

func (f *FeedAccrual) Close() {
	close(f.close)
}

func (f *FeedAccrual) reviseOrders(ctx context.Context) error {
	unfinishedOrders, err :=
		f.storage.GetUnfinishedOrdersStates(ctx)
	if err != nil {
		return err
	}

	if len(unfinishedOrders) > 0 {
		f.logger.Infof("found %d unchecked orders", len(unfinishedOrders))
	}

	for order, state := range unfinishedOrders {
		nextState, value, err := f.service.GetStateOfOrder(ctx, order)

		var delay *entities.AccrualTooManyRequestsError

		switch {
		case errors.Is(err, entities.ErrAccrualOrderIsNotRegistered):
			return fmt.Errorf(
				"order %d is not registered in accrual: %w", order, err)
		case errors.As(err, &delay):
			f.ticker.Reset(delay.Period)
			f.logger.Infof("ticker reset to %v", delay.Period)
			return nil
		case err != nil:
			return fmt.Errorf("unexpected error from accrual: %w", err)
		}

		if state != nextState {
			f.logger.Infof("order %d changed state from %q to %q",
				order, state, nextState)

			var accrual *float64

			if nextState == entities.OrderStatusProcessed {
				f.logger.Infof("order %d value will be %v",
					order, value)
				accrual = &value
			}

			err := f.storage.UpdateOrderState(ctx, order, nextState, accrual)
			if err != nil {
				return fmt.Errorf("failed to update order state: %w", err)
			}
			f.logger.Infof("order %d state updated successfully",
				order)
		}
	}

	return nil
}
