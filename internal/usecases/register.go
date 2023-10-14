package usecases

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Register struct {
}

func NewRegister() *Register {
	return &Register{}
}

// register user
func (a *Register) Auth(context.Context, entities.AuthCredentials) (
	entities.AuthToken, error,
) {
	return entities.AuthToken(""), errors.New("TODO")
}
