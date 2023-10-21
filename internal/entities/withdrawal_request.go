package entities

type WithdrawalRequest struct {
	User  Login
	Order OrderID  `json:"order"`
	Sum   Currency `json:"sum"`
}
