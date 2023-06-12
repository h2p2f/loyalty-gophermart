package models

type Account struct {
	Balance   float64 `json:"current"`
	Withdraws float64 `json:"withdrawn,omitempty"`
}
