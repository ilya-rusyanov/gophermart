package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		c.errorHandler(rw,
			fmt.Errorf(
				"failed to read data from request body: %w",
				entities.ErrInvalidOrder))
		return
	}

	id := buf.String()
	if len(id) == 0 {
		c.errorHandler(rw,
			fmt.Errorf(
				"failed to parse order ID from request body: %w",
				entities.ErrInvalidOrder))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.errorHandler(rw,
			fmt.Errorf("failed to read integer from order ID: %w",
				entities.ErrInvalidOrder))
	}

	if !luhn.Valid(int(idInt)) {
		c.errorHandler(rw, entities.ErrLuhnValidation)
		return
	}

	createRequest := entities.CreateOrderRequest{
		ID:   entities.OrderID(id),
		User: getUser(r.Context()),
		Time: time.Now(),
	}

	c.logger.Infof("request to create order id %q for user %q",
		createRequest.ID, createRequest.User)

	err = c.usecase.CreateOrder(r.Context(), createRequest)
	if err != nil {
		c.errorHandler(rw, fmt.Errorf(
			"usecase order creation error: %w", err))
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	c.logger.Infof("order %q created", id)
}
