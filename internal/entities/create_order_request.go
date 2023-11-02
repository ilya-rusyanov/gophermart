package entities

import "time"

type CreateOrderRequest struct {
	ID   OrderID
	User Login
	Time time.Time
}
