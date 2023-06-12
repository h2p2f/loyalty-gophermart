package models

import (
	"time"
)

const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

type Order struct {
	Number      string    `json:"number"`
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual,omitempty"`
	TimeCreated time.Time `json:"uploaded_at,omitempty"`
}
