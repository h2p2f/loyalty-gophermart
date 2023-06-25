package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/config"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/database"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/httpserver"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/orderprocessor"
)

func Run() {
	//init conf

	conf := config.GetConfig()

	if err := logger.InitLogger(conf.SystemParams.LogLevel); err != nil {
		fmt.Println(err)
	}

	db := database.NewPostgresDB(conf.Database.Database, logger.Log)
	defer db.Close()
	//create table
	err := db.Create()
	if err != nil {
		logger.Log.Sugar().Errorf("Error creating DB: %s", err)
	}
	//create context for server
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	//create router
	router := httpserver.RequestRouter(db, logger.Log, conf.SystemParams.Key)
	//create server
	server := &http.Server{Addr: conf.Server.Address, Handler: router}
	//create order processor
	orderProcessor := orderprocessor.NewOrderProcessor(db, logger.Log)
	go orderProcessor.Process(serverCtx, conf.AccrualSystem.Address)
	//create channel for stop server
	sig := make(chan os.Signal, 1)
	//notify about stop server
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//goroutine to stop server
	go func() {
		<-sig
		logger.Log.Sugar().Infof("Stopping server")
		err := server.Shutdown(context.Background())
		if err != nil {
			logger.Log.Sugar().Errorf("Error stopping server: %s", err)
		}
		serverStopCtx()
	}()
	//start server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Sugar().Errorf("Error starting server: %s", err)
	}
	//wait for stop server
	<-serverCtx.Done()
}
