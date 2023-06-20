package models

// ExternalData - struct for external loyalty system
type ExternalData struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
