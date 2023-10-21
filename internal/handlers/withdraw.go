package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Withdrawer interface {
	Withdraw(context.Context, entities.WithdrawalRequest) error
}

type Withdraw struct {
	logger       Logger
	withdrawer   Withdrawer
	errorHandler ErrorHandler
}

func NewWithdraw(
	logger Logger, withdrawer Withdrawer, errorHandler ErrorHandler,
) *Withdraw {
	return &Withdraw{
		logger:       logger,
		withdrawer:   withdrawer,
		errorHandler: errorHandler,
	}
}

func (w *Withdraw) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var wr entities.WithdrawalRequest

	err := json.NewDecoder(r.Body).Decode(&wr)
	if err != nil {
		w.errorHandler(rw,
			fmt.Errorf("failed to parse JSON: %w", errParsing),
		)
		return
	}

	wr.User = getUser(r.Context())

	err = w.withdrawer.Withdraw(r.Context(), wr)
	if err != nil {
		w.errorHandler(rw, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
