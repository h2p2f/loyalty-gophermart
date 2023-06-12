package app

import (
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/config"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/database"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/httpserver"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/orderprocessor"
	"go.uber.org/zap"
	"net/http"
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

	logger.Sugar().Fatalf("Server stopped with error: %s",
		http.ListenAndServe(config.ServerAddress, router))

}
