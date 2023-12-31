package handlers

import (
	"context"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type ListOrdersUsecase interface {
	ListOrders(context.Context, entities.ListOrdersRequest,
	) (entities.OrderList, error)
}

type ListOrders struct {
	usecase      ListOrdersUsecase
	errorHandler ErrorHandler
	logger       Logger
}

func NewListOrders(
	logger Logger, usecase ListOrdersUsecase, errorHandler ErrorHandler,
) *ListOrders {
	return &ListOrders{
		usecase:      usecase,
		errorHandler: errorHandler,
		logger:       logger,
	}
}

func (l *ListOrders) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	l.logger.Info("listing orders")
	req := entities.ListOrdersRequest{
		Login: getUser(r.Context()),
	}
	list, err := l.usecase.ListOrders(r.Context(), req)
	if err != nil {
		l.errorHandler(rw, err)
		return
	}
	if len(list) == 0 {
		http.Error(rw, "no data", http.StatusNoContent)
		return
	}
	l.logger.Info("order list success")
	encodeJSON(rw, l.errorHandler, &list, http.StatusOK)
}
