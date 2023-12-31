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
	IncreaseBalance(context.Context, entities.Login, entities.Currency) error
}

type AccrualService interface {
	GetStateOfOrder(ctx context.Context, orderID entities.OrderID) (
		status entities.OrderStatus, value entities.Currency, err error,
	)
}

type FeedAccrual struct {
	storage AccrualStorage
	service AccrualService
	logger  Logger
}

func NewFeedAccrual(
	logger Logger, storage AccrualStorage, service AccrualService,
) *FeedAccrual {
	return &FeedAccrual{
		storage: storage,
		service: service,
		logger:  logger,
	}
}

func (f *FeedAccrual) Run(ctx context.Context, basePeriod time.Duration) {
	ticker := time.NewTicker(basePeriod)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := f.reviseOrders(ctx)
			var delay *entities.AccrualTooManyRequestsError
			if errors.As(err, &delay) {
				ticker.Reset(delay.Period)
				f.logger.Infof("ticker reset to %v", delay.Period)
				break
			} else if err != nil {
				f.logger.Error(err)
			}

			ticker.Reset(basePeriod)
		}
	}
}

func (f *FeedAccrual) reviseOrders(ctx context.Context) error {
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

		switch {
		case errors.Is(err, entities.ErrAccrualOrderIsNotRegistered):
			return fmt.Errorf(
				"order %q is not registered in accrual: %w", order.ID, err)
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
				err = f.increaseBalance(ctx, update)
				if err != nil {
					return fmt.Errorf("cannot increase balance: %w", err)
				}
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

func (f *FeedAccrual) increaseBalance(
	ctx context.Context, order entities.Order,
) error {
	if order.Accrual == nil {
		return nil
	}

	err := f.storage.IncreaseBalance(
		ctx, order.User, *order.Accrual,
	)
	if err != nil {
		return fmt.Errorf("storage failed to increase balance: %w", err)
	}
	f.logger.Infof(
		"user %q balance increased by %v",
		order.User, order.Accrual)
	return nil
}
