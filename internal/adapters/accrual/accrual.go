package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Adapter struct {
	addr string
}

type orderStatus string

const (
	orderStatusRegistered orderStatus = "REGISTERED"
	orderStatusInvalid    orderStatus = "INVALID"
	orderStatusProcessing orderStatus = "PROCESSING"
	orderStatusProcessed  orderStatus = "PROCESSED"
)

type accrualResponse struct {
	//Order string `json:"order"`
	Status  orderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}

var toEntityStatus = map[orderStatus]entities.OrderStatus{
	orderStatusRegistered: entities.OrderStatusNew,
	orderStatusInvalid:    entities.OrderStatusInvalid,
	orderStatusProcessing: entities.OrderStatusProcessing,
	orderStatusProcessed:  entities.OrderStatusProcessed,
}

func New(addr string) *Adapter {
	return &Adapter{
		addr: addr,
	}
}

func (a *Adapter) GetStateOfOrder(ctx context.Context, orderID entities.OrderID) (
	entities.OrderStatus, error,
) {
	var currentStatus entities.OrderStatus
	resp, err := http.Get(a.addr)
	if err != nil {
		return currentStatus, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		currentStatus, err = a.readOrderState(resp.Body)
		if err != nil {
			return currentStatus, fmt.Errorf(
				"failed to read status from body: %w", err)
		}
	case 204:
		return currentStatus, entities.ErrAccrualOrderIsNotRegistered
	case 429:
		period, err := a.readPeriod(resp.Header)
		if err != nil {
			return currentStatus, errors.New("failed to read delay from response header")
		}
		return currentStatus, &entities.AccrualTooManyRequestsError{
			Period: period,
		}
	default:
		return currentStatus, fmt.Errorf("unknown accrual response status %d", resp.StatusCode)
	}
	return currentStatus, errors.New("TODO")
}

func (a *Adapter) readOrderState(reader io.Reader) (entities.OrderStatus, error) {
	var (
		result   entities.OrderStatus
		response accrualResponse
		ok       bool
	)

	err := json.NewDecoder(reader).Decode(&response)
	if err != nil {
		return result, fmt.Errorf("error decoding JSON: %w", err)
	}

	result, ok = toEntityStatus[response.Status]
	if !ok {
		return result, fmt.Errorf("unknown order status %q", response.Status)
	}

	return result, nil
}

func (a *Adapter) readPeriod(header http.Header) (time.Duration, error) {
	var result time.Duration

	retryAfter := header.Get("Retry-After")
	seconds, err := strconv.Atoi(retryAfter)
	if err != nil {
		return result, fmt.Errorf("failed to decode Retry-After header: %w", err)
	}

	result = time.Duration(seconds) * time.Second

	return result, nil
}
