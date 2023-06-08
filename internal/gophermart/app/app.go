package app

import (
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/database"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/httpserver"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/orderprocessor"
	"go.uber.org/zap"
	"net/http"
)

func Run(logger *zap.Logger) {
	//init config
	param := "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable"
	//create model
	db := database.NewPostgresDB(param, logger)

	defer db.Close()

	err := db.Create()
	if err != nil {
		logger.Sugar().Errorf("Error creating DB: %s", err)
	}
	orderProcessor := orderprocessor.NewOrderProcessor(db, logger)
	go orderProcessor.Process()

	router := httpserver.RequestRouter(db, logger)
	ServerAddress := ":8080"
	logger.Sugar().Fatalf("Server stopped with error: %s",
		http.ListenAndServe(ServerAddress, router))

}
