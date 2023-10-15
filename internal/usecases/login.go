package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type LoginStorage interface {
	FindCredentials(context.Context, entities.AuthCredentials) error
}

type Login struct {
	storage    LoginStorage
	logger     Logger
	expiration time.Duration
	key        string
}

func NewLogin(logger Logger, tokenExpiration time.Duration, signingKey string, storage LoginStorage) *Login {
	return &Login{
		storage:    storage,
		logger:     logger,
		expiration: tokenExpiration,
		key:        signingKey,
	}
}

func (l *Login) Auth(ctx context.Context, creds entities.AuthCredentials) (
	entities.AuthToken, error,
) {
	var result entities.AuthToken

	l.logger.Infof("attempt to login user %q", creds.Login)

	err := l.storage.FindCredentials(ctx, creds)
	switch {
	case errors.Is(err, entities.ErrNotFound):
		return result, entities.ErrLoginIncorrect
	case err == nil:
		// all okay
		break
	default:
		return result, fmt.Errorf("unexpected storage error: %w", err)
	}

	result, err = buildAuthToken(l.expiration, creds.Login, l.key)
	if err != nil {
		return result, fmt.Errorf("failed to build auth token: %w", err)
	}

	l.logger.Infof("user %q authenticated", creds.Login)

	return result, nil
}
