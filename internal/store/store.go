package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"time"
)


type Store struct {
	TypeDB string
	db    *sql.DB
	dbURL string
	cch *cache.Cache
}


func Create(TypeDB, dbURL string) *Store {
	return &Store{
		TypeDB: TypeDB,
		db:    &sql.DB{},
		dbURL: dbURL,
		cch : cache.New(5*time.Minute, 10*time.Minute),
	}
}


func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.dbURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}


func (s *Store) Close() {
	s.db.Close()
}


func CreateTEST() *Store{
	return &Store{
		TypeDB: "Postgres",
		db: &sql.DB{},
		dbURL: "postgres://postgres:qwerty@database:5432/postgres?sslmode=disable",
	}
}
