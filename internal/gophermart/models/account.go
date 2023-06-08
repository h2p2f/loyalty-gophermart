package models

type Account struct {
	Balance   int `json:"current"`
	Withdraws int `json:"withdrawn,omitempty"`
}
