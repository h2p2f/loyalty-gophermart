package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
	"go.uber.org/zap"
)

func RequestRouter(db DataBaser, log *zap.Logger) chi.Router {
	handler := NewGopherMartHandler(db, log)

	r := chi.NewRouter()

	r.Use(logger.WithLogging, JWTAuth)
	r.Post("/api/user/register", handler.Register)
	r.Post("/api/user/login", handler.Login)
	r.Post("/api/user/orders", handler.AddOrder)
	r.Get("/api/user/orders", handler.Orders)
	r.Get("/api/user/balance", handler.Balance)
	r.Post("/api/user/balance/withdraw", handler.Withdraw)
	r.Get("/api/user/withdrawals", handler.Withdrawals)

	return r
}
