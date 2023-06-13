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
	"time"
)

func Run(logger *zap.Logger) {
	//init config
	config := config.NewConfig()
	config.SetConfig()

	//param := "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable"
	//create model
	db := database.NewPostgresDB(config.Database, logger)
	defer db.Close()

	err := db.Create()
	if err != nil {
		logger.Sugar().Errorf("Error creating DB: %s", err)
	}
	orderProcessor := orderprocessor.NewOrderProcessor(db, logger)
	go orderProcessor.Process(config.AccrualSystemAddress)

	router := httpserver.RequestRouter(db, logger)

	server := &http.Server{Addr: config.ServerAddress, Handler: router}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second) //nolint:govet
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Sugar().Errorf("Shutdown timeout exceeded")
			}
		}()
		logger.Sugar().Infof("Stopping server")
		err := server.Shutdown(shutdownCtx)
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
