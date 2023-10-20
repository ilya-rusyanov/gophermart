package entities

import "time"

type Order struct {
	ID         string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *Currency   `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
