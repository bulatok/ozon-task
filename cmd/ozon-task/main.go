package main

import (
	"log"

	"github.com/bulatok/ozon-task/internal/ozon-task/app"
	"github.com/bulatok/ozon-task/internal/ozon-task/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(conf)
}
