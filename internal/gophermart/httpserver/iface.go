package httpserver

import "time"

type DataBaser interface {
	NewUser(login, password string) error
	NewOrder(id, login, status string, accrual int, timeCreated time.Time) error
	GetOrdersByUser(login string) ([]byte, error)
	CheckUniqueOrder(order string) (string, bool)
	FindPassByLogin(login string) (string, error)
	GetBalance(login string) (int, error)
	GetSumOfAllWithdraws(login string) int
	NewWithdraw(login, order string, amount int, timeCreated time.Time) error
	GetAllWithdraws(login string) []byte
}
