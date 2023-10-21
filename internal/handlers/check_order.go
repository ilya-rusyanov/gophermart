package handlers

import (
	"fmt"
	"strconv"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/theplant/luhn"
)

func checkOrder(id entities.OrderID) error {
	idInt, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to read integer from order ID: %w",
			entities.ErrInvalidOrder)
	}

	if !luhn.Valid(int(idInt)) {
		return entities.ErrLuhnValidation
	}

	return nil
}
