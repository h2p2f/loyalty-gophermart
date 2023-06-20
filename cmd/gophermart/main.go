package main

import (
	"fmt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/app" //nolint:typecheck
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
)

func main() {
	// Init logger
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}
	//start app
	app.Run(logger.Log)

}
