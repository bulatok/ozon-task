package main

import (
	"flag"
	"github.com/bulatok/ozon-task/configs"
	"github.com/bulatok/ozon-task/internal/server"
	"github.com/bulatok/ozon-task/internal/store"
	_ "github.com/lib/pq"
	"log"
)

var (
	StoreType  string
)

func init() {
	flag.StringVar(&StoreType, "store_type", "Postgres", "choose the type of data base")
}
func main() {
	flag.Parse()

	config, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	var str store.Store
	switch StoreType {
	case "Postgres":
		str, err = store.CreatePostgreDB(config)
		if err != nil{
			log.Fatal(err)
		}
	case "in-memory":
		str = store.CreateInMemory()
	default:
		log.Fatal("unknown database")
	}

	if err := server.Start(config, str); err != nil {
		log.Fatal(err)
	}
}
