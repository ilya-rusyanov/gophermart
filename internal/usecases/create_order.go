package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/ports"
)

type CreateOrder struct {
	logger  Logger
	storage ports.CreateOrderTransactioner
}

func NewCreateOrder(logger Logger, storage ports.CreateOrderTransactioner) *CreateOrder {
	return &CreateOrder{
		logger:  logger,
		storage: storage,
	}
}

func (o *CreateOrder) CreateOrder(
	ctx context.Context,
	req entities.CreateOrderRequest,
) error {
	tx, err := o.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start storage transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	storageUser, err := tx.FindUserForOrder(ctx, req.ID)
	switch {
	case errors.Is(err, entities.ErrNotFound):
		// all okay, may proceed
		break
	case err == nil:
		if storageUser == req.User {
			return entities.ErrAlreadyUploaded
		} else {
			return entities.ErrAlreadyUploadedOtherUser
		}
	default:
		return fmt.Errorf("unexpected storage error: %w", err)
	}

	err = tx.CreateOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("storage failed to create order: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit create order transaction: %w", err)
	}

	// TODO enqueue to accrual

	return nil
}
