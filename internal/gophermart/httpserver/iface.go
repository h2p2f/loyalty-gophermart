package httpserver

import (
	"context"

	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
)

// DataBaser interface for database
//
//go:generate mockery --name DataBaser --output ./mocks --filename mocks_databaser.go
type DataBaser interface {
	NewUser(ctx context.Context, user models.User) error
	NewOrder(ctx context.Context, login string, order models.Order) error
	GetOrdersByUser(ctx context.Context, login string) ([]byte, error)
	CheckUniqueOrder(ctx context.Context, order string) (string, error)
	FindPassByLogin(ctx context.Context, login string) (string, error)
	GetBalance(ctx context.Context, login string) (float64, error)
	GetSumOfAllWithdraws(ctx context.Context, login string) float64
	NewWithdraw(ctx context.Context, login string, withdraw models.Withdraw) error
	GetAllWithdraws(ctx context.Context, login string) []byte
}
