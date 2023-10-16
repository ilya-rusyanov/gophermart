package usecases

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type CreateOrder struct {
	logger Logger
}

func NewCreateOrder(logger Logger) *CreateOrder {
	return &CreateOrder{
		logger: logger,
	}
}

func (o *CreateOrder) CreateOrder(
	ctx context.Context,
	req entities.CreateOrderRequest,
) error {
	return errors.New("TODO")
}
