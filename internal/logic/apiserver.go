package logic

import (
	"github.com/bulatok/ozon-task/internal/server"
	"github.com/bulatok/ozon-task/internal/store"
	"log"
	"net/http"
)
func Start(config *Config, dbType string) error {
	log.Println("Start listening on port :8080")
	s := &store.Store{}
	switch dbType {
	case "Postgres":
		s = store.Create("Postgres", config.DatabaseURL)
		if err := s.Open(); err != nil {
			return err
		}

		defer s.Close()

		srv := server.NewServer(s)
		return http.ListenAndServe(config.Port, srv)
	default:
		s = store.Create("In-memory", config.DatabaseURL)
	}
	srv := server.NewServer(s)
	return http.ListenAndServe(config.Port, srv)
}
