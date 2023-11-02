package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type BalanceDisplayer interface {
	Show(context.Context, entities.Login) (entities.Balance, error)
}

type ShowBalance struct {
	logger       Logger
	displayer    BalanceDisplayer
	errorHandler ErrorHandler
}

func NewShowBalance(
	logger Logger, displayer BalanceDisplayer, errorHandler ErrorHandler,
) *ShowBalance {
	return &ShowBalance{
		logger:       logger,
		displayer:    displayer,
		errorHandler: errorHandler,
	}
}

func (s *ShowBalance) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.logger.Info("current balance request")

	user := getUser(r.Context())

	balance, err := s.displayer.Show(r.Context(), user)
	if err != nil {
		err = fmt.Errorf("failed to retrieve balance for user %q: %w", user, err)
		s.errorHandler(rw, err)
		return
	}

	s.logger.Infof("show balance success for user %q", user)

	encodeJSON(rw, s.errorHandler, &balance, http.StatusOK)
}
