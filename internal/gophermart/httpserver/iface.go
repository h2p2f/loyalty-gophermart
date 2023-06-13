package httpserver

import (
	"context"
	"time"
)

//go:generate mockery --name DataBaser --output ./mocks --filename mocks_databaser.go
type DataBaser interface {
	NewUser(ctx context.Context, login, password string) error
	NewOrder(ctx context.Context, id, login, status string, accrual float64, timeCreated time.Time) error
	GetOrdersByUser(ctx context.Context, login string) ([]byte, error)
	CheckUniqueOrder(ctx context.Context, order string) (string, bool)
	FindPassByLogin(ctx context.Context, login string) (string, error)
	GetBalance(ctx context.Context, login string) (float64, error)
	GetSumOfAllWithdraws(ctx context.Context, login string) float64
	NewWithdraw(ctx context.Context, login, order string, amount float64, timeCreated time.Time) error
	GetAllWithdraws(ctx context.Context, login string) []byte
}
