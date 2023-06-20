package app

import (
	"context"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/config"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/database"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/httpserver"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/orderprocessor"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(logger *zap.Logger) {
	//init conf
	conf := config.NewConfig()
	conf.SetConfig()

	//param := "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable"
	//create model
	db := database.NewPostgresDB(conf.Database, logger)
	defer db.Close()

	err := db.Create()
	if err != nil {
		logger.Sugar().Errorf("Error creating DB: %s", err)
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	router := httpserver.RequestRouter(db, logger)
	server := &http.Server{Addr: conf.ServerAddress, Handler: router}

	orderProcessor := orderprocessor.NewOrderProcessor(db, logger)
	go orderProcessor.Process(serverCtx, conf.AccrualSystemAddress)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		logger.Sugar().Infof("Stopping server")
		err := server.Shutdown(context.Background())
		if err != nil {
			logger.Sugar().Errorf("Error stopping server: %s", err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Sugar().Errorf("Error starting server: %s", err)
	}

	<-serverCtx.Done()
}
