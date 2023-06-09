package models

import "time"

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	TimeCreated time.Time `json:"processed_at,omitempty"`
}
