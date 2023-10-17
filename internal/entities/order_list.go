package entities

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         OrderID     `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *Currency   `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type OrderList []Order
