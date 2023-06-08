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
	Number      string    `json:"order"`
	Status      string    `json:"status"`
	Accrual     int       `json:"accrual,omitempty"`
	TimeCreated time.Time `json:"uploaded_at,omitempty"`
}
