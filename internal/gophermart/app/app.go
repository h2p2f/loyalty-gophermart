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
	//create postgresql db
	db := database.NewPostgresDB(conf.Database, logger)
	defer db.Close()
	//create table
	err := db.Create()
	if err != nil {
		logger.Sugar().Errorf("Error creating DB: %s", err)
	}
	//create context for server
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	//create router
	router := httpserver.RequestRouter(db, logger)
	//create server
	server := &http.Server{Addr: conf.ServerAddress, Handler: router}
	//create order processor
	orderProcessor := orderprocessor.NewOrderProcessor(db, logger)
	go orderProcessor.Process(serverCtx, conf.AccrualSystemAddress)
	//create channel for stop server
	sig := make(chan os.Signal, 1)
	//notify about stop server
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//goroutine to stop server
	go func() {
		<-sig
		logger.Sugar().Infof("Stopping server")
		err := server.Shutdown(context.Background())
		if err != nil {
			logger.Sugar().Errorf("Error stopping server: %s", err)
		}
		serverStopCtx()
	}()
	//start server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Sugar().Errorf("Error starting server: %s", err)
	}
	//wait for stop server
	<-serverCtx.Done()
}
