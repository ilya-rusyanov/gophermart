package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/ports"
)

type Withdraw struct {
	storage ports.WithdrawalTransactioner
}

func NewWithdraw(storage ports.WithdrawalTransactioner) *Withdraw {
	return &Withdraw{
		storage: storage,
	}
}

func (w *Withdraw) Withdraw(
	ctx context.Context, request entities.WithdrawalRequest,
) error {
	tx, err := w.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	currentBalance, err := tx.GetCurrentBalance(ctx, request.User)
	if err != nil {
		return fmt.Errorf("failed to retrieve current balance: %w", err)
	}

	if currentBalance < request.Sum {
		return entities.ErrInsufficientBalance
	}

	err = tx.DecreaseBalance(ctx, request.User, request.Sum)
	if err != nil {
		return fmt.Errorf("failed to decrease balance: %w", err)
	}

	err = tx.IncreaseWithdrawn(ctx, request.User, request.Sum)
	if err != nil {
		return fmt.Errorf("failed to increase withdrawn amount: %w", err)
	}

	record := entities.WithdrawalRecord{
		User:        request.User,
		Order:       request.Order,
		Sum:         request.Sum,
		ProcessedAt: time.Now(),
	}

	err = tx.RecordWithdrawal(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to record withdrawal: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
