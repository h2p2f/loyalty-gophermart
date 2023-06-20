package main

import (
	"fmt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/app" //nolint:typecheck
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
)

func main() {

	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}

	app.Run(logger.Log)

}
