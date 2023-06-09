package orderprocessor

type processor interface {
	GetUnfinishedOrders() (map[string]string, error)
	UpdateOrderStatus(order, status string, accrual int) error
}
