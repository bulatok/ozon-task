package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
)

const (
	storeName = "postgres"
)

type Postgres struct {
	db *sql.DB
	*links
}

func Provide(conf *config.Postgres) (*Postgres, error) {
	db, err := sql.Open("postgres", conf.DstUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{
		db:    db,
		links: &links{db},
	}, nil
}

func (p *Postgres) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Name() string {
	return storeName
}
