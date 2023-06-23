package httpserver

import (
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// RequestRouter create router
func RequestRouter(db DataBaser, log *zap.Logger, key string) chi.Router {
	handler := NewGopherMartHandler(db, log)

	r := chi.NewRouter()

	r.Use(logger.WithLogging, GzipHanler)

	r.Route("/api/user", func(r chi.Router) {

		r.Use(middleware.WithValue("key", key))

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(JWTAuth)

			r.Post("/orders", handler.AddOrder)
			r.Post("/balance/withdraw", handler.DoWithdraw)

			r.Get("/orders", handler.GetOrders)
			r.Get("/balance", handler.GetBalance)
			r.Get("/withdrawals", handler.GetWithdrawals)
		})
	})
	return r
}
