package orderprocessor

// processor - interface for order processor
//
//go:generate mockery --name processor --output ./mocks --filename mocks_processor.go
type processor interface {
	GetUnfinishedOrders() (map[string]string, error)
	UpdateOrderStatus(order, status string, accrual float64) error
}
