package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ilya-rusyanov/gophermart/internal/entities"

	"github.com/theplant/luhn"
)

type OrderCreator interface {
	CreateOrder(context.Context, entities.CreateOrderRequest) error
}

type OrderCreation struct {
	usecase      OrderCreator
	errorHandler ErrorHandler
	logger       Logger
}

func NewOrderCreation(logger Logger, usecase OrderCreator, errorHandler ErrorHandler) *OrderCreation {
	return &OrderCreation{
		usecase:      usecase,
		errorHandler: errorHandler,
		logger:       logger,
	}
}

func (c *OrderCreation) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	buf := &strings.Builder{}
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		c.errorHandler(rw, fmt.Errorf("failed to read data from request body: %w", entities.ErrInvalidOrder))
		return
	}
	id, err := strconv.ParseInt(buf.String(), 10, 64)
	if err != nil {
		c.errorHandler(rw, fmt.Errorf("failed to parse order ID from body: %w", entities.ErrInvalidOrder))
		return
	}
	if !luhn.Valid(int(id)) {
		c.errorHandler(rw, fmt.Errorf("failed to verify Luhn: %w", err))
		return
	}
	createRequest := entities.CreateOrderRequest{
		ID:   entities.OrderID(id),
		User: getUser(r.Context()),
	}
	c.logger.Infof("request to create order id %d for user %q",
		createRequest.ID, createRequest.User)
	err = c.usecase.CreateOrder(r.Context(), createRequest)
	if err != nil {
		c.errorHandler(rw, fmt.Errorf(
			"usecase order creation error: %w", err))
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}
