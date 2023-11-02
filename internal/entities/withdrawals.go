package entities

import "time"

type Withdrawal struct {
	Order       OrderID   `json:"order"`
	Sum         Currency  `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdrawals []Withdrawal
