package main

import (
	"fmt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/app"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
)

func main() {

	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}

	app.Run(logger.Log)

}

//TODO:
//1. Add unit tests
//2. Make processing INVALID order's status
//3. Gracefully shutdown Processing goroutine
