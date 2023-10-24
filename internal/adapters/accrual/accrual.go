package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Logger interface {
	Infof(string, ...any)
}

type Adapter struct {
	addr   string
	logger Logger
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

func New(logger Logger, addr string) *Adapter {
	return &Adapter{
		addr:   addr + "/api/orders/",
		logger: logger,
	}
}

func (a *Adapter) GetStateOfOrder(ctx context.Context, orderID entities.OrderID) (
	status entities.OrderStatus, value entities.Currency, err error,
) {
	var path string

	path, err = url.JoinPath(a.addr, string(orderID))
	if err != nil {
		err = fmt.Errorf("failed to construct url: %w", err)
		return
	}
	a.logger.Infof("getting state of order %q", path)

	resp, err := http.Get(path)
	if err != nil {
		err = fmt.Errorf("http request error: %w", err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		status, value, err = a.readOrderState(resp.Body)
		if err != nil {
			err = fmt.Errorf(
				"failed to read status from body: %w", err)
			return
		}
	case http.StatusNoContent:
		err = entities.ErrAccrualOrderIsNotRegistered
		return
	case http.StatusTooManyRequests:
		var period time.Duration
		period, err = a.readPeriod(resp.Header)
		if err != nil {
			err = errors.New("failed to read delay from response header")
			return
		}
		err = &entities.AccrualTooManyRequestsError{
			Period: period,
		}
		return
	default:
		err = fmt.Errorf("unknown accrual response status %d", resp.StatusCode)
		return
	}

	return
}

func (a *Adapter) readOrderState(reader io.Reader) (
	status entities.OrderStatus, value entities.Currency, err error,
) {
	var (
		response accrualResponse
		ok       bool
	)

	err = json.NewDecoder(reader).Decode(&response)
	if err != nil {
		err = fmt.Errorf("error decoding JSON: %w", err)
		return
	}

	status, ok = toEntityStatus[response.Status]
	if !ok {
		err = fmt.Errorf("unknown order status %q", response.Status)
		return
	}

	value = entities.Currency(response.Accrual)

	return
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
