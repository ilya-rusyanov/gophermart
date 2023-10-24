package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type AccrualStorage interface {
	GetUnfinishedOrders(context.Context) (
		entities.OrderList, error)
	UpdateOrder(ctx context.Context, order entities.Order) error
}

type AccrualService interface {
	GetStateOfOrder(ctx context.Context, orderID entities.OrderID) (
		status entities.OrderStatus, value entities.Currency, err error,
	)
}

type FeedAccrual struct {
	storage         AccrualStorage
	service         AccrualService
	logger          Logger
	processedOrders chan entities.Order
}

func NewFeedAccrual(
	logger Logger, storage AccrualStorage, service AccrualService,
) *FeedAccrual {
	return &FeedAccrual{
		processedOrders: make(chan entities.Order, 1),
		storage:         storage,
		service:         service,
		logger:          logger,
	}
}

func (f *FeedAccrual) Run(ctx context.Context, basePeriod time.Duration) (
	<-chan entities.Order, <-chan error,
) {
	errors := make(chan error, 1)

	ticker := time.NewTicker(basePeriod)

	go func() {
		defer ticker.Stop()
		defer close(errors)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := f.reviseOrders(ctx, ticker)
				if err != nil {
					errors <- err
				}
			}
		}
	}()

	return f.processedOrders, errors
}

func (f *FeedAccrual) reviseOrders(ctx context.Context, ticker *time.Ticker) error {
	unfinishedOrders, err :=
		f.storage.GetUnfinishedOrders(ctx)
	if err != nil {
		return err
	}

	if len(unfinishedOrders) > 0 {
		f.logger.Infof("found %d unchecked orders", len(unfinishedOrders))
	}

	for _, order := range unfinishedOrders {
		nextStatus, value, err := f.service.GetStateOfOrder(ctx, order.ID)

		var delay *entities.AccrualTooManyRequestsError

		switch {
		case errors.Is(err, entities.ErrAccrualOrderIsNotRegistered):
			return fmt.Errorf(
				"order %q is not registered in accrual: %w", order.ID, err)
		case errors.As(err, &delay):
			ticker.Reset(delay.Period)
			f.logger.Infof("ticker reset to %v", delay.Period)
			return nil
		case err != nil:
			return fmt.Errorf("unexpected error from accrual: %w", err)
		}

		if order.Status != nextStatus {
			update := order
			update.Status = nextStatus

			f.logger.Infof("order %q changed state from %q to %q",
				order.ID, order.Status, nextStatus)

			if nextStatus == entities.OrderStatusProcessed {
				f.logger.Infof("order %q value will be %v",
					order.ID, value)
				update.Accrual = &value
				f.processedOrders <- update
			}

			err := f.storage.UpdateOrder(ctx, update)
			if err != nil {
				return fmt.Errorf(
					"failed to update order state: %w", err)
			}
			f.logger.Infof("order %q state updated successfully",
				order.ID)
		}
	}

	return nil
}
