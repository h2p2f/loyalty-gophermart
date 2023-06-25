package models

// User - struct for user
type Account struct {
	Balance   float64  `json:"current"`
	Withdraws *float64 `json:"withdrawn,omitempty"`
}
