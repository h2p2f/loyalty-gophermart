package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
	"go.uber.org/zap"
)

// RequestRouter create router
func RequestRouter(db DataBaser, log *zap.Logger, key string) chi.Router {
	handler := NewGopherMartHandler(db, log)

	r := chi.NewRouter()
	//use middlewares
	r.Use(middleware.WithValue("key", key))
	r.Use(logger.WithLogging, GzipHanler)
	r.Use(JWTAuth)
	//add routes
	r.Post("/api/user/register", handler.Register)
	r.Post("/api/user/login", handler.Login)
	r.Post("/api/user/orders", handler.AddOrder)
	r.Get("/api/user/orders", handler.Orders)
	r.Get("/api/user/balance", handler.Balance)
	r.Post("/api/user/balance/withdraw", handler.Withdraw)
	r.Get("/api/user/withdrawals", handler.Withdrawals)

	return r
}
