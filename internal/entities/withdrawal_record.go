package entities

import "time"

type WithdrawalRecord struct {
	User        Login
	Order       OrderID
	Sum         Currency
	ProcessedAt time.Time
}
