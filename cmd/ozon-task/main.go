package main

import (
	"flag"
	"github.com/bulatok/ozon-task/internal/logic"
	_ "github.com/lib/pq"
	"log"
	//"go.mongodb.org/mongo-driver/mongo"
)

var (
	configPath string = "configs/config.yml"
	StoreType  string
)

func init() {
	flag.StringVar(&StoreType, "store_type", "Postgres", "choose the type of data base")
}
func main() {

	flag.Parse()

	if StoreType != "Postgres" && StoreType != "Inmemory" && StoreType != "In-memory" && StoreType != "inmemory" && StoreType != "in-memory"{
		log.Fatal("Invalid database")
	}

	config, err := logic.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := logic.Start(config, StoreType); err != nil {
		log.Fatal(err)
	}
}
