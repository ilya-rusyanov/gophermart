package handlers

import (
	"context"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type WithdrawalsStorage interface {
	ListWithdrawals(context.Context, entities.Login) (entities.Withdrawals, error)
}

type ListWithdrawals struct {
	logger     Logger
	storage    WithdrawalsStorage
	errHandler ErrorHandler
}

func NewListWithdrawals(
	logger Logger, storage WithdrawalsStorage, errHandler ErrorHandler,
) *ListWithdrawals {
	return &ListWithdrawals{
		logger:     logger,
		storage:    storage,
		errHandler: errHandler,
	}
}

func (l *ListWithdrawals) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	user := getUser(r.Context())
	resp, err := l.storage.ListWithdrawals(r.Context(), user)
	if err != nil {
		l.errHandler(rw, err)
		return
	}

	if len(resp) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	encodeJSON(rw, l.errHandler, &resp, http.StatusOK)
}
