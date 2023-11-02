package handlers

import (
	"context"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

func getUser(ctx context.Context) entities.Login {
	return ctx.Value(ContextKeyLogin).(entities.Login)
}
