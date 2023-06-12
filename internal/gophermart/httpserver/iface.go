package httpserver

import "time"

type DataBaser interface {
	NewUser(login, password string) error
	NewOrder(id, login, status string, accrual float64, timeCreated time.Time) error
	GetOrdersByUser(login string) ([]byte, error)
	CheckUniqueOrder(order string) (string, bool)
	FindPassByLogin(login string) (string, error)
	GetBalance(login string) (float64, error)
	GetSumOfAllWithdraws(login string) float64
	NewWithdraw(login, order string, amount float64, timeCreated time.Time) error
	GetAllWithdraws(login string) []byte
}
