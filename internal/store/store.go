package store

import (
	"database/sql"
	"github.com/bulatok/ozon-task/configs"
	"github.com/bulatok/ozon-task/internal/store/SQLdb"
	in_memory "github.com/bulatok/ozon-task/internal/store/in-memory"
	_ "github.com/lib/pq"
	cc "github.com/patrickmn/go-cache"
	"time"
)

type Store interface {
	Open() error
	Close() error
	AddLink(string, string) error
	FindByURL(string) (string, error)
}

func CreateInMemory() *in_memory.InMemory{
	return &in_memory.InMemory{
		Cache: cc.New(time.Minute * 5, time.Minute*10),
		Timeout: time.Minute * 5,
	}
}

func CreatePostgreDB(config *configs.Config) (*SQLdb.PostgreDB, error){
	db :=  &SQLdb.PostgreDB{
		DB : &sql.DB{},
		DbURL: config.DatabaseURL,
	}
	if err := db.Open(); err != nil{
		return nil, err
	}
	return db, nil
}
